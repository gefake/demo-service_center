package webhookutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"telegram_bot/pkg/application"
	"telegram_bot/pkg/bot"
	"telegram_bot/pkg/database"
	"telegram_bot/pkg/logger"
	"telegram_bot/pkg/service"
	"time"

	"github.com/joho/godotenv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Error(".env файл не существует")
	}
}

func sendMessage(b *bot.Bot, app *application.ApplicationForCall) {
	users, err := database.GetUsers()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for _, user := range users {
		date := time.Unix(int64(app.Date), 0).Format("02/01/2006 15:04")
		applicationData := fmt.Sprintf("📞 Номер телефона: `%s`\n🙍🏻‍♂️ Имя клиента: `%s`\n🕒 Дата подачи заявки: `%s`", app.PhoneNumber, app.Name, date)
		botMsg := tg.NewMessage(user.TelegramID, "`📋 Поступила новая заявка`\n\nДанные по заявке:\n"+applicationData)
		botMsg.ParseMode = "markdown"

		b.GetBot().Send(botMsg)
	}
}

func handleWebhookMessage(body []byte, b *bot.Bot) error {
	var webhookData map[string]interface{}

	err := json.Unmarshal(body, &webhookData)
	if err != nil {
		return err
	}

	postingType, ok := webhookData["type"].(string)
	if !ok {
		return errors.New("bad webhook data")
	}

	switch postingType {
	case "newTask":
		application, err := unpackTask(webhookData)

		if err != nil {
			return err
		}

		sendMessage(b, application)
		fmt.Println("Received newTask webhook data:", string(body))
	case "newTrustedUser":
		unpackedUser, err := unpackAddTrustUser(webhookData)
		if err != nil {
			return err
		}

		service.AddTrustedUser(unpackedUser)

		fmt.Println("Новый доверенный пользователь добавлен в массив")
	case "removeTrustedUser":
		unpackRemoveTrustUser(webhookData)
	default:
		return errors.New("unknown message type")
	}

	return nil
}

func handleWebhook(w http.ResponseWriter, r *http.Request, b *bot.Bot) {
	trustedIP := os.Getenv("TRUSTED_IP")

	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	if clientIP != trustedIP {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only post, my son", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app := application.ApplicationForCall{}
	err = json.Unmarshal(body, &app)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := handleWebhookMessage(body, b); err != nil {
		logger.Log.Error(err.Error())
		return
	}

	fmt.Println("Received webhook:", string(body))

	w.WriteHeader(http.StatusOK)
}

func StartListen(c chan *bot.Bot) {
	logger.Log.Info("Starting webhook listener")

	bot := <-c
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		handleWebhook(w, r, bot)
	})

	listenPort := os.Getenv("LISTEN_PORT")

	http.ListenAndServe(":"+listenPort, nil)
}

package main

import (
	"fmt"
	"os"
	"telegram_bot/pkg/bot"
	"telegram_bot/pkg/logger"
	"telegram_bot/pkg/service"
	"telegram_bot/pkg/webhookutil"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/joho/godotenv"
)
//
func main() {
	botChan := make(chan *bot.Bot)

	//service.GetTrustedUsers()
	fmt.Println("JWT TOKEN:")
	fmt.Println(service.TrustedUsers)

	if err := godotenv.Load(); err != nil {
		logger.Log.Error(".env файл не существует")
	}

	logger.Log.Info("Старт обработчика бота")

	apiKey, _ := os.LookupEnv("TELEGRAM_BOT_KEY")
	botAPI, err := tg.NewBotAPI(apiKey)
	if err != nil {
		logger.Log.Fatal(err)
	}

	botHandler := bot.New(botAPI)

	go webhookutil.StartListen(botChan)
	botChan <- botHandler

	if err := botHandler.Start(); err != nil {
		logger.Log.Fatal(err)
	}
}

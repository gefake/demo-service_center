package bot

import (
	"fmt"
	"telegram_bot/pkg/database"
	"telegram_bot/pkg/logger"
	"telegram_bot/pkg/service"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommand(userMsg *tg.Message) {
	logger.Log.Info("COMMAND [" + userMsg.From.UserName + "] " + userMsg.Text)

	switch userMsg.Command() {
	case "start":
		if !service.IsTrustedUsername(userMsg.From.UserName) {
			botMsg := tg.NewMessage(userMsg.From.ID, "Вы не находитесь в списке доверенных пользователей для оповещений.\n\nЕсли вы являетесь администратором вы можете добавить себя в админ-панели.")
			b.bot.Send(botMsg)
		} else {
			botMsg := tg.NewMessage(userMsg.From.ID, "Привет, "+userMsg.From.FirstName+"!\n\nДобро пожаловать в бота для рассылки оповещений")

			b.bot.Send(botMsg)

			if !database.IsAuthorizedUser(userMsg.From.ID) {
				usr := &database.TelegramAuthUser{
					TelegramID:   userMsg.From.ID,
					TelegramName: userMsg.From.UserName,
				}

				botMsg = tg.NewMessage(userMsg.From.ID, "Авторизация проведена успешно под ид: *"+fmt.Sprintf("%d", userMsg.From.ID)+"*")
				botMsg.ParseMode = "markdown"
				b.bot.Send(botMsg)

				database.InsertUser(usr)
			}
		}
	}
}

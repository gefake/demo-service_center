package bot

import (
	"telegram_bot/pkg/logger"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot *tg.BotAPI
}

func New(bot *tg.BotAPI) *Bot {
	return &Bot{bot}
}

func (b *Bot) GetBot() *tg.BotAPI {
	return b.bot
}

func (b *Bot) Start() error {
	b.bot.Debug = false

	logger.Log.Info("Авторизация под аккаунтом " + b.bot.Self.UserName)

	//b.initCommands()
	updatesChan := b.initUpdatesChan()
	b.updateHandler(updatesChan)

	return nil
}

func (b *Bot) initUpdatesChan() tg.UpdatesChannel {
	u := tg.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) updateHandler(updates tg.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := update.Message

		if msg.IsCommand() {
			b.handleCommand(update.Message)
		}
	}
}

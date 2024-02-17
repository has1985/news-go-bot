package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/has1985/news-go-bot/internal/botkit"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот для чтения новостей.")
		if _, err := bot.Send(msg); err != nil {
			return err
		}
		return nil
	}

}

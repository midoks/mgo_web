package notify

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Notification struct {
	telegramBot *tgbotapi.BotAPI
	ChatID      int64
	Enabled     bool
}

func NewNotification(token string, chat_id int64, enabled bool) (*Notification, error) {
	if !enabled {
		return &Notification{Enabled: false}, nil
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	return &Notification{
		telegramBot: bot,
		ChatID:      chat_id,
		Enabled:     true,
	}, nil
}

func (n *Notification) Send(ctx context.Context, title, content string) error {
	if !n.Enabled {
		return nil
	}

	message := fmt.Sprintf("%s\n\n%s", title, content)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		msg := tgbotapi.NewMessage(n.ChatID, message)
		msg.ParseMode = "HTML"

		_, err := n.telegramBot.Send(msg)
		return err
	}
}

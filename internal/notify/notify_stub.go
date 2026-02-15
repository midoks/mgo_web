//go:build !telegram

package notify

import (
	"context"
)

type Notification struct {
	ChatID  int64
	Enabled bool
}

func NewNotification(token string, chat_id int64, enabled bool) (*Notification, error) {
	return &Notification{
		ChatID:  chat_id,
		Enabled: enabled,
	}, nil
}

func (n *Notification) Send(ctx context.Context, title, content string) error {
	if !n.Enabled {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

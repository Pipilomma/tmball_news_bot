package service

import (
	"context"

	"tmballNews/internal/domain"
)

func (s *service) OneByChatID(ctx context.Context, chatID int64) (*domain.Subs, error) {
	return s.db.OneByChatID(ctx, chatID)
}

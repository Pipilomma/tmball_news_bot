package service

import (
	"context"
	"tmballNews/internal/domain"
)

func (s *service) OneByChatIDAndUsername(ctx context.Context, chatID int64, username string) (*domain.Subs, error) {
	return s.db.OneByChatIDAndUsername(ctx, chatID, username)
}

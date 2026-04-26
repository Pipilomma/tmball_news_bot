package repository

import (
	"context"
	"time"

	"tmballNews/internal/domain"
)

type Postgres interface {
	Close(ctx context.Context) error
	SaveNews(ctx context.Context, news []domain.News) error
	SaveSubs(ctx context.Context, subs *domain.Subs) error
	LatestNews(ctx context.Context, now time.Time) (*domain.News, error)
	WeekNews(ctx context.Context, week time.Time) ([]domain.News, error)
	NewsExists(ctx context.Context, newsID string) bool
	FindNews(ctx context.Context, query string) (*domain.News, error)
	GetSubs(ctx context.Context) ([]domain.Subs, error)
	OneByChatID(ctx context.Context, chatID int64) (*domain.Subs, error)
}

type Parser interface {
	GetLatestNews(ctx context.Context) ([]domain.News, error)
}

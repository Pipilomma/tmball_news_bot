package repository

import (
	"context"
	"time"

	"tmballNews/internal/domain"
	"tmballNews/internal/repository/postgres/dao"
)

type Postgres interface {
	Close(ctx context.Context) error
	SaveNews(ctx context.Context, news []domain.News) error
	SaveSubs(ctx context.Context, subs *domain.Subs) error
	News(ctx context.Context, week time.Time) ([]domain.News, error)
	FindNews(ctx context.Context, query string) ([]domain.News, error)
	GetSubs(ctx context.Context) ([]domain.Subs, error)
	OneByChatID(ctx context.Context, chatID int64) (*domain.Subs, error)
	ListNews(ctx context.Context, filters *dao.NewsFilter) ([]domain.News, error)
}

type Parser interface {
	GetLatestNews(ctx context.Context) ([]domain.News, error)
}

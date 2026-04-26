package parser

import (
	"context"
	"tmballNews/internal/domain"
)

type Parser interface {
	GetLatestNews(ctx context.Context) ([]domain.News, error)
}

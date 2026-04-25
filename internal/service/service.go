package service

import (
	"context"
	"tmballNews/internal/domain"
	"tmballNews/internal/repository"
)

type Parser interface {
	GetLatestNews(ctx context.Context) ([]domain.News, error)
}

type service struct {
	db     repository.Postgres
	parser Parser
}

func New(db repository.Postgres, parser Parser) *service {
	return &service{
		db:     db,
		parser: parser,
	}
}

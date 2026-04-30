package service

import (
	"tmballNews/internal/repository"
)

type service struct {
	db     repository.Postgres
	parser repository.Parser
}

func New(db repository.Postgres, parser repository.Parser) *service {
	return &service{
		db:     db,
		parser: parser,
	}
}

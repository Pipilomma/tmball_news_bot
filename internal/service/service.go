package service

import (
	"tmballNews/internal/repository"
	"tmballNews/internal/repository/parser"
)

type service struct {
	db     repository.Postgres
	parser parser.Parser
}

func New(db repository.Postgres, parser parser.Parser) *service {
	return &service{
		db:     db,
		parser: parser,
	}
}

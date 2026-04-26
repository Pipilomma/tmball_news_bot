package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tmballNews/internal/domain"
	"tmballNews/internal/lib/errs"
	"tmballNews/internal/service/dto"
)

func (s *service) Subcribe(ctx context.Context, input dto.InputSubs) error {
	sb, err := s.db.OneByChatID(ctx, input.ChatID)
	if sb != nil {
		return errs.ErrUserUlreadySub
	}

	sub, err := domain.NewSubs(input.ChatID, input.Username, input.FirstName, true)
	if err != nil {
		return errors.New("failed to create new sub")
	}

	if err := s.db.SaveSubs(ctx, sub); err != nil {
		return err
	}

	return nil
}

func (s *service) ParseTmball(ctx context.Context) ([]domain.News, []domain.Subs, error) {
	news, err := s.parser.GetLatestNews(ctx)
	if err != nil {
		return nil, nil, errors.New("failed to parsed news")
	}

	if len(news) == 0 {
		return nil, nil, errors.New("zero news")
	}

	var newNews []domain.News

	if err := s.db.SaveNews(ctx, news); err != nil {
		return nil, nil, errors.New("failed to save news")
	}

	subs, err := s.db.GetSubs(ctx)
	if err != nil {
		return nil, nil, errors.New("failed to get subs")
	}

	return newNews, subs, nil
}

func (s *service) LastNews(ctx context.Context) (*domain.News, error) {
	now := time.Now()

	news, err := s.db.LatestNews(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest news: %w", err)
	}

	return news, nil
}

func (s *service) LastWeekNews(ctx context.Context) ([]domain.News, error) {
	now := time.Now()
	fromDate := now.AddDate(0, 0, -7)

	news, err := s.db.WeekNews(ctx, fromDate)
	if err != nil {
		return nil, errors.New("failed to get weeks news")
	}

	if len(news) == 0 {
		return nil, errors.New("zero news")
	}

	return news, nil
}

func (s *service) FindNews(ctx context.Context, message string) (*domain.News, error) {
	if message == "" {
		return nil, errors.New("empty query for search")
	}

	return s.db.FindNews(ctx, message)
}

package service

import (
	"context"
	"fmt"
	"time"

	"tmballNews/internal/domain"
	"tmballNews/internal/lib/errs"
	"tmballNews/internal/repository/postgres/dao"
	"tmballNews/internal/service/dto"
)

func (s *service) Subcribe(ctx context.Context, input dto.InputSubs) error {
	sb, err := s.db.OneByChatID(ctx, input.ChatID)
	if sb != nil {
		return errs.ErrUserUlreadySub
	}

	sub, err := domain.NewSubs(input.ChatID, input.Username, input.FirstName, true)
	if err != nil {
		return fmt.Errorf("failed to create new sub: %w", err)
	}

	if err := s.db.SaveSubs(ctx, sub); err != nil {
		return err
	}

	return nil
}

func (s *service) ParseTmball(ctx context.Context) ([]domain.News, []domain.Subs, error) {
	news, err := s.parser.GetLatestNews(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parsed news: %w", err)
	}

	if len(news) == 0 {
		return nil, nil, ErrZeroNews
	}

	var newNews []domain.News

	if err := s.db.SaveNews(ctx, news); err != nil {
		return nil, nil, fmt.Errorf("failed to save news: %w", err)
	}

	subs, err := s.db.GetSubs(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get subs: %w", err)
	}

	return newNews, subs, nil
}

func (s *service) LastNews(ctx context.Context) (*domain.News, error) {
	filters := dao.NewFilter(&time.Time{}, &dao.Now, "", dao.NewsSearchOff, dao.NewsOrderPublishedAtDesc)

	news, err := s.db.ListNews(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get last news: %w", err)
	}

	if len(news) == 0 {
		return nil, ErrZeroNews
	}

	return &news[0], nil
}

func (s *service) LastWeekNews(ctx context.Context) ([]domain.News, error) {
	filters := dao.NewFilter(&dao.WeekAgo, &dao.Now, "", dao.NewsSearchOff, dao.NewsOrderPublishedAtDesc)

	news, err := s.db.ListNews(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get weeks news: %w", err)
	}

	if len(news) == 0 {
		return nil, ErrZeroNews
	}

	return news, nil
}

func (s *service) FindNews(ctx context.Context, message string) (*domain.News, error) {
	if message == "" {
		return nil, fmt.Errorf("empty query for search")
	}

	filters := dao.NewFilter(nil, nil, message, dao.NewsSearchOn, dao.NewsOrderPublishedAtDesc)
	news, err := s.db.ListNews(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("cannot find news: %w", err)
	}

	if len(news) == 0 {
		return nil, errs.ErrNewsNotFound
	}

	return &news[0], nil
}

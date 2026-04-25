package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"tmballNews/internal/config"
	"tmballNews/internal/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var (
	newsTable = "news"
	subsTable = "subscribers"
)

type postgres struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
	cfg     *config.DBConfig
}

func NewPostgres(db *sqlx.DB, cfg *config.DBConfig) *postgres {
	return &postgres{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		cfg:     cfg,
	}
}

func (r *postgres) SaveNews(ctx context.Context, news *domain.News) error {
	now := time.Now().UTC()

	query, args, err := r.builder.Insert(newsTable).Columns(
		"id",
		"title",
		"content",
		"author_id",
		"status",
		"image_url",
		"video_url",
		"published_at",
		"created_at",
		"updated_at",
		"saved_at").Values(
		news.ID, news.Title, news.Content, news.AuthorID, news.Status, news.ImageURL,
		news.VideoURL, news.PublishedAt, news.CreatedAt, now, now).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgres) LatestNews(ctx context.Context, now time.Time) (*domain.News, error) {
	query, args, err := r.builder.
		Select(
			"id",
			"title",
			"content",
			"COALESCE(image_url, '') AS image_url",
			"COALESCE(video_url, '') AS video_url",
			"published_at",
		).
		From("news").
		OrderByClause(sq.Expr("ABS(EXTRACT(EPOCH FROM (published_at - ?)))", now)).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, errors.New("failed to get latest news")
	}

	row := r.db.QueryRowContext(ctx, query, args...)

	var n domain.News
	if err := row.Scan(&n.ID, &n.Title, &n.Content, &n.ImageURL, &n.VideoURL, &n.PublishedAt); err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *postgres) WeekNews(ctx context.Context, fromDate time.Time) ([]domain.News, error) {
	query, args, err := r.builder.
		Select(
			"id",
			"title",
			"content",
			"COALESCE(image_url, '') AS image_url",
			"COALESCE(video_url, '') AS video_url",
			"published_at",
		).
		From("news").
		Where(sq.GtOrEq{"published_at": fromDate}).
		OrderBy("published_at DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.News

	for rows.Next() {
		var n domain.News

		err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Content,
			&n.ImageURL,
			&n.VideoURL,
			&n.PublishedAt,
		)

		if err != nil {
			return nil, errors.New("failed to get one of week news")
		}

		result = append(result, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *postgres) NewsExists(ctx context.Context, newsID string) bool {
	query, args, err := r.builder.
		Select().
		Column(sq.Expr("EXISTS (SELECT 1 FROM news WHERE id = ?)", newsID)).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return false
	}

	var exists bool
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func (r *postgres) FindNews(ctx context.Context, query string) (*domain.News, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, nil
	}

	if news, err := r.findExact(ctx, q); err != nil {
		return nil, err
	} else if news != nil {
		return news, nil
	}

	if news, err := r.findByWords(ctx, q); err != nil {
		return nil, err
	} else if news != nil {
		return news, nil
	}

	return r.findFuzzy(ctx, q)
}

func (r *postgres) SaveSubs(ctx context.Context, subs *domain.Subs) error {
	now := time.Now().UTC()

	query, args, err := r.builder.Insert(subsTable).Columns(
		"id",
		"chat_id",
		"username",
		"first_name",
		"is_active",
		"created_at",
		"updated_at").Values(
		subs.ID, subs.ChatID, subs.Username, subs.FirstName, subs.IsActive, now, now).
		ToSql()

	if err != nil {
		return fmt.Errorf("хуй")
	}

	log.Println(subs.Username, subs.IsActive, subs.ChatID)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("пизда: %w", err)
	}

	return nil
}

func (r *postgres) OneByChatIDAndUsername(ctx context.Context, chatID int64, username string) (*domain.Subs, error) {
	query, args, err := r.builder.Select(
		"id",
		"chat_id",
		"username",
		"first_name",
		"is_active",
		"created_at",
		"updated_at",
	).From(subsTable).
		Where(sq.Eq{
			"chat_id": chatID, "username": username}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var res domain.Subs

	if err := r.db.GetContext(ctx, &res, query, args...); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *postgres) GetSubs(ctx context.Context) ([]domain.Subs, error) {
	query, args, err := r.builder.
		Select(
			"chat_id",
			"username",
			"first_name",
			"is_active",
			"created_at",
			"updated_at",
		).
		From("subscribers").
		Where("is_active = TRUE").
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Subs

	for rows.Next() {
		var s domain.Subs

		if err := rows.Scan(
			&s.ChatID,
			&s.Username,
			&s.FirstName,
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *postgres) findExact(ctx context.Context, q string) (*domain.News, error) {
	builder := sq.
		Select(
			"id",
			"title",
			"content",
			"COALESCE(image_url, '') AS image_url",
			"COALESCE(video_url, '') AS video_url",
			"published_at",
		).
		From("news").
		Where("LOWER(title) = LOWER(?)", q).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	return r.scanOne(ctx, builder)
}

func (r *postgres) findByWords(ctx context.Context, q string) (*domain.News, error) {
	words := strings.Fields(strings.ToLower(q))
	if len(words) == 0 {
		return nil, nil
	}

	conds := make(sq.And, 0, len(words))
	for _, w := range words {
		conds = append(conds, sq.Expr("LOWER(title) LIKE ?", "%"+w+"%"))
	}

	builder := sq.
		Select(
			"id",
			"title",
			"content",
			"COALESCE(image_url, '') AS image_url",
			"COALESCE(video_url, '') AS video_url",
			"published_at",
		).
		From("news").
		Where(conds).
		OrderBy("published_at DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	return r.scanOne(ctx, builder)
}

func (r *postgres) findFuzzy(ctx context.Context, q string) (*domain.News, error) {
	builder := sq.
		Select(
			"id",
			"title",
			"content",
			"COALESCE(image_url, '') AS image_url",
			"COALESCE(video_url, '') AS video_url",
			"published_at",
		).
		From("news").
		Where("title % ?", q).
		OrderByClause(
			sq.Expr("similarity(title, ?) DESC", q),
		).
		OrderBy("published_at DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	return r.scanOne(ctx, builder)
}

func (r *postgres) scanOne(ctx context.Context, builder sq.SelectBuilder) (*domain.News, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var n domain.News
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&n.ID,
		&n.Title,
		&n.Content,
		&n.ImageURL,
		&n.VideoURL,
		&n.PublishedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &n, nil
}

package postgres

import (
	"context"
	"database/sql"
	"errors"
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

func (r *postgres) SaveNews(ctx context.Context, news []domain.News) error {
	now := time.Now().UTC()

	builder := sq.Insert(newsTable).Columns(
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
		"saved_at").
		PlaceholderFormat(sq.Dollar).
		Suffix("ON CONFLICT (id) DO NOTHING")

	for _, n := range news {
		builder = builder.Values(
			n.ID,
			n.Title,
			n.Content,
			n.AuthorID,
			n.Status,
			n.ImageURL,
			n.VideoURL,
			n.PublishedAt,
			n.CreatedAt,
			n.UpdatedAt,
			now,
		)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.New("failed to create query for upsert news")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
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

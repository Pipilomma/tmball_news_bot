package postgres

import (
	"context"
	"log"
	"time"
	"tmballNews/internal/domain"

	sq "github.com/Masterminds/squirrel"
)

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
		return err
	}

	log.Println(subs.Username, subs.IsActive, subs.ChatID)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgres) OneByChatID(ctx context.Context, chatID int64) (*domain.Subs, error) {
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
			"chat_id": chatID}).
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

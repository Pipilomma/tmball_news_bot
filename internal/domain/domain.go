package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserState string

const (
	StateNone                  UserState = "none"
	StateAwaitingFindNewsInput UserState = "awaiting_find_news_input"
)

type News struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	AuthorID    string    `json:"authorId"`
	Status      string    `json:"status"`
	ImageURL    string    `json:"imageUrl"`
	VideoURL    string    `json:"videoUrl"`
	PublishedAt time.Time `json:"publishedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Subs struct {
	ID        uuid.UUID `db:"id"`
	ChatID    int64     `db:"chat_id"`
	Username  string    `db:"username"`
	FirstName string    `db:"first_name"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewSubs(chatID int64, username string, firstName string, isActive bool) (*Subs, error) {
	now := time.Now().UTC()

	sub := &Subs{
		ID:        uuid.New(),
		ChatID:    chatID,
		Username:  username,
		FirstName: firstName,
		IsActive:  isActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return sub, nil
}

type NewsListResponse struct {
	Articles []News `json:"articles"`
}

type NewsDetailResponse struct {
	Article News `json:"article"`
}

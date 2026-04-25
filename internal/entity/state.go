package entity

type UserState string

const (
	StateNone                  UserState = "none"
	StateAwaitingFindNewsInput UserState = "awaiting_find_news_input"
)

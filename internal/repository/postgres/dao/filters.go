package dao

import "time"

const (
	NewsOrderPublishedAtDesc = "published_at_desc"
)

const (
	NewsSearchOff = "off"
	NewsSearchOn  = "on"
)

var (
	Now     = time.Now()
	WeekAgo = Now.AddDate(0, 0, -7)
)

type NewsFilter struct {
	PublishedFrom *time.Time
	PublishedTo   *time.Time
	Query         string
	SearchMode    string
	Order         string
}

func NewFilter(pubFrom *time.Time, PubTo *time.Time, query string, searchMode string, order string) *NewsFilter {
	return &NewsFilter{
		PublishedFrom: pubFrom,
		PublishedTo:   PubTo,
		Query:         query,
		SearchMode:    searchMode,
		Order:         order,
	}
}

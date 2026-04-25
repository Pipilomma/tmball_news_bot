package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"tmballNews/internal/domain"
)

var (
	baseURL = "https://tmball.online/api/v1/public/news"
	agent   = "User-Agent"
	browser = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)"
)

type Parser struct {
	httpClient *http.Client
	baseURL    string
}

func New() *Parser {
	return &Parser{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (p *Parser) GetLatestNews(ctx context.Context) ([]domain.News, error) {
	url := fmt.Sprintf("%s?limit=12&page=1", p.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(agent, browser)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response domain.NewsListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Articles, nil
}

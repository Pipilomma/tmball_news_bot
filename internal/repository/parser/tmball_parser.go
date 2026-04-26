package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"tmballNews/internal/config"
	"tmballNews/internal/domain"
)

type TmballParser struct {
	httpClient *http.Client
	cfg        *config.ParserConfig
}

func New(cfg *config.ParserConfig) *TmballParser {
	return &TmballParser{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cfg: cfg,
	}
}

func (p *TmballParser) GetLatestNews(ctx context.Context) ([]domain.News, error) {
	url := fmt.Sprintf("%s?limit=12&page=1", p.cfg.TmballUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(url, err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to fetch news: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode, url)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response domain.NewsListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Articles, nil
}

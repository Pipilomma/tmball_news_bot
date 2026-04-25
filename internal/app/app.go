package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tmballNews/internal/config"
	"tmballNews/internal/domain"
	"tmballNews/internal/handler/tg"
	"tmballNews/internal/lib/closer"
	"tmballNews/internal/parser"
	"tmballNews/internal/repository"
	postgres "tmballNews/internal/repository/postgres"
	"tmballNews/internal/service"
)

type app struct {
	ctx       context.Context
	cancelCtx context.CancelFunc

	cfg *config.Config

	tgAPI *tg.API
	db    repository.Postgres

	parser  service.Parser
	service tg.Service
}

func New() (*app, error) {
	ctx, cancel := context.WithCancel(context.Background())

	a := &app{
		ctx:       ctx,
		cancelCtx: cancel,
	}

	err := a.initDeps()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *app) initDeps() error {
	fns := []func() error{
		a.initConfig,
		a.initDB,
		a.initServices,
	}

	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (a *app) initConfig() error {
	cfg := config.MustLoad()
	a.cfg = cfg

	return nil
}

func (a *app) initDB() error {
	db, err := postgres.New(a.ctx, &a.cfg.DB)
	if err != nil {
		return fmt.Errorf("init db error: %w", err)
	}

	a.db = db

	return nil
}

func (a *app) initServices() error {
	a.parser = parser.New()
	a.service = service.New(a.db, a.parser)

	if a.cfg.Telegram.Enabled {
		a.tgAPI = tg.NewAPI(a.ctx, &a.cfg.Telegram, a.service)
	}

	return nil
}

func (a *app) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	go func() {
		if a.tgAPI != nil {
			a.tgAPI.Register()
		}
	}()

	go func() {
		a.runParserLoop(a.ctx, a.cfg.ParserTimeout)
	}()

	log.Println("Server is running...")

	waitSignalAndShutdown(a.cancelCtx)

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.GracefulTimeout)
	defer cancel()

	return a.shutdown(ctx)
}

func (a *app) shutdown(ctx context.Context) error {
	if err := a.db.Close(ctx); err != nil {
		log.Printf("db close err: %v\n", err)
	}

	log.Println("Shutdown complete")

	return nil
}

func waitSignalAndShutdown(cancelApp context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGHUP, os.Interrupt,
	)

	sig := <-quit

	log.Printf("Caught signal %s. Shutting down...\n", sig)

	cancelApp()
}

func (a *app) runParserLoop(ctx context.Context, interval time.Duration) {
	runOnce := func() {
		news, subs, err := a.service.ParseTmball(ctx)
		if err != nil {
			log.Printf("parser run failed: %v", err)
			return
		}
		log.Printf("parser run done, new news: %d count subs: %d", len(news), len(subs))

		a.autoSendNews(news, subs)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			runOnce()
		}
	}
}

func (a *app) autoSendNews(news []domain.News, subs []domain.Subs) {
	if a.tgAPI == nil {
		return
	}

	for _, s := range subs {
		go func() {
			for _, n := range news {
				if err := a.tgAPI.SendTelegramNews(s.ChatID, &n); err != nil {
					log.Printf("failed send %s to %s: %v", n.Title, s.Username, err)
				}
			}
		}()
	}
}

package tg

import (
	"context"
	"log"
	"sync"

	"tmballNews/internal/config"
	"tmballNews/internal/domain"
	"tmballNews/internal/service/dto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler func(*tgbotapi.Message)
type CallbackHandler func(*tgbotapi.CallbackQuery)

type Service interface {
	ParseTmball(ctx context.Context) ([]domain.News, []domain.Subs, error)
	LastNews(ctx context.Context) (*domain.News, error)
	LastWeekNews(ctx context.Context) ([]domain.News, error)
	FindNews(ctx context.Context, message string) (*domain.News, error)
	Subcribe(ctx context.Context, input dto.InputSubs) error
	OneByChatID(ctx context.Context, chatID int64) (*domain.Subs, error)
}

type API struct {
	ctx              context.Context
	cfg              *config.TelegramConfig
	bot              *tgbotapi.BotAPI
	Service          Service
	commandHandlers  map[string]CommandHandler
	callbackHandlers map[string]CallbackHandler
	userStates       map[string]domain.UserState
	stateHandlers    map[domain.UserState]func(*tgbotapi.Message)
	mu               sync.Mutex
}

func NewAPI(ctx context.Context, cfg *config.TelegramConfig, service Service) *API {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalln(err)
	}

	bot.Debug = cfg.Debug

	api := &API{
		ctx:              ctx,
		cfg:              cfg,
		bot:              bot,
		Service:          service,
		commandHandlers:  make(map[string]CommandHandler),
		callbackHandlers: make(map[string]CallbackHandler),
		userStates:       make(map[string]domain.UserState),
		stateHandlers:    make(map[domain.UserState]func(*tgbotapi.Message)),
		mu:               sync.Mutex{},
	}

	api.setBotCommands()
	api.registerHandlers()

	return api
}

func (a *API) registerHandlers() {
	a.commandHandlers = map[string]CommandHandler{
		startCommand:        a.StartHandler,
		seeNewsCommand:      a.SeeNewsHandler,
		helpCommand:         a.HelpHandler,
		lastNewsCommand:     a.LastNewsHandler,
		lastWeekNewsCommand: a.LastWeekNewsHandler,
		findNewsCommand:     a.StartFindNewsHandler,
	}

	a.stateHandlers = map[domain.UserState]func(*tgbotapi.Message){
		domain.StateAwaitingFindNewsInput: a.FindNewsHandler,
	}
}

func (a *API) setBotCommands() {
	commands := []tgbotapi.BotCommand{
		{Command: startCommand, Description: "Запустить бота"},
		{Command: lastNewsCommand, Description: "Показать последнюю новость"},
		{Command: lastWeekNewsCommand, Description: "Показать новости текущей недели"},
		{Command: findNewsCommand, Description: "Найти новость по названию"},
		{Command: helpCommand, Description: "Помощь и команды"},
		{Command: seeNewsCommand, Description: "автоматическое получение новых новостей"},
	}

	if _, err := a.bot.Request(tgbotapi.NewSetMyCommands(commands...)); err != nil {
		log.Printf("Set commands error %v", err)
	}
}

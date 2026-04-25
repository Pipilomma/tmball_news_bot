package tg

import (
	"log"
	"runtime/debug"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *API) Register() {
	for update := range a.bot.GetUpdatesChan(tgbotapi.UpdateConfig{}) {
		go func(update tgbotapi.Update) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovering from panic: %v\nStack trace: %s", r, debug.Stack())
				}
			}()

			switch {
			case update.Message != nil && update.Message.IsCommand():
				a.handleCommand(update.Message)
			case update.Message != nil:
				a.handleState(update.Message)
			default:
				log.Printf("Unknown update type: %+v\n", update)
			}
		}(update)
	}
}

func (a *API) handleCommand(message *tgbotapi.Message) {
	if handler, exists := a.commandHandlers[message.Command()]; exists {
		handler(message)
	} else {
		a.UnknownCommandHandler(message)
	}
}

func (a *API) handleState(message *tgbotapi.Message) {
	userID := a.getUserID(message)
	state := a.getUserState(userID)

	if handler, exists := a.stateHandlers[state]; exists {
		handler(message)
	} else {
		a.UnknownCommandHandler(message)
	}
}

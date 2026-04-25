package tg

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *API) getInternalUserID(chatID int64, username string) (string, error) {
	subs, err := a.Service.OneByChatIDAndUsername(a.ctx, chatID, username)
	if err != nil {
		return "", err
	}

	id := subs.ID.String()

	return id, nil
}

func (a *API) getUserID(message *tgbotapi.Message) string {
	chatID := message.Chat.ID
	username := message.From.UserName

	internalUserID, err := a.getInternalUserID(chatID, username)
	if err != nil {
		log.Printf("Error getting internal user ID: %v", err)

		return ""
	}

	return internalUserID
}

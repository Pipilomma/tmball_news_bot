package tg

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"tmballNews/internal/domain"
	"tmballNews/internal/entity"
	"tmballNews/internal/lib/errs"
	"tmballNews/internal/service/dto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	reImg = regexp.MustCompile(`(?is)<img[^>]*>`)
	reA   = regexp.MustCompile(`(?is)<a[^>]*href="([^"]+)"[^>]*>`)
)

func (a *API) StartHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	text := startTextTelegram

	msg := tgbotapi.NewMessage(chatID, text)
	_, _ = a.bot.Send(msg)
}

func (a *API) HelpHandler(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, helpTextTelegram)
	_, _ = a.bot.Send(msg)
}

func (a *API) UnknownCommandHandler(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, unknowCommandTextTelegram)
	_, _ = a.bot.Send(msg)
}

func (a *API) LastNewsHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	msg := tgbotapi.NewMessage(chatID, lastNewsTextTelegram)
	_, _ = a.bot.Send(msg)

	news, err := a.Service.LastNews(a.ctx)
	if err != nil {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, lastNewsErr))
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, bringToSub))
		return
	}

	if news == nil {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, emptyNewsError))
		return
	}

	if err := a.SendTelegramNews(message.Chat.ID, news); err != nil {
		return
	}
}

func (a *API) LastWeekNewsHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	text := lastWeekNewsTextTelegram
	msg := tgbotapi.NewMessage(chatID, text)
	_, _ = a.bot.Send(msg)

	news, err := a.Service.LastWeekNews(a.ctx)
	if err != nil {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, lastWeekNewsErr))
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, bringToSub))
		return
	}

	if len(news) == 0 {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, zeroNewsLastWeekTextTelegram))
		return
	}

	for _, n := range news {

		if err := a.SendTelegramNews(message.Chat.ID, &n); err != nil {
			return
		}
	}

}

func (a *API) SeeNewsHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	err := a.Service.LetsSeeNews(a.ctx, dto.InputSubs{
		ChatID:    message.Chat.ID,
		Username:  message.Chat.UserName,
		FirstName: message.Chat.FirstName,
	})

	if err != nil {
		if errors.Is(err, errs.ErrUserUlreadySub) {
			_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, ulreadySubTextTelegram))
		}
		log.Println(err)
		return
	}

	msg := tgbotapi.NewMessage(chatID, letsSeeNewsTextTelegram)
	_, _ = a.bot.Send(msg)
}

func (a *API) FindNewsHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := a.getUserID(message)
	defer a.clearUserState(userID)

	input := strings.TrimSpace(message.Text)

	news, err := a.Service.FindNews(a.ctx, input)
	if err != nil {
		log.Println(err)
		return
	}

	if news == nil {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, unknowNewsTextTelegram))
		return
	}

	log.Println(news.Title)

	if err := a.SendTelegramNews(chatID, news); err != nil {
		log.Println(err)
		return
	}
}

func (a *API) StartFindNewsHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := a.getUserID(message)

	a.setUserState(userID, entity.StateAwaitingFindNewsInput)
	msg := tgbotapi.NewMessage(chatID, findNewsTextTelegram)
	_, _ = a.bot.Send(msg)
}

func (a *API) SendTelegramNews(chatID int64, news *domain.News) error {
	news.Content = toTelegramHTML(news.Content)
	news.Content = splitText(news.Content)

	text := toTelegramHTML(news.Title + "\n\n" + news.Content + "\n\n" + "открыть полностью: https://tmball.online/news/" + news.ID)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	_, err := a.bot.Send(msg)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func toTelegramHTML(raw string) string {
	s := raw
	s = reImg.ReplaceAllString(s, "")
	s = reA.ReplaceAllString(s, `<a href="$1">`)

	s = strings.NewReplacer(
		"<h1>", "<b>", "</h1>", "</b>\n\n",
		"<h2>", "<b>", "</h2>", "</b>\n\n",
		"<p>", "", "</p>", "\n\n",
		"<ul>", "", "</ul>", "\n",
		"<li>", "• ", "</li>", "\n",
		"<strong>", "<b>", "</strong>", "</b>",
		"<em>", "<i>", "</em>", "</i>",
		"<br>", "\n", "<br/>", "\n", "<br />", "\n",
	).Replace(s)

	return strings.TrimSpace(s)
}

func splitText(text string) string {
	lines := strings.Split(text, "\n")

	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		result = append(result, line)

		if len(result) == 2 {
			break
		}
	}

	return strings.Join(result, "\n")
}

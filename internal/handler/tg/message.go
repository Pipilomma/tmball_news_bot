package tg

// commands
const (
	startCommand        = "start"
	helpCommand         = "help"
	lastNewsCommand     = "last_news"
	lastWeekNewsCommand = "last_week_news"
	findNewsCommand     = "find_news"
	seeNewsCommand      = "lets_see_news"
)

// messages
const (
	startTextTelegram            = "Привет 👋! я бот для отправки новостей о футболе в Берендеевке из сервиса tmball.online ⚽️!\nЗдесь ты сможешь автоматически получать новые новости и анонсы предстоящих матчей 📱\nПросматривать последнюю новость из сервиса 👀\nА также читать новости текущей недели 📰!\nОтправь команду /lets_see_news, если хочешь получать уведомления о новых событиях противостояния магнита и пятерочки 😉"
	helpTextTelegram             = "Я умею делать следующие вещи:\n1 - автоматически отправлять последнюю новость из tmball.online 📱\n2 - по команде /last_news я пришлю самое свежее событие 👀\n3 - по команде /last_week_news я пришлю новости текущей недели 📰\n4 - отправь мне /find_news, после чего введи запрос на поиск новости по ее названию. Например, новость называется 'релиз 2026. Ответы на вопросы' - ты можешь отправить мне 'релиз 2026' и я найду соответсвующую новость 🔎\n\nЧтобы начать получать новости автоматически, введи /lets_see_news 🙌"
	lastNewsTextTelegram         = "Самая свежая новость 👀"
	lastWeekNewsTextTelegram     = "Новости текущей недели 📰"
	findNewsTextTelegram         = "Напиши название новости, которую ты хочешь увидеть, а я постараюсь ее найти 🔎"
	zeroNewsLastWeekTextTelegram = "На текущей недели пока нет новостей 🤷‍♂️"
	unknowCommandTextTelegram    = "Неизвестная комманда 🤔"
	letsSeeNewsTextTelegram      = "Вы подписались на рассылку новостей ✅"
	unknowNewsTextTelegram       = "Не смог найти данную новость 😔"
	ulreadySubTextTelegram       = "Вы уже подписаны на рассылку новостей 😊"
	bringToSub                   = "Возможно, новостей еще нет, и вам стоит подписаться на рассылку"
)

// errors
const (
	lastNewsErr      = "❌ Ошибка при получении последней новости ❌"
	lastWeekNewsErr  = "❌ Ошибка при получении новостей текущей недели❌"
	notFoundErr      = "❌ Данной новости не существует. Возвожно, вы ввели неверное название ❌"
	emptyNewsError   = "❌ Пустая новость ❌"
	letsSeeNewsError = "❌ Ошибка при автоматической отправке новости ❌"
)

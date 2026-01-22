package bot

import (
	"os"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CreateBot() {
	// считываем токен
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("Токен BOT_TOKEN не найден")
	}

	// создаем бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// включаем дебаг для отладки
	bot.Debug = true

	log.Printf("Авторизован как %s", bot.Self.UserName)

	// создаем канал для получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// логируем входящее сообщение
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// получаем входящее сообщение
		msgText := update.Message.Text
		chatId := update.Message.Chat.ID

		var replyText string

		switch {
			case msgText == "/start":
				replyText = "🚀 Привет! Я простой Telegram бот на Go. Напиши /help, чтобы узнать, что я умею."
			case msgText == "/help":
				replyText = `📚 Список доступных команд:
/start - запустить бота
/help - показать это сообщение`
			default:
				replyText = "🤖 Я не знаю, что ответить на это сообщение. Напиши /help, чтобы увидеть список команд."
		}
		
		// отправляем ответ
		msg := tgbotapi.NewMessage(chatId, replyText)
		bot.Send(msg)
	}
}

package bot

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sspyataev/simpleReminderBot/parser"
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
		fields := strings.SplitN(msgText, " ", 2)
		chatId := update.Message.Chat.ID

		var replyText string

		switch fields[0] {
		case "/start":
			replyText = "🚀 Привет! Я простой Telegram бот на Go. Напиши /help, чтобы узнать, что я умею."
		case "/help":
			replyText = `📚 Список доступных команд:
/start - запустить бота
/help - показать это сообщение
/add текст напоминания и когда напомнить - добавить напоминание
форматы времени: "через 15 минут", "завтра в 10 вечера", "в 18:00", "завтра", "через 2 часа", "в 9 вечера"`
		case "/add":
			prsr := parser.NewReminderParser()
			text, timePart := prsr.Parse(fields[1])
			fmt.Printf("Исходная фраза: %q\n", fields[1])
			fmt.Printf("  → Текст: %q\n", text)
			if timePart != "" {
				fmt.Printf("  → Время: %q\n", timePart)
			} else {
				fmt.Printf("  → Время: не указано\n")
			}
			fmt.Println()
		default:
			replyText = "🤖 Я не знаю, что ответить на это сообщение. Напиши /help, чтобы увидеть список команд."
		}

		// отправляем ответ
		msg := tgbotapi.NewMessage(chatId, replyText)
		bot.Send(msg)
	}
}

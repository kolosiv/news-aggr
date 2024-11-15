package telegram

import (
	"fmt"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kolosiv/news-aggr/internal/repository"
	"github.com/sirupsen/logrus"
)

func TelegramBot(nr repository.NewsRepository, stop <-chan struct{}) {
	botToken := os.Getenv("TG_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logrus.Error("error tg roken", err)
		return
	}

	// Установите уровень логирования в Debug, чтобы увидеть все запросы и ответы
	// bot.Debug = true

	logrus.Debug("Authorized on account " + bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message != nil {
				go handleUpdate(bot, nr, update)
			}
		case <-stop:
			logrus.Println("Telegram bot stopped.")
			return
		}
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, nr repository.NewsRepository, update tgbotapi.Update) {
	logrus.Info("[" + update.Message.From.UserName + "] " + update.Message.Text)

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "todaynews":
			news, err := nr.GetTodayNews()
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error")
				if _, err := bot.Send(msg); err != nil {
					logrus.Error("error send the msg", err)
				}
			} else {
				groupedNews := groupNewsBySource(news)
				sendGroupedNews(bot, update.Message.Chat.ID, groupedNews)
			}
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know such a command.")
			bot.Send(msg)
		}
	}
}

func groupNewsBySource(news []repository.News) map[string][]repository.News {
	groupedNews := make(map[string][]repository.News)
	for _, n := range news {
		groupedNews[n.SourceName] = append(groupedNews[n.SourceName], n)
	}
	return groupedNews
}

func formatGroupedNews(sourceName string, news []repository.News) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n", sourceName))
	for _, n := range news {
		sb.WriteString(fmt.Sprintf("[%s](%s)\nDescription: %s\n\n",
			n.Title, n.Link, n.Description))
	}
	return sb.String()
}

func sendGroupedNews(bot *tgbotapi.BotAPI, chatID int64, groupedNews map[string][]repository.News) {
	for sourceName, news := range groupedNews {
		formattedNews := formatGroupedNews(sourceName, news)
		sendLongMessage(bot, chatID, formattedNews)
		time.Sleep(1 * time.Second) // Можно убрать или уменьшить задержку
	}
}

func sendLongMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	const maxMessageLength = 4096
	for len(text) > maxMessageLength {
		part := text[:maxMessageLength]
		msg := tgbotapi.NewMessage(chatID, part)
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = true
		bot.Send(msg)
		text = text[maxMessageLength:]
	}
	if len(text) > 0 {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = true
		bot.Send(msg)
	}
}

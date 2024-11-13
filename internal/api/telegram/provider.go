package telegram

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kolosiv/news-aggr/internal/repository"
	"github.com/sirupsen/logrus"
)

func TelegramBot(nr repository.NewsRepository) {
	botToken := os.Getenv("TG_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logrus.Error("error tg roken", err)
		return
	}

	// Установите уровень логирования в Debug, чтобы увидеть все запросы и ответы
	// bot.Debug = true

	logrus.Debug("Authorized on account " + bot.Self.UserName)

	var news []repository.News

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			logrus.Info("[ " + update.Message.From.UserName + "] " + update.Message.Text)

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "todaynews":
					if news, err = nr.GetTodayNews(); err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error")
						if _, err := bot.Send(msg); err != nil {
							logrus.Error("error send th msg", err)
						}
					} else {
						for _, n := range news {
							formattedNews := repository.FormatSingleNews(n)
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, formattedNews)

							if _, err := bot.Send(msg); err != nil {
								logrus.Error("error send th msg", err)
							}
						}
					}
				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know such a command.")
					bot.Send(msg)
				}
			}
		}
	}
}

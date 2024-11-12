package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kolosiv/news-aggr/internal/repository"
)

func TelegramBot(nr repository.NewsRepository) {
	botToken := "7802974180:AAE1Df4osPZeeSpKb1DHTT_SsQB0Opkh7l8"
	// chatID := int64(572778494) // поменять начитку токена и id

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	// Установите уровень логирования в Debug, чтобы увидеть все запросы и ответы
	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	var news []repository.News

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "todaynews":
					if news, err = nr.GetTodayNews(); err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка поиска") //мультиязычность

						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
					} else {
						for _, n := range news {
							formattedNews := repository.FormatSingleNews(n)
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, formattedNews)

							if _, err := bot.Send(msg); err != nil {
								log.Panic(err)
							}
						}
					}
				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не знаю такой команды.") //мультиязычность
					bot.Send(msg)
				}
			}
		}
	}
}

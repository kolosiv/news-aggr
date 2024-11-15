package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kolosiv/news-aggr/internal/api/rest"
	"github.com/kolosiv/news-aggr/internal/api/telegram"
	"github.com/kolosiv/news-aggr/internal/database"
	"github.com/kolosiv/news-aggr/internal/parser"
	"github.com/kolosiv/news-aggr/internal/repository"
	"github.com/kolosiv/news-aggr/pkg/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.InitLogger()

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	dbpool := database.ConnectPostgres()
	defer dbpool.Close()

	nr := repository.CreateNewsRepository(dbpool)
	pc := parser.CreateParserController(nr)
	rc := rest.CreateRestController(nr)

	stopParser := make(chan struct{})
	stopTelegram := make(chan struct{})
	stopRest := make(chan struct{})

	doneParser := make(chan struct{})
	doneTelegram := make(chan struct{})
	doneRest := make(chan struct{})

	go func() {
		pc.MainParser(stopParser)
		close(doneParser)
	}()

	go func() {
		telegram.TelegramBot(nr, stopTelegram)
		close(doneTelegram)
	}()

	go func() {
		rest.Setup(rc, stopRest)
		close(doneRest)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan

	close(stopParser)
	close(stopTelegram)
	close(stopRest)

	<-doneParser
	<-doneTelegram
	<-doneRest

	logrus.Println("All services stopped. Exiting.")
}

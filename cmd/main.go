package main

import (
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

	go pc.MainParser()
	go telegram.TelegramBot(nr)
	go rest.Setup(rc)

	select {}
}

package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kolosiv/news-aggr/internal/api/rest"
	"github.com/kolosiv/news-aggr/internal/api/telegram"
	"github.com/kolosiv/news-aggr/internal/database"
	"github.com/kolosiv/news-aggr/internal/parser"
	"github.com/kolosiv/news-aggr/internal/repository"
	"github.com/robfig/cron/v3"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbpool := database.ConnectPostgres()
	defer dbpool.Close()

	nr := repository.CreateNewsRepository(dbpool)
	pc := parser.CreateParserController(nr)
	rc := rest.CreateRestController(nr)

	c := cron.New()
	_, err = c.AddFunc("41 19 * * *", pc.MainParser)
	if err != nil {
		fmt.Println("Ошибка при добавлении задания:", err)
		return
	}
	c.Start()

	go telegram.TelegramBot(nr)
	go rest.Setup(rc)

	// pc.MainParser()

	select {}
}

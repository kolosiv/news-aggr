package source

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

type Source struct {
	URL       string    `json:"url"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_ad"`
}

func GetSources() []Source {
	data, err := ioutil.ReadFile("C:/Projects/Go/code/news-aggr/sources.json") // поменять начитку пути
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	var sources []Source

	err = json.Unmarshal(data, &sources)
	if err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	return sources
}

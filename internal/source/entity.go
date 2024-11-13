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
	data, err := ioutil.ReadFile("sources.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var sources []Source

	err = json.Unmarshal(data, &sources)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	return sources
}

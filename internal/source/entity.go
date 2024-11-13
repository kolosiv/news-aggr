package source

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
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
		logrus.Error("Error reading file", err)
		return nil
	}

	var sources []Source

	err = json.Unmarshal(data, &sources)
	if err != nil {
		logrus.Error("Error decoding JSON", err)
		return nil
	}

	return sources
}

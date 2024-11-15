package parser

import (
	"sync"
	"time"

	repository "github.com/kolosiv/news-aggr/internal/repository"
	"github.com/kolosiv/news-aggr/internal/source"
	"github.com/sirupsen/logrus"
)

type ParserController interface {
	rssParser(rssURL string, sourceName string, wg *sync.WaitGroup)
	MainParser(stop <-chan struct{})
	parseAndStoreRSS(rssURL string, sourceName string)
}

type parserController struct {
	nr repository.NewsRepository
	sr []source.Source
}

func CreateParserController(nr repository.NewsRepository) ParserController {
	return &parserController{
		nr: nr,
		sr: source.GetSources()}
}

func (pc *parserController) MainParser(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			logrus.Println("Parser stopped.")
			return
		default:
			var wg sync.WaitGroup
			for _, source := range pc.sr {
				wg.Add(1)
				switch source.Type {
				case "1": //RSS
					go pc.rssParser(source.URL, source.Name, &wg)
				}
			}
			wg.Wait()
			time.Sleep(5 * time.Minute) // начитывать из среды
		}
	}
}

package parser

import (
	"sync"

	repository "github.com/kolosiv/news-aggr/internal/repository"
	"github.com/kolosiv/news-aggr/internal/source"
)

type ParserController interface {
	rssParser(rssURL string, sourceName string, wg *sync.WaitGroup)
	MainParser()
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

func (pc *parserController) MainParser() {
	var wg sync.WaitGroup
	for _, source := range pc.sr {
		wg.Add(1)
		switch source.Type {
		case "1": //RSS
			go pc.rssParser(source.URL, source.Name, &wg)
		}
	}
	wg.Wait()
}

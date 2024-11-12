package parser

import (
	repository "github.com/kolosiv/news-aggr/internal/repository"
	"github.com/kolosiv/news-aggr/internal/source"
)

type ParserController interface {
	RssParser(rssURL string)
	MainParser()
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
	for _, source := range pc.sr {
		switch source.Type {
		case "1": //RSS
			pc.RssParser(source.URL)
		}
	}
}

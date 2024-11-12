package parser

import (
	"fmt"
	"time"

	"github.com/kolosiv/news-aggr/internal/repository"
	"github.com/mmcdole/gofeed"
)

func (pc *parserController) RssParser(rssURL string) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		// log.("Ошибка при парсинге RSS ленты: %v", err)
		return
	}

	var (
		news    []repository.News
		article repository.News
	)

	now := time.Now()

	fmt.Println("Заголовки новостей:")
	for _, item := range feed.Items {
		artpubl := *item.PublishedParsed
		if artpubl.Year() == now.Year() && artpubl.Month() == now.Month() && artpubl.Day() == now.Day() {
			article.Title = item.Title
			article.Description = item.Description
			article.Link = item.Link
			article.PubDate = *item.PublishedParsed
			news = append(news, article) //можно потом добавить указатель на слайс с помощью reflect
		}
	}

	if len(news) < 1 {
		return
	}

	if err = pc.nr.CreateNews(news); err != nil {
		return
	}
}

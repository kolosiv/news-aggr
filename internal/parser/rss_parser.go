package parser

import (
	"time"

	"github.com/kolosiv/news-aggr/internal/repository"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

func (pc *parserController) RssParser(rssURL string) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		logrus.Error("Error parsing RSS feed: "+rssURL, err)
		return
	}

	logrus.Debug("RSS feed parsed successfully: " + rssURL)

	var (
		news    []repository.News
		article repository.News
	)

	now := time.Now()

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

	pc.nr.CreateNews(news)
}

package parser

import (
	"sync"
	"time"

	repository "github.com/kolosiv/news-aggr/internal/repository"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

const (
	checkInterval = 10 * time.Minute // начитывать из среды
)

func (pc *parserController) rssParser(rssURL string, sourceName string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		pc.parseAndStoreRSS(rssURL, sourceName)
		time.Sleep(checkInterval)
	}
}

func (pc *parserController) parseAndStoreRSS(rssURL string, sourceName string) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		logrus.Error("Error parsing RSS feed: "+rssURL, err)
		return
	}

	logrus.Info("RSS feed parsed successfully: " + rssURL)

	for _, item := range feed.Items {
		newsItem := repository.News{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			PubDate:     *item.PublishedParsed,
			SourceName:  sourceName,
		}

		exists, err := pc.nr.NewsExists(newsItem)
		if err != nil {
			continue
		}
		if !exists {
			pc.nr.CreateNewsItem(newsItem)
		}
	}

	// var (
	// 	news    []repository.News
	// 	article repository.News
	// )

	// now := time.Now()

	// for _, item := range feed.Items {
	// artpubl := *item.PublishedParsed
	// if artpubl.Year() == now.Year() && artpubl.Month() == now.Month() && artpubl.Day() == now.Day() {
	// article.Title = item.Title
	// article.Description = item.Description
	// article.Link = item.Link
	// article.PubDate = *item.PublishedParsed
	// article.SourceName = sourceName
	// news = append(news, article) //можно потом добавить указатель на слайс с помощью reflect
	// }
	// }

	// if len(news) < 1 {
	// 	logrus.Debug("No news on RSS feed: " + rssURL)
	// 	return
	// }

	// pc.nr.CreateNews(news)
}

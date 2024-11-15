package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type NewsRepository interface {
	CreateNews(news []News) error
	GetTodayNews() ([]News, error)
	GetNewsByInterval(startDate time.Time, endtDate time.Time) ([]News, error)
	NewsExists(news News) (bool, error)
	CreateNewsItem(newsItem News) error
}

type newsRepository struct {
	db *pgxpool.Pool
}

func CreateNewsRepository(db *pgxpool.Pool) NewsRepository {
	return &newsRepository{db: db}
}

func (nr *newsRepository) NewsExists(newsItem News) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM news WHERE link=$1)`
	err := nr.db.QueryRow(context.Background(), query, newsItem.Link).Scan(&exists)
	if err != nil {
		logrus.Error("failed to check if news exists", err)
		return false, err
	}
	return exists, err
}

func (nr *newsRepository) GetNewsByInterval(startDate time.Time, endtDate time.Time) ([]News, error) {
	query := `SELECT id, pub_date, title, description, link, source_name
				FROM news
				WHERE pub_date::date BETWEEN $1 AND $2;`

	rows, err := nr.db.Query(context.Background(), query, startDate, endtDate)
	if err != nil {
		logrus.Error("error select process", err)
		return nil, err
	}
	defer rows.Close()

	var newsList []News

	for rows.Next() {
		var news News
		if err := rows.Scan(&news.ID, &news.PubDate, &news.Title,
			&news.Description, &news.Link, &news.SourceName); err != nil {
			logrus.Error("error select process", err)
			return nil, err
		}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		logrus.Error("error select process", err)
		return nil, err
	}

	return newsList, nil
}

func (nr *newsRepository) CreateNewsItem(newsItem News) error {
	query := `INSERT INTO news (title, link, description, pub_date, source_name) VALUES ($1, $2, $3, $4, $5)`
	_, err := nr.db.Exec(context.Background(), query, newsItem.Title, newsItem.Link,
		newsItem.Description, newsItem.PubDate, newsItem.SourceName)
	if err != nil {
		logrus.Error("failed to store news item", err)
		return err
	} else {
		logrus.Info("New news item stored: " + newsItem.Title)
	}
	return nil
}

func (nr *newsRepository) CreateNews(news []News) error {
	var rows [][]interface{}

	row2ins := len(news)

	for _, article := range news {
		rows = append(rows, []interface{}{article.PubDate, article.Title, article.Description,
			article.Link, article.SourceName})
	}

	copyCount, err := nr.db.CopyFrom(context.Background(),
		pgx.Identifier{"news"},
		[]string{"pub_date", "title", "description", "link", "source_name"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		logrus.Error("failed to insert data with CopyFrom: ", err)
		return err
	}

	logrus.Info("Successfully inserted rows." + strconv.FormatInt(copyCount, 10))

	if row2ins != int(copyCount) {
		logrus.Error("rows missed")
		return nil
	}

	return nil
}

func (nr *newsRepository) GetTodayNews() ([]News, error) {
	query := `SELECT id, pub_date, title, description, link, source_name
				FROM news
				WHERE pub_date::date = CURRENT_DATE`

	rows, err := nr.db.Query(context.Background(), query)
	if err != nil {
		logrus.Error("error select process", err)
		return nil, err
	}
	defer rows.Close()

	var newsList []News

	for rows.Next() {
		var news News
		if err := rows.Scan(&news.ID, &news.PubDate, &news.Title,
			&news.Description, &news.Link, &news.SourceName); err != nil {
			logrus.Error("error scan rows", err)
			return nil, err
		}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		logrus.Error("error row process", err)
		return nil, err
	}

	return newsList, nil
}

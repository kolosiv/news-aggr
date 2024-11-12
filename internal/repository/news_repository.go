package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsRepository interface {
	CreateNews(news []News) error
	GetTodayNews() ([]News, error)
	GetNewsByInterval(startDate time.Time, endtDate time.Time) ([]News, error)
}

type newsRepository struct {
	db *pgxpool.Pool
}

func CreateNewsRepository(db *pgxpool.Pool) NewsRepository {
	return &newsRepository{db: db}
}

func (nr *newsRepository) GetNewsByInterval(startDate time.Time, endtDate time.Time) ([]News, error) {
	query := `SELECT id, pub_date, title, description, link
				FROM news
				WHERE pub_date::date BETWEEN $1 AND $2;`

	rows, err := nr.db.Query(context.Background(), query, startDate, endtDate)
	if err != nil {
		return nil, fmt.Errorf("error select process: %v", err)
	}
	defer rows.Close()

	var newsList []News

	for rows.Next() {
		var news News
		if err := rows.Scan(&news.ID, &news.PubDate, &news.Title, &news.Description, &news.Link); err != nil {
			return nil, fmt.Errorf("error scan rows: %v", err)
		}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error row process: %v", err)
	}

	return newsList, nil
}

func (nr *newsRepository) CreateNews(news []News) error {
	var rows [][]interface{}

	row2ins := len(news)

	for _, article := range news {
		rows = append(rows, []interface{}{article.PubDate, article.Title, article.Description, article.Link})
	}

	copyCount, err := nr.db.CopyFrom(context.Background(),
		pgx.Identifier{"news"},
		[]string{"pub_date", "title", "description", "link"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Fatal("Failed to insert data with CopyFrom: ", err)
	}

	fmt.Printf("Successfully inserted %d rows.\n", copyCount)

	if row2ins != int(copyCount) {
		log.Fatal("Rows missed")
	}

	return nil
}

func (nr *newsRepository) GetTodayNews() ([]News, error) {
	query := `SELECT id, pub_date, title, description, link
				FROM news
				WHERE pub_date::date = CURRENT_DATE`

	rows, err := nr.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error select process: %v", err)
	}
	defer rows.Close()

	var newsList []News

	for rows.Next() {
		var news News
		if err := rows.Scan(&news.ID, &news.PubDate, &news.Title, &news.Description, &news.Link); err != nil {
			return nil, fmt.Errorf("error scan rows: %v", err)
		}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error row process: %v", err)
	}

	return newsList, nil
}

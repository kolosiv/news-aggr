package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func ConnectPostgres() *pgxpool.Pool {

	dbConfig := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		logrus.Fatalf("unable to connect to database: %v", err)
		return nil
	}

	if err := dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		logrus.Fatalf("unable to connect to database: %v", err)
		return nil
	}

	logrus.Info("Successfully connected to PostgreSQL!")

	migrateDir := "/app/internal/database/migrations"
	if _, err := os.Stat(migrateDir); os.IsNotExist(err) {
		logrus.Fatalf("Путь к миграциям не найден: %s", migrateDir)
	}

	migrateURL := fmt.Sprintf("file://%s", migrateDir)
	dbURL := fmt.Sprintf("%s?sslmode=disable", connString)
	m, err := migrate.New(migrateURL, dbURL)
	if err != nil {
		logrus.Fatalf("Ошибка создания миграции: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Ошибка выполнения миграции: %v", err)
	}

	logrus.Info("Миграция успешно выполнена!")

	return dbpool
}

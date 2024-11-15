package database

import (
	"context"
	"fmt"
	"os"
	"time"

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
	return dbpool
}

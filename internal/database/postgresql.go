package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
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

	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
		return nil
	}

	if err := dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		log.Fatalf("unable to connect to database: %v", err)
		return nil
	}

	fmt.Printf("Successfully connected to PostgreSQL!")
	return dbpool
}

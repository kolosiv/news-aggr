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

	logrus.Infof("Try to connect to database: %s", connString)
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

	var migrateDir string
	if os.Getenv("APP_MODE") == "debug" {
		migrateDir = "internal/database/migrations"
	} else {
		migrateDir = os.Getenv("MIGRATE_DIR")
	}
	if _, err := os.Stat(migrateDir); os.IsNotExist(err) {
		logrus.Fatalf("Migration path not found: %s", migrateDir)
	}

	migrateURL := fmt.Sprintf("file://%s", migrateDir)
	dbURL := fmt.Sprintf("%s?sslmode=disable", connString)
	m, err := migrate.New(migrateURL, dbURL)
	if err != nil {
		logrus.Fatalf("Error creating migration: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Error running migration: %v", err)
	}

	logrus.Info("Migration completed successfully!")

	return dbpool
}

package db

import (
	"Solflora/logger"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type credentials struct {
	Username string
	Password string
	Database string
	Host     string
	Port     string
}

var DB *sql.DB

func Init() {
	var log = logger.Logger()
	log.Info("[START] db.init")

	dbCreds := credentials{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Baku",
		dbCreds.Host, dbCreds.Username, dbCreds.Password, dbCreds.Database, dbCreds.Port)
	log.Debugf("[DEBUG] conn-string: %s", dsn)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.WithError(err).Fatal("[ERROR] db.init | Failed to connect to database")
	}

	if err = DB.Ping(); err != nil {
		log.WithError(err).Fatal("[ERROR] db.init | Failed to ping database")
	}

	DB.SetMaxIdleConns(10)
	DB.SetMaxOpenConns(100)
	DB.SetConnMaxLifetime(time.Hour)

	log.Info("[END] db.init")
}

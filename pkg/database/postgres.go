package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rusystem/crm-api/pkg/logger"
	"github.com/tinrab/retry"
	"time"
)

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresConnection(cfg PostgresConfig) (*sql.DB, error) {
	var db *sql.DB
	var err error

	retry.ForeverSleep(time.Second*2, func(i int) error {
		db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.SSLMode, cfg.Password))
		if err != nil {
			logger.Error(fmt.Sprintf("postgres connection, err - %v", err))
			return err
		}

		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(5)

		err = db.Ping()
		if err != nil {
			return err
		}

		return nil
	})

	return db, nil
}

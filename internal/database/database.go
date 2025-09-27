package database

import (
	"PattyWagon/internal/utils"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var (
	DatabaseName    = os.Getenv("DB_DATABASE")
	Password        = os.Getenv("DB_PASSWORD")
	Username        = os.Getenv("DB_USERNAME")
	Port            = os.Getenv("DB_PORT")
	Host            = os.Getenv("DB_HOST")
	Schema          = os.Getenv("DB_SCHEMA")
	MaxOpenConns    = utils.GetEnvInt64("DB_MAX_OPEN_CONNS", 20)
	MaxIdleConns    = utils.GetEnvInt64("DB_MAX_IDLE_CONNS", 10)
	ConnMaxIdleTime = utils.GetEnvInt64("DB_CONN_MAX_IDLE_TIME_IN_SECONDS", 60)
	ConnMaxLifeTime = utils.GetEnvInt64("DB_CONN_MAX_LIFE_TIME_IN_SECONDS", 300)
)

type ConnectionPoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifeTime time.Duration
}

func New(
	host, port string,
	database string,
	username, password string,
	schema string,
	connPoolConfig *ConnectionPoolConfig,
) *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if connPoolConfig != nil {
		db.SetMaxOpenConns(connPoolConfig.MaxOpenConns)
		db.SetMaxIdleConns(connPoolConfig.MaxIdleConns)
		db.SetConnMaxIdleTime(connPoolConfig.ConnMaxIdleTime)
		db.SetConnMaxLifetime(connPoolConfig.ConnMaxLifeTime)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

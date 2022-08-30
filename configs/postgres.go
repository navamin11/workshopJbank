package configs

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func ConnectPostgreSQL() *bun.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error(err)
	}

	// Conncet to PostgreSQL
	DbHost := os.Getenv("DB_HOST")
	DbPort := os.Getenv("DB_PORT")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DbUser, DbPassword, DbHost, DbPort, DbName)

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Error(err)
	}
	config.PreferSimpleProtocol = true

	sqldb := stdlib.OpenDB(*config)
	db := bun.NewDB(sqldb, pgdialect.New())

	// db.AddQueryHook(bundebug.NewQueryHook(
	// 	bundebug.WithVerbose(true),
	// 	bundebug.FromEnv("BUNDEBUG"),
	// ))

	return db
}

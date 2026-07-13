package database

import (
	"context"
	"fmt"
	"log"

	"github.com/Mikhail-Tal63/Orbit/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg configs.Configs) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to PostgreSQL")

	return db
}

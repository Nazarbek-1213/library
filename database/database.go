package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(url string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal("baza sozlamasi xato: ", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal("bazaga ulanib bo'lmadi: ", err)
	}

	log.Println("bazaga ulandi")
	return pool
}

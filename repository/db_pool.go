package repository

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbPool struct {
	Pool *pgxpool.Pool
}

var (
	once   sync.Once
	dbPool *DbPool
)

func GetDbPoolConnection(connString string) *DbPool {
	once.Do(func() {
		pool, err := pgxpool.New(context.Background(), connString)
		if err != nil {
			panic(err)
		}
		dbPool = &DbPool{pool}
	})

	return dbPool
}

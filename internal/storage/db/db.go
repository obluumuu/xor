package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/obluumuu/xor/gen/sqlc"
	"github.com/obluumuu/xor/internal/storage"
)

var _ (storage.Storage) = (*DbStorage)(nil)

type DbStorage struct {
	pool  *pgxpool.Pool
	query *sqlc.Queries
}

func NewDbStorage(ctx context.Context, dbDsn string) (*DbStorage, error) {
	// TODO: add ctx or edit timeouts
	ctxConn, cancelConn := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelConn()
	pool, err := pgxpool.New(ctxConn, dbDsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	ctxPing, cancelPing := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelPing()
	if err := pool.Ping(ctxPing); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	query := sqlc.New(pool)

	return &DbStorage{pool: pool, query: query}, nil
}

func (s *DbStorage) Close() {
	s.pool.Close()
}

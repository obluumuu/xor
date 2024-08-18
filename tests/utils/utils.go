package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	proxy_pb "github.com/obluumuu/xor/gen/proto/proxy"
	"github.com/obluumuu/xor/internal/server"
	"github.com/obluumuu/xor/internal/storage"
	"github.com/obluumuu/xor/internal/storage/db"
	"github.com/obluumuu/xor/internal/storage/inmem"
)

func createPool(pgDsn string) (*pgxpool.Pool, error) {
	ctxConn, cancelConn := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelConn()

	pool, err := pgxpool.New(ctxConn, pgDsn)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	return pool, nil
}

func createConn(pgDsn string) (*pgconn.PgConn, error) {
	ctxConn, cancelConn := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelConn()

	conn, err := pgconn.Connect(ctxConn, pgDsn)
	if err != nil {
		return nil, fmt.Errorf("create conn: %w", err)
	}
	return conn, nil
}

func closeConnect(conn *pgconn.PgConn) error {
	ctxClose, cancelClose := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelClose()

	err := conn.Close(ctxClose)
	if err != nil {
		return fmt.Errorf("close connect: %w", err)
	}
	return nil
}

func setupDb(t *testing.T) (string, func()) {
	t.Helper()

	const appDsn = "postgresql://app:app@localhost:5432/app?sslmode=disable"

	suff := strconv.FormatInt(int64(rand.Int()), 10)
	user := "testuser_" + suff
	db := "testdb_" + suff

	connToApp, err := createConn(appDsn)
	require.NoError(t, err)
	defer func() {
		err := closeConnect(connToApp)
		require.NoError(t, err)
	}()

	_, err = connToApp.Exec(context.Background(), fmt.Sprintf("CREATE USER %s WITH PASSWORD 'test'", user)).ReadAll()
	require.NoError(t, err)
	_, err = connToApp.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", db)).ReadAll()
	require.NoError(t, err)
	_, err = connToApp.Exec(context.Background(), fmt.Sprintf("ALTER DATABASE %s OWNER TO %s", db, user)).ReadAll()
	require.NoError(t, err)

	testDbDsn := fmt.Sprintf("postgresql://%s:test@localhost:5432/%s?sslmode=disable", user, db)

	poolToTest, err := createPool(testDbDsn)
	require.NoError(t, err)

	sqlDb := stdlib.OpenDBFromPool(poolToTest)
	err = goose.Up(sqlDb, "../migrations")
	require.NoError(t, err)

	poolToTest.Close()

	return testDbDsn, func() {
		connToApp, err := createConn(appDsn)
		require.NoError(t, err)
		defer func() {
			err := closeConnect(connToApp)
			require.NoError(t, err)
		}()

		batch := &pgconn.Batch{}
		batch.ExecParams(fmt.Sprintf("DROP DATABASE %s", db), nil, nil, nil, nil)
		batch.ExecParams(fmt.Sprintf("DROP USER %s", user), nil, nil, nil, nil)

		ctxBatch, cancelBatch := context.WithTimeout(context.Background(), time.Second*5)
		defer cancelBatch()
		_, err = connToApp.ExecBatch(ctxBatch, batch).ReadAll()
		require.NoError(t, err)
	}
}

func SetupServer(t *testing.T, opts ...Option) (*server.Server, func()) {
	t.Helper()

	options := &options{}
	for _, opt := range opts {
		opt(options)
	}

	var storage storage.Storage
	var closeStorage func()

	serverType := os.Getenv("TESTING_SERVER_TYPE")
	switch serverType {
	case "PRODUCTION":
		dbDsn, closeDb := setupDb(t)

		var err error
		storage, err = db.NewDbStorage(context.Background(), dbDsn)
		require.NoError(t, err)

		closeStorage = func() {
			storage.Close()
			closeDb()
		}
	case "STANDALONE":
		storage = inmem.NewInmemStorage()
		closeStorage = storage.Close
	default:
		t.Fatalf("env TESTING_SERVER_TYPE must match (PRODUCTION|STANDALONE)")
	}

	srv := server.New(storage)
	s := grpc.NewServer()
	proxy_pb.RegisterProxyServiceServer(s, srv)

	lis := options.lis
	if lis == nil {
		var err error
		lis, err = net.Listen("tcp", "localhost:0")
		require.NoError(t, err)
	}
	log.Printf("listening on %v", lis.Addr())

	go func() {
		err := s.Serve(lis)
		require.True(t, err == nil || errors.Is(err, grpc.ErrServerStopped))
	}()

	return srv, func() {
		s.GracefulStop()
		closeStorage()
	}
}

type options struct {
	lis net.Listener
}

type Option func(*options)

func WithListener(lis net.Listener) Option {
	return func(opts *options) {
		opts.lis = lis
	}
}

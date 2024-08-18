package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	proxy_pb "github.com/obluumuu/xor/gen/proto/proxy"
	"github.com/obluumuu/xor/internal/config"
	"github.com/obluumuu/xor/internal/server"
	"github.com/obluumuu/xor/internal/storage/db"
	"github.com/obluumuu/xor/internal/storage/inmem"
)

func setupLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableLevelTruncation: true,
		PadLevelText:           true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			splitted := strings.Split(f.Function, ".")
			funcname := splitted[len(splitted)-1]
			filename := filepath.Base(f.File)
			return "", fmt.Sprintf(" %s:%d:%s", filename, f.Line, funcname)
		},
	})
	logrus.SetReportCaller(true)
}

func main() {
	setupLogger()

	var (
		addr       = flag.String("addr", "localhost:1234", "host:port to serve")
		configPath = flag.String("config", "./config.json", "Path to config file for standalone mode")
		standalone = flag.Bool("standalone", false, "Run in standalone mode (no db, specify proxies and proxy_blocks in config)")
	)
	flag.Parse()

	if *standalone {
		if err := runStandalone(*addr, *configPath); err != nil {
			logrus.Fatalf("failed to run standalone: %v", err)
		}
	} else {
		if err := runProduction(*addr); err != nil {
			logrus.Fatalf("failed to run production: %v", err)
		}
	}
}

func runStandalone(addr, configPath string) error {
	cfg, err := config.ReadAndParseJsonFile(configPath)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	storage := inmem.NewInmemStorage()
	for _, proxy := range cfg.Proxies {
		if err := storage.CreateProxy(context.Background(), proxy.ToModel()); err != nil {
			return fmt.Errorf("create proxy from config: %w", err)
		}
	}
	for _, proxyBlock := range cfg.ProxyBlocks {
		if err := storage.CreateProxyBlock(context.Background(), proxyBlock.ToModel()); err != nil {
			return fmt.Errorf("create proxy_block from config: %w", err)
		}
	}

	srv := server.New(storage)
	s := grpc.NewServer()
	proxy_pb.RegisterProxyServiceServer(s, srv)
	reflection.Register(s)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	logrus.Infof("listening on %v", addr)

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to ")
	}

	return nil
}

func runProduction(addr string) error {
	dbDsn, ok := os.LookupEnv("PG_DSN")
	if !ok {
		return errors.New("env PG_DSN is not set")
	}

	storage, err := db.NewDbStorage(context.Background(), dbDsn)
	if err != nil {
		return fmt.Errorf("create db storage: %w", err)
	}
	defer storage.Close()

	srv := server.New(storage)
	s := grpc.NewServer()
	proxy_pb.RegisterProxyServiceServer(s, srv)
	reflection.Register(s)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	logrus.Infof("listening on %v", addr)

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

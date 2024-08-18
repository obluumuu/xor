package server

import (
	proxy_pb "github.com/obluumuu/xor/gen/proto/proxy"
	"github.com/obluumuu/xor/internal/storage"
)

type Server struct {
	proxy_pb.UnimplementedProxyServiceServer

	storage storage.Storage
}

func New(storage storage.Storage) *Server {
	return &Server{storage: storage}
}

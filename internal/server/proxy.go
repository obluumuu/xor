package server

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proxy_pb "github.com/obluumuu/xor/gen/proto/proxy"
	"github.com/obluumuu/xor/internal/models"
	"github.com/obluumuu/xor/internal/server/utils"
	"github.com/obluumuu/xor/internal/storage"
)

func (s *Server) CreateProxy(ctx context.Context, req *proxy_pb.CreateProxyRequest) (*proxy_pb.CreateProxyResponse, error) {
	proxy := &models.Proxy{
		Id:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Schema:      utils.SchemaProtoToString(req.Schema),
		Host:        req.Host,
		Port:        int32(req.Port),
		Username:    req.Username,
		Password:    req.Password,
		Tags:        models.NewTagsFromNamesList(req.Tags),
	}
	if err := s.storage.CreateProxy(ctx, proxy); err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error")
	}

	return &proxy_pb.CreateProxyResponse{Id: proxy.Id.String()}, nil
}

func (s *Server) GetProxy(ctx context.Context, req *proxy_pb.GetProxyRequest) (*proxy_pb.GetProxyResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	proxy, err := s.storage.GetProxy(ctx, id)
	if errors.Is(err, storage.ErrProxyNotFound) {
		return nil, status.Errorf(codes.NotFound, "Proxy with id %s not found", req.Id)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	resp := &proxy_pb.GetProxyResponse{
		Id:          proxy.Id.String(),
		Name:        proxy.Name,
		Description: proxy.Description,
		Schema:      utils.SchemaStringToProto(proxy.Schema),
		Host:        proxy.Host,
		Port:        uint32(proxy.Port),
		Username:    proxy.Username,
		Password:    proxy.Password,
	}
	resp.Tags = make([]*proxy_pb.GetProxyResponse_Tag, 0, len(proxy.Tags))
	for _, tag := range proxy.Tags {
		resp.Tags = append(resp.Tags, &proxy_pb.GetProxyResponse_Tag{
			Id:    tag.Id.String(),
			Name:  tag.Name,
			Color: tag.Color,
		})
	}
	return resp, nil
}

func (s *Server) UpdateProxy(ctx context.Context, req *proxy_pb.UpdateProxyRequest) (*proxy_pb.UpdateProxyResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	proxy := &models.Proxy{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		Schema:      utils.SchemaProtoToString(req.Schema),
		Host:        req.Host,
		Port:        int32(req.Port),
		Username:    &req.Username,
		Password:    &req.Password,
		Tags:        models.NewTagsFromNamesList(req.Tags),
	}
	err = s.storage.UpdateProxy(ctx, proxy, req.FieldMask)
	if errors.Is(err, storage.ErrProxyNotFound) {
		return nil, status.Errorf(codes.NotFound, "Proxy with id %s not found", req.Id)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &proxy_pb.UpdateProxyResponse{}, nil
}

func (s *Server) DeleteProxy(ctx context.Context, req *proxy_pb.DeleteProxyRequest) (*proxy_pb.DeleteProxyResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	err = s.storage.DeleteProxy(ctx, id)
	if errors.Is(err, storage.ErrProxyNotFound) {
		return nil, status.Errorf(codes.NotFound, "Proxy with id %s not found", req.Id)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}
	return &proxy_pb.DeleteProxyResponse{}, nil
}

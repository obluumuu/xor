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

func (s *Server) CreateProxyBlock(ctx context.Context, req *proxy_pb.CreateProxyBlockRequest) (*proxy_pb.CreateProxyBlockResponse, error) {
	proxyBlock := &models.ProxyBlock{
		Id:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Tags:        models.NewTagsFromNamesList(req.Tags),
	}

	if err := s.storage.CreateProxyBlock(ctx, proxyBlock); err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error")
	}

	return &proxy_pb.CreateProxyBlockResponse{Id: proxyBlock.Id.String()}, nil
}

func (s *Server) GetProxyBlock(ctx context.Context, req *proxy_pb.GetProxyBlockRequest) (*proxy_pb.GetProxyBlockResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	proxyBlock, err := s.storage.GetProxyBlock(ctx, id)
	if errors.Is(err, storage.ErrProxyBlockNotFound) {
		return nil, status.Errorf(codes.NotFound, "ProxyBlock with id %s not found", req.Id)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	resp := &proxy_pb.GetProxyBlockResponse{
		Id:          proxyBlock.Id.String(),
		Name:        proxyBlock.Name,
		Description: proxyBlock.Description,
	}
	resp.Tags = make([]*proxy_pb.GetProxyBlockResponse_Tag, 0, len(resp.Tags))
	for _, tag := range proxyBlock.Tags {
		resp.Tags = append(resp.Tags, &proxy_pb.GetProxyBlockResponse_Tag{
			Id:    tag.Id.String(),
			Name:  tag.Name,
			Color: tag.Color,
		})
	}

	return resp, nil
}

func (s *Server) UpdateProxyBlock(ctx context.Context, req *proxy_pb.UpdateProxyBlockRequest) (*proxy_pb.UpdateProxyBlockResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	proxyBlock := &models.ProxyBlock{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		Tags:        models.NewTagsFromNamesList(req.Tags),
	}
	err = s.storage.UpdateProxyBlock(ctx, proxyBlock, req.FieldMask)
	if errors.Is(err, storage.ErrProxyBlockNotFound) {
		return nil, status.Errorf(codes.NotFound, "ProxyBlock with id %s not found", req.Id)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &proxy_pb.UpdateProxyBlockResponse{}, nil
}

func (s *Server) DeleteProxyBlock(ctx context.Context, req *proxy_pb.DeleteProxyBlockRequest) (*proxy_pb.DeleteProxyBlockResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	err = s.storage.DeleteProxyBlock(ctx, id)
	if errors.Is(err, storage.ErrProxyBlockNotFound) {
		return nil, status.Errorf(codes.NotFound, "ProxyBlock with id %s not found", req.Id)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &proxy_pb.DeleteProxyBlockResponse{}, nil
}

func (s *Server) GetProxiesByProxyBlockId(ctx context.Context, req *proxy_pb.GetProxiesByProxyBlockIdRequest) (*proxy_pb.GetProxiesByProxyBlockIdResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	proxies, err := s.storage.GetProxiesByProxyBlockId(ctx, id)
	if errors.Is(err, storage.ErrProxyBlockNotFound) {
		return nil, status.Errorf(codes.NotFound, "proxy_block with id %s not found", req.Id)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	res := &proxy_pb.GetProxiesByProxyBlockIdResponse{Proxies: make([]*proxy_pb.GetProxiesByProxyBlockIdResponse_Proxy, 0, len(proxies))}
	for _, proxy := range proxies {
		res.Proxies = append(res.Proxies, &proxy_pb.GetProxiesByProxyBlockIdResponse_Proxy{
			Id:          proxy.Id.String(),
			Name:        proxy.Name,
			Description: proxy.Description,
			Schema:      utils.SchemaStringToProto(proxy.Schema),
			Host:        proxy.Host,
			Port:        uint32(proxy.Port),
			Username:    proxy.Username,
			Password:    proxy.Password,
		})
	}

	return res, nil
}

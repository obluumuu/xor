package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/obluumuu/xor/internal/models"
)

var (
	ErrProxyNotFound      = errors.New("proxy not found")
	ErrProxyBlockNotFound = errors.New("proxy_block not found")
)

type Storage interface {
	CreateProxy(ctx context.Context, proxy *models.Proxy) error
	GetProxy(ctx context.Context, id uuid.UUID) (*models.Proxy, error)
	UpdateProxy(ctx context.Context, proxy *models.Proxy, fieldMask []string) error
	DeleteProxy(ctx context.Context, id uuid.UUID) error

	CreateProxyBlock(ctx context.Context, proxyBlock *models.ProxyBlock) error
	GetProxyBlock(ctx context.Context, id uuid.UUID) (*models.ProxyBlock, error)
	UpdateProxyBlock(ctx context.Context, proxyBlock *models.ProxyBlock, fieldMask []string) error
	DeleteProxyBlock(ctx context.Context, id uuid.UUID) error

	GetProxiesByProxyBlockId(ctx context.Context, id uuid.UUID) ([]*models.Proxy, error)

	Close()
}

package inmem

import (
	"sync"

	"github.com/google/uuid"

	"github.com/obluumuu/xor/internal/models"
	"github.com/obluumuu/xor/internal/storage"
)

var _ (storage.Storage) = (*InmemStorage)(nil)

type InmemStorage struct {
	proxyMu      sync.RWMutex
	proxyBlockMu sync.RWMutex
	tagMu        sync.RWMutex

	proxies     map[uuid.UUID]models.Proxy
	proxyBlocks map[uuid.UUID]models.ProxyBlock
	tags        map[uuid.UUID]models.Tag
	tagNameToId map[string]uuid.UUID

	proxyToTag map[uuid.UUID]map[uuid.UUID]struct{}
	tagToProxy map[uuid.UUID]map[uuid.UUID]struct{}

	proxyBlockToTag map[uuid.UUID]map[uuid.UUID]struct{}
	tagToProxyBlock map[uuid.UUID]map[uuid.UUID]struct{}
}

func NewInmemStorage() *InmemStorage {
	return &InmemStorage{
		proxies:     make(map[uuid.UUID]models.Proxy),
		proxyBlocks: make(map[uuid.UUID]models.ProxyBlock),
		tags:        make(map[uuid.UUID]models.Tag),
		tagNameToId: make(map[string]uuid.UUID),

		proxyToTag: make(map[uuid.UUID]map[uuid.UUID]struct{}),
		tagToProxy: make(map[uuid.UUID]map[uuid.UUID]struct{}),

		proxyBlockToTag: make(map[uuid.UUID]map[uuid.UUID]struct{}),
		tagToProxyBlock: make(map[uuid.UUID]map[uuid.UUID]struct{}),
	}
}

func (s *InmemStorage) Close() {
}

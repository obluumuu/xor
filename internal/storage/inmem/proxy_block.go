package inmem

import (
	"context"
	"slices"

	"github.com/google/uuid"

	"github.com/obluumuu/xor/internal/models"
	"github.com/obluumuu/xor/internal/storage"
)

func (s *InmemStorage) CreateProxyBlock(ctx context.Context, proxyBlock *models.ProxyBlock) error {
	s.proxyBlockMu.Lock()
	defer s.proxyBlockMu.Unlock()

	s.tagMu.Lock()
	defer s.tagMu.Unlock()

	tags := proxyBlock.Tags
	proxyBlock.Tags = nil

	s.proxyBlocks[proxyBlock.Id] = *proxyBlock

	if len(tags) > 0 {
		s.proxyBlockToTag[proxyBlock.Id] = make(map[uuid.UUID]struct{})
	}

	for _, tag := range tags {
		if _, ok := s.tagNameToId[tag.Name]; !ok {
			s.tagNameToId[tag.Name] = tag.Id
			s.tags[tag.Id] = tag
		} else {
			tag = s.tags[s.tagNameToId[tag.Name]]
		}

		if s.tagToProxyBlock[tag.Id] == nil {
			s.tagToProxyBlock[tag.Id] = make(map[uuid.UUID]struct{})
		}

		s.tagToProxyBlock[tag.Id][proxyBlock.Id] = struct{}{}
		s.proxyBlockToTag[proxyBlock.Id][tag.Id] = struct{}{}
	}
	proxyBlock.Tags = nil

	return nil
}

func (s *InmemStorage) GetProxyBlock(ctx context.Context, id uuid.UUID) (*models.ProxyBlock, error) {
	s.proxyBlockMu.RLock()
	defer s.proxyBlockMu.RUnlock()

	s.tagMu.RLock()
	defer s.tagMu.RUnlock()

	proxyBlock, ok := s.proxyBlocks[id]
	if !ok {
		return nil, storage.ErrProxyBlockNotFound
	}

	if len(s.proxyBlockToTag[proxyBlock.Id]) > 0 {
		proxyBlock.Tags = make([]models.Tag, 0, len(s.proxyBlockToTag[proxyBlock.Id]))
	}

	for tagId := range s.proxyBlockToTag[proxyBlock.Id] {
		proxyBlock.Tags = append(proxyBlock.Tags, s.tags[tagId])
	}

	return &proxyBlock, nil
}

func (s *InmemStorage) UpdateProxyBlock(ctx context.Context, proxyBlock *models.ProxyBlock, fieldMask []string) error {
	s.proxyBlockMu.Lock()
	defer s.proxyBlockMu.Unlock()

	s.tagMu.Lock()
	defer s.tagMu.Unlock()

	curProxyBlock, ok := s.proxyBlocks[proxyBlock.Id]
	if !ok {
		return storage.ErrProxyBlockNotFound
	}

	if slices.Contains(fieldMask, "name") {
		curProxyBlock.Name = proxyBlock.Name
	}
	if slices.Contains(fieldMask, "description") {
		curProxyBlock.Description = proxyBlock.Description
	}

	s.proxyBlocks[curProxyBlock.Id] = curProxyBlock

	if slices.Contains(fieldMask, "tags") {
		for tagId := range s.proxyBlockToTag[curProxyBlock.Id] {
			delete(s.tagToProxyBlock[tagId], curProxyBlock.Id)

			if len(s.tagToProxyBlock[tagId]) == 0 {
				delete(s.tagToProxyBlock, tagId)
			}
		}
		clear(s.proxyBlockToTag[curProxyBlock.Id])

		for _, tag := range proxyBlock.Tags {
			if _, ok := s.tagNameToId[tag.Name]; !ok {
				s.tagNameToId[tag.Name] = tag.Id
				s.tags[tag.Id] = tag
			} else {
				tag = s.tags[s.tagNameToId[tag.Name]]
			}

			if s.tagToProxyBlock[tag.Id] == nil {
				s.tagToProxyBlock[tag.Id] = make(map[uuid.UUID]struct{})
			}

			s.tagToProxyBlock[tag.Id][curProxyBlock.Id] = struct{}{}
			s.proxyBlockToTag[curProxyBlock.Id][tag.Id] = struct{}{}
		}
	}

	return nil
}

func (s *InmemStorage) DeleteProxyBlock(ctx context.Context, id uuid.UUID) error {
	s.proxyBlockMu.Lock()
	defer s.proxyBlockMu.Unlock()

	s.tagMu.Lock()
	defer s.tagMu.Unlock()

	proxyBlock, ok := s.proxyBlocks[id]
	if !ok {
		return storage.ErrProxyBlockNotFound
	}

	for _, tag := range proxyBlock.Tags {
		delete(s.tagToProxyBlock[tag.Id], proxyBlock.Id)

		if len(s.tagToProxyBlock[tag.Id]) == 0 {
			delete(s.tagToProxyBlock, tag.Id)
		}
	}
	s.proxyBlockToTag[proxyBlock.Id] = nil
	delete(s.proxyBlocks, id)

	return nil
}

func (s *InmemStorage) GetProxiesByProxyBlockId(ctx context.Context, id uuid.UUID) ([]*models.Proxy, error) {
	s.proxyMu.RLock()
	defer s.proxyMu.RUnlock()

	s.proxyBlockMu.RLock()
	defer s.proxyBlockMu.RUnlock()

	s.tagMu.RLock()
	defer s.tagMu.RUnlock()

	proxyBlock, ok := s.proxyBlocks[id]
	if !ok {
		return nil, storage.ErrProxyBlockNotFound
	}

	var matchedProxies map[uuid.UUID]struct{}

	for tagId := range s.proxyBlockToTag[proxyBlock.Id] {
		if matchedProxies == nil {
			matchedProxies = make(map[uuid.UUID]struct{})
			for proxyId := range s.tagToProxy[tagId] {
				matchedProxies[proxyId] = struct{}{}
			}
			continue
		}

		updatedMatchedProxies := make(map[uuid.UUID]struct{})
		for proxyId := range s.tagToProxy[tagId] {
			if _, ok := matchedProxies[proxyId]; ok {
				updatedMatchedProxies[proxyId] = struct{}{}
			}
		}
		matchedProxies = updatedMatchedProxies
	}

	res := make([]*models.Proxy, 0, len(matchedProxies))
	for proxyId := range matchedProxies {
		proxy := s.proxies[proxyId]
		res = append(res, &proxy)
	}
	return res, nil
}

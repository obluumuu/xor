package inmem

import (
	"context"
	"slices"

	"github.com/google/uuid"

	"github.com/obluumuu/xor/internal/models"
	"github.com/obluumuu/xor/internal/storage"
)

func (s *InmemStorage) CreateProxy(ctx context.Context, proxy *models.Proxy) error {
	s.proxyMu.Lock()
	defer s.proxyMu.Unlock()

	s.tagMu.Lock()
	defer s.tagMu.Unlock()

	tags := proxy.Tags
	proxy.Tags = nil

	s.proxies[proxy.Id] = *proxy

	if len(tags) > 0 {
		s.proxyToTag[proxy.Id] = make(map[uuid.UUID]struct{})
	}

	for _, tag := range tags {
		if _, ok := s.tagNameToId[tag.Name]; !ok {
			s.tags[tag.Id] = tag
			s.tagNameToId[tag.Name] = tag.Id
		} else {
			tag = s.tags[s.tagNameToId[tag.Name]]
		}

		if s.tagToProxy[tag.Id] == nil {
			s.tagToProxy[tag.Id] = make(map[uuid.UUID]struct{})
		}

		s.tagToProxy[tag.Id][proxy.Id] = struct{}{}
		s.proxyToTag[proxy.Id][tag.Id] = struct{}{}
	}

	return nil
}

func (s *InmemStorage) GetProxy(ctx context.Context, id uuid.UUID) (*models.Proxy, error) {
	s.proxyMu.RLock()
	defer s.proxyMu.RUnlock()

	s.tagMu.RLock()
	defer s.tagMu.RUnlock()

	proxy, ok := s.proxies[id]
	if !ok {
		return nil, storage.ErrProxyNotFound
	}

	if len(s.proxyToTag[proxy.Id]) > 0 {
		proxy.Tags = make([]models.Tag, 0, len(s.proxyToTag[proxy.Id]))
	}

	for tagId := range s.proxyToTag[proxy.Id] {
		proxy.Tags = append(proxy.Tags, s.tags[tagId])
	}

	return &proxy, nil
}

func (s *InmemStorage) UpdateProxy(ctx context.Context, proxy *models.Proxy, fieldMask []string) error {
	s.proxyMu.Lock()
	defer s.proxyMu.Unlock()

	s.tagMu.Lock()
	defer s.tagMu.Unlock()

	curProxy, ok := s.proxies[proxy.Id]
	if !ok {
		return storage.ErrProxyNotFound
	}

	if slices.Contains(fieldMask, "name") {
		curProxy.Name = proxy.Name
	}
	if slices.Contains(fieldMask, "description") {
		curProxy.Description = proxy.Description
	}
	if slices.Contains(fieldMask, "schema") {
		curProxy.Schema = proxy.Schema
	}
	if slices.Contains(fieldMask, "host") {
		curProxy.Host = proxy.Host
	}
	if slices.Contains(fieldMask, "port") {
		curProxy.Port = proxy.Port
	}
	if slices.Contains(fieldMask, "username") {
		curProxy.Username = proxy.Username
	}
	if slices.Contains(fieldMask, "password") {
		curProxy.Password = proxy.Password
	}

	s.proxies[curProxy.Id] = curProxy

	if slices.Contains(fieldMask, "tags") {
		for tagId := range s.proxyToTag[proxy.Id] {
			delete(s.tagToProxy[tagId], curProxy.Id)

			if len(s.tagToProxy[tagId]) == 0 {
				delete(s.tagToProxy, tagId)
			}
		}
		clear(s.proxyToTag[proxy.Id])

		for _, tag := range proxy.Tags {
			if _, ok := s.tagNameToId[tag.Name]; !ok {
				s.tagNameToId[tag.Name] = tag.Id
				s.tags[tag.Id] = tag
			} else {
				tag = s.tags[s.tagNameToId[tag.Name]]
			}

			if s.tagToProxy[tag.Id] == nil {
				s.tagToProxy[tag.Id] = make(map[uuid.UUID]struct{})
			}

			s.tagToProxy[tag.Id][curProxy.Id] = struct{}{}
			s.proxyToTag[curProxy.Id][tag.Id] = struct{}{}
		}
	}

	return nil
}

func (s *InmemStorage) DeleteProxy(ctx context.Context, id uuid.UUID) error {
	s.proxyMu.Lock()
	defer s.proxyMu.Unlock()

	s.tagMu.Lock()
	defer s.tagMu.Unlock()

	proxy, ok := s.proxies[id]
	if !ok {
		return storage.ErrProxyNotFound
	}

	for _, tag := range proxy.Tags {
		delete(s.tagToProxy[tag.Id], proxy.Id)

		if len(s.tagToProxy[tag.Id]) == 0 {
			delete(s.tagToProxy, tag.Id)
		}
	}
	s.proxyToTag[proxy.Id] = nil
	delete(s.proxies, id)

	return nil
}

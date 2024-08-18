package models

import (
	"github.com/google/uuid"
	"github.com/obluumuu/xor/gen/sqlc"
)

type ProxyBlock struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []Tag     `json:"tags"`
}

func NewFromSqlcProxyBlock(proxy *sqlc.ProxyBlock) *ProxyBlock {
	return &ProxyBlock{
		Id:          proxy.Id,
		Name:        proxy.Name,
		Description: proxy.Description,
	}
}

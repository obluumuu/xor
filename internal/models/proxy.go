package models

import (
	"github.com/google/uuid"

	"github.com/obluumuu/xor/gen/sqlc"
)

type Proxy struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Schema      string    `json:"schema"`
	Host        string    `json:"host"`
	Port        int32     `json:"port"`
	Username    *string   `json:"username"`
	Password    *string   `json:"password"`
	Tags        []Tag     `json:"tags"`
}

func NewFromSqlcProxy(proxy *sqlc.Proxy) *Proxy {
	return &Proxy{
		Id:          proxy.Id,
		Name:        proxy.Name,
		Description: proxy.Description,
		Schema:      proxy.Schema,
		Host:        proxy.Host,
		Port:        proxy.Port,
		Username:    proxy.Username.Ptr(),
		Password:    proxy.Password.Ptr(),
	}
}

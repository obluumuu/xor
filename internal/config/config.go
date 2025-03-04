package config

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/obluumuu/xor/internal/models"
)

type Proxy struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Schema      string    `json:"schema"`
	Host        string    `json:"host"`
	Port        uint16    `json:"port"`
	Username    *string   `json:"username"`
	Password    *string   `json:"password"`
	Tags        []string  `json:"tags"`
}

func (p *Proxy) ToModel() *models.Proxy {
	return &models.Proxy{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Schema:      p.Schema,
		Host:        p.Host,
		Port:        int32(p.Port),
		Username:    p.Username,
		Password:    p.Password,
		Tags:        models.NewTagsFromNamesList(p.Tags),
	}
}

type ProxyBlock struct {
	Id          uuid.UUID `koanf:"id"`
	Name        string    `koanf:"name"`
	Description string    `koanf:"description"`
	Tags        []string  `koanf:"tags"`
}

func (p *ProxyBlock) ToModel() *models.ProxyBlock {
	return &models.ProxyBlock{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Tags:        models.NewTagsFromNamesList(p.Tags),
	}
}

type Config struct {
	Proxies     []Proxy      `koanf:"proxies"`
	ProxyBlocks []ProxyBlock `koanf:"proxy_blocks"`
}

func ReadAndParseJsonFile(filename string) (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(filename), json.Parser()); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

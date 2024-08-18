package config

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/obluumuu/xor/internal/models"
)

type Proxy struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Schema      string   `json:"schema"`
	Host        string   `json:"host"`
	Port        uint16   `json:"port"`
	Username    *string  `json:"username"`
	Password    *string  `json:"password"`
	Tags        []string `json:"tags"`
}

func (p *Proxy) ToModel() *models.Proxy {
	proxy := &models.Proxy{
		Id:          uuid.New(),
		Name:        p.Name,
		Description: p.Description,
		Schema:      p.Schema,
		Host:        p.Host,
		Port:        int32(p.Port),
		Username:    p.Username,
		Password:    p.Password,
	}
	proxy.Tags = models.NewTagsFromNamesList(p.Tags)

	return proxy
}

type ProxyBlock struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func (p *ProxyBlock) ToModel() *models.ProxyBlock {
	proxyBlock := &models.ProxyBlock{
		Id:          uuid.New(),
		Name:        p.Name,
		Description: p.Description,
	}
	proxyBlock.Tags = models.NewTagsFromNamesList(p.Tags)

	return proxyBlock
}

type Config struct {
	Proxies     []Proxy      `json:"proxies"`
	ProxyBlocks []ProxyBlock `json:"proxy_blocks"`
}

func ReadAndParseJsonFile(filename string) (*Config, error) {
	viper.SetConfigFile(filename)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

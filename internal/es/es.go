package es

import (
	"github.com/elastic/go-elasticsearch/v8"
	"nba-remake/internal/config"
)

// 根据conf的配置返回es的客户端
func NewEsClient(conf *config.ElasticsearchConfig) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: conf.Addresses,
		Username:  conf.Username,
		Password:  conf.Password,
	}
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return esClient
}

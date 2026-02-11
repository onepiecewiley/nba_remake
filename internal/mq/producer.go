package mq

import (
	"log"

	"github.com/IBM/sarama"
	"nba-remake/internal/config"
)

type Producer struct {
	client sarama.SyncProducer
	topic  string // 存一下 topic，发送时就不用每次传了
}

// NewProducer 现在接收 config.KafkaConfig
func NewProducer(conf config.KafkaConfig) (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Partitioner = sarama.NewHashPartitioner

	// 使用配置里的 Brokers
	client, err := sarama.NewSyncProducer(conf.Brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &Producer{
		client: client,
		topic:  conf.Topic, // 从配置里拿 topic
	}, nil
}

// Send 现在只需要 key 和 value，topic 自动用配置里的
func (p *Producer) Send(key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic, // 使用配置好的 Topic
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := p.client.SendMessage(msg)
	if err != nil {
		log.Printf("Kafka Send Error: %v", err)
		return err
	}
	return nil
}

func (p *Producer) Close() {
	p.client.Close()
}

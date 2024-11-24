// Package kafka implementation for producing events
package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	CreateCompanyEvent = "create_company"
	DeleteCompanyEvent = "delete_company"
	UpdateCompanyEvent = "update_company"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(cfg *Config) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}

	return &Producer{
		writer: writer,
		topic:  cfg.Topic,
	}
}

func (p *Producer) Publish(ctx context.Context, key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}
	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("Failed to publish message to Kafka: %v", err)
	}
	return err
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

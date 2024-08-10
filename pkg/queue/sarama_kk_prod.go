package queue

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

// Producer is the interface that wraps the basic Produce method
type Producer interface {
	SendMessage(ctx context.Context, topic string, key string, message []byte) error
	Close() error
}

type SaramaKafkaProducer struct {
	producer  sarama.AsyncProducer
	successes chan *sarama.ProducerMessage
	errors    chan *sarama.ProducerError
}

func NewSaramaKafkaProducer(brokers []string) (*SaramaKafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	config.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	sp := SaramaKafkaProducer{
		producer:  producer,
		successes: make(chan *sarama.ProducerMessage),
		errors:    make(chan *sarama.ProducerError),
	}

	go sp.handleSuccesses()
	go sp.handleErrors()

	return &sp, nil
}

func (sp *SaramaKafkaProducer) SendMessage(ctx context.Context, topic string, key string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
		Key:   sarama.StringEncoder(key),
	}

	select {
	case sp.producer.Input() <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (sp *SaramaKafkaProducer) handleSuccesses() {
	for msg := range sp.producer.Successes() {
		log.Printf("Sent message to topic %s partition %d and offset %d", msg.Topic, msg.Partition, msg.Offset)
		sp.successes <- msg
	}
}

func (sp *SaramaKafkaProducer) handleErrors() {
	for err := range sp.producer.Errors() {
		log.Printf("error sending message %v", err.Err)
		sp.errors <- err
	}
}

func (sp *SaramaKafkaProducer) Close() error {
	if err := sp.producer.Close(); err != nil {
		return err
	}
	close(sp.successes)
	close(sp.errors)
	return nil
}

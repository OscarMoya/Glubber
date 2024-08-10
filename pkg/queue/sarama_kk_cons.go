package queue

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type SaramaKafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	messages      chan *sarama.ConsumerMessage
	errors        chan error
}

// NewSaramaKafkaConsumer crea una nueva instancia del consumidor Kafka usando Sarama.
func NewSaramaKafkaConsumer(brokers []string, groupID string) (*SaramaKafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategySticky()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	sc := &SaramaKafkaConsumer{
		consumerGroup: consumerGroup,
		messages:      make(chan *sarama.ConsumerMessage),
		errors:        make(chan error),
	}

	go sc.handleErrors()

	return sc, nil
}

// ConsumeMessages inicia la lectura de mensajes desde un tópico de forma asincrónica.
func (sc *SaramaKafkaConsumer) ConsumeMessages(ctx context.Context, topics []string) error {
	handler := &consumerGroupHandler{
		messages: sc.messages,
		errors:   sc.errors,
	}

	// Consumir mensajes en una goroutine
	go func() {
		for {
			if err := sc.consumerGroup.Consume(ctx, topics, handler); err != nil {
				sc.errors <- err
			}

			// Si el contexto se cancela, se sale del bucle
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

// handleErrors maneja los errores ocurridos durante la operación del consumidor.
func (sc *SaramaKafkaConsumer) handleErrors() {
	for err := range sc.consumerGroup.Errors() {
		log.Printf("Error en el consumidor: %v", err)
		sc.errors <- err
	}
}

// Close cierra el consumidor Kafka de forma segura.
func (sc *SaramaKafkaConsumer) Close() error {
	if err := sc.consumerGroup.Close(); err != nil {
		return err
	}
	close(sc.messages)
	close(sc.errors)
	return nil
}

// consumerGroupHandler implementa sarama.ConsumerGroupHandler.
type consumerGroupHandler struct {
	messages chan<- *sarama.ConsumerMessage
	errors   chan<- error
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.messages <- msg
		session.MarkMessage(msg, "")
	}
	return nil
}

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

func (sc *SaramaKafkaConsumer) ConsumeMessages(ctx context.Context, topics []string) error {
	handler := &consumerGroupHandler{
		messages: sc.messages,
		errors:   sc.errors,
	}

	go func() {
		for {
			if err := sc.consumerGroup.Consume(ctx, topics, handler); err != nil {
				sc.errors <- err
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

func (sc *SaramaKafkaConsumer) handleErrors() {
	for err := range sc.consumerGroup.Errors() {
		log.Printf("consumer error: %v", err)
		sc.errors <- err
	}
}

func (sc *SaramaKafkaConsumer) Close() error {
	if err := sc.consumerGroup.Close(); err != nil {
		return err
	}
	close(sc.messages)
	close(sc.errors)
	return nil
}

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

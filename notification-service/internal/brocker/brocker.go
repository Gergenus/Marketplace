package brocker

import (
	"fmt"

	"github.com/IBM/sarama"
)

type KafkaBrocker struct {
	EventConsumer sarama.Consumer
}

type BrockerInterface interface {
	RecieveMessages(topic string) (sarama.PartitionConsumer, error)
}

func NewKafkaBrocker(EventConsumer sarama.Consumer) KafkaBrocker {
	return KafkaBrocker{EventConsumer: EventConsumer}
}

func (k KafkaBrocker) RecieveMessages(topic string) (sarama.PartitionConsumer, error) {
	const op = "brocker.RecieveMessages"
	consumer, err := k.EventConsumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return consumer, nil
}

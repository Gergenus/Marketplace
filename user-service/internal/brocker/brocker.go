package brocker

import (
	"fmt"

	"github.com/IBM/sarama"
)

type KafkaBrocker struct {
	EventProducer sarama.SyncProducer
}

type BrockerInterface interface {
	// returns partition, offset and error
	SendMailOrder(topic string, message []byte) (int32, int64, error)
}

// Close the produces connection
func NewKafkaBrocker(EventProducer sarama.SyncProducer) KafkaBrocker {
	return KafkaBrocker{EventProducer: EventProducer}
}

func (k KafkaBrocker) SendMailOrder(topic string, message []byte) (int32, int64, error) {
	const op = "brocker.SendMailOrder"
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	part, offset, err := k.EventProducer.SendMessage(msg)
	if err != nil {
		return -1, -1, fmt.Errorf("%s: %w", op, err)
	}
	return part, offset, nil
}

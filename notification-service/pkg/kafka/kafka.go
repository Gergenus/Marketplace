package kafka

import (
	"github.com/IBM/sarama"
)

func ConnectConsumer(brokers []string) sarama.Consumer {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	cons, err := sarama.NewConsumer(brokers, cfg)
	if err != nil {
		panic(err)
	}
	return cons
}

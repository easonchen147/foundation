package kafka

import (
	"errors"
	"fmt"

	"github.com/easonchen147/foundation/cfg"

	"github.com/segmentio/kafka-go"
)

var (
	producers map[string]*kafka.Writer
	consumers map[string]*kafka.Reader
)

func init() {
	err := InitProducer(cfg.AppConf)
	if err != nil {
		panic(fmt.Sprintf("init kafka producer failed: %s", err))
	}

	err = InitConsumer(cfg.AppConf)
	if err != nil {
		panic(fmt.Sprintf("init kafka consumer failed: %s", err))
	}
}

func InitProducer(cfg *cfg.AppConfig) error {
	if cfg.KafkaConfig == nil {
		return nil
	}
	producers = make(map[string]*kafka.Writer)
	for name, kafkaCfg := range cfg.KafkaConfig.Producers {
		producers[name] = &kafka.Writer{
			Addr:     kafka.TCP(kafkaCfg.Broker),
			Topic:    kafkaCfg.Topic,
			Balancer: &kafka.LeastBytes{},
		}
	}
	return nil
}

func InitConsumer(cfg *cfg.AppConfig) error {
	if cfg.KafkaConfig == nil {
		return nil
	}
	for name, kafkaCfg := range cfg.KafkaConfig.Consumers {
		consumers[name] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{kafkaCfg.Broker},
			GroupID:   kafkaCfg.Group,
			Topic:     kafkaCfg.Topic,
			Partition: kafkaCfg.Partition,
		})
	}
	return nil
}

func Producer(name string) *kafka.Writer {
	if len(producers) == 0 {
		panic(errors.New("kafka producer is not ready"))
	}
	return producers[name]
}

func Writer(name string) *kafka.Reader {
	if len(consumers) == 0 {
		panic(errors.New("kafka consumer is not ready"))
	}
	return consumers[name]
}

func Close() {
	for _, producer := range producers {
		_ = producer.Close()
	}
	for _, consumer := range consumers {
		_ = consumer.Close()
	}
}

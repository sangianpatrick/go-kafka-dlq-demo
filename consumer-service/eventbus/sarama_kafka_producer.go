package eventbus

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

// SaramaKafkaProducerAdapterConfig is a configuration of sarama kafka adapter.
type SaramaKafkaProducerAdapterConfig struct {
	AsyncProducer sarama.AsyncProducer
}

// SaramaKafkaProducerAdapter is a concrete struct of sarma kafka adapter.
type SaramaKafkaProducerAdapter struct {
	logger *logrus.Logger
	config *SaramaKafkaProducerAdapterConfig
}

// NewSaramaKafkaProducerAdapter is a constructor.
func NewSaramaKafkaProducerAdapter(logger *logrus.Logger, config *SaramaKafkaProducerAdapterConfig) Publisher {
	p := &SaramaKafkaProducerAdapter{
		logger, config,
	}
	go p.run()

	return p
}

func (skpa *SaramaKafkaProducerAdapter) run() {
	for producerError := range skpa.config.AsyncProducer.Errors() {
		skpa.logger.Errorf("[Sarama] %s", producerError.Error())
	}
}

// Send will send the message to the brokers
func (skpa *SaramaKafkaProducerAdapter) Send(ctx context.Context, topic string, key string, headers MessageHeaders, message []byte) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			skpa.logger.Errorf("[Sarama] %v", r)
		}
	}()
	bunchOfRecordHeaders := make([]sarama.RecordHeader, 0)
	if headers != nil {
		for key, value := range headers {
			recordHeader := sarama.RecordHeader{
				Key: []byte(key), Value: []byte(value),
			}
			bunchOfRecordHeaders = append(bunchOfRecordHeaders, recordHeader)
		}
	}

	producerMessage := &sarama.ProducerMessage{
		Headers: bunchOfRecordHeaders,
		Key:     sarama.ByteEncoder(key),
		Topic:   topic,
		Value:   sarama.ByteEncoder(message),
	}

	channel := skpa.config.AsyncProducer.Input()
	channel <- producerMessage

	return
}

// Close will stop the producer
func (skpa *SaramaKafkaProducerAdapter) Close() (err error) {
	err = skpa.config.AsyncProducer.Close()
	skpa.logger.Info("[Sarama] Producer is gracefully shutdown")
	return
}

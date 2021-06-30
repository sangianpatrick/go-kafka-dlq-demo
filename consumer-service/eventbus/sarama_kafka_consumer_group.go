package eventbus

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

// SaramaKafkaConsumserGroupAdapterConfig is a configuration.
//
// FIELDS:
//
// `ConsumerGroupClient` is client that returned from `sarama.NewConsumerGroup()`.
//
//
// `ConsumerGroupHandler` is an implementation of `sarama.ConsumerGroupHandler`.
//
//
// `Topics` is a kafka topic to be subscribed.
type SaramaKafkaConsumserGroupAdapterConfig struct {
	ConsumerGroupClient  sarama.ConsumerGroup
	ConsumerGroupHandler sarama.ConsumerGroupHandler
	Topics               []string
}

// SaramaKafkaConsumserGroupAdapter is an adapter for eventbus's subcriber
type SaramaKafkaConsumserGroupAdapter struct {
	logger    *logrus.Logger
	closeChan chan struct{}
	config    *SaramaKafkaConsumserGroupAdapterConfig
}

// NewSaramaKafkaConsumerGroupFullConfigAdapter is constructor that immediately returns subscriber.
func NewSaramaKafkaConsumerGroupFullConfigAdapter(
	logger *logrus.Logger, addresses []string, groupID string, topics []string,
	consumerGroupHandler sarama.ConsumerGroupHandler,
	saramaConfig *sarama.Config,
) (subscriber Subscriber, err error) {

	consumerGroupClient, err := sarama.NewConsumerGroup(addresses, groupID, saramaConfig)
	if err != nil {
		return
	}

	closeChan := make(chan struct{}, 1)
	config := &SaramaKafkaConsumserGroupAdapterConfig{
		ConsumerGroupClient:  consumerGroupClient,
		ConsumerGroupHandler: consumerGroupHandler,
		Topics:               topics,
	}

	subscriber = &SaramaKafkaConsumserGroupAdapter{
		closeChan: closeChan,
		logger:    logger,
		config:    config,
	}

	return
}

// NewSaramaKafkaConsumserGroupAdapter is a constructor
//
// This Constructor is deprecated and use `NewSaramaKafkaConsumerGroupFullConfigAdapter` instead.
func NewSaramaKafkaConsumserGroupAdapter(logger *logrus.Logger, config *SaramaKafkaConsumserGroupAdapterConfig) Subscriber {
	closeChan := make(chan struct{}, 1)
	return &SaramaKafkaConsumserGroupAdapter{logger, closeChan, config}
}

// Subscribe will consume the published message
func (skcga *SaramaKafkaConsumserGroupAdapter) Subscribe() {
	go func() {
	POLL:
		for {
			select {
			case <-skcga.closeChan:
				break POLL
			default:
				err := skcga.config.ConsumerGroupClient.Consume(context.Background(), skcga.config.Topics, skcga.config.ConsumerGroupHandler)
				if err != nil {
					skcga.logger.Errorf("[Sarama] %s", err.Error())
				}
			}
		}
	}()

	return
}

// Close will stop the kafka consumer
func (skcga *SaramaKafkaConsumserGroupAdapter) Close() (err error) {
	defer close(skcga.closeChan)

	skcga.closeChan <- struct{}{}

	if err = skcga.config.ConsumerGroupClient.Close(); err != nil {
		skcga.logger.Errorf("[Sarama] Consumer is closed with error. | %s", err.Error())
		return
	}

	skcga.logger.Info("[Sarama] Consumer is gracefully shut down.")
	return
}

// SaramaConsumerGroup is an interface that purposed for mock creation for unit testing.
// Do not use this for an implementation.
type SaramaConsumerGroup interface {
	Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error
	Errors() <-chan error
	Close() error
}

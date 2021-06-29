package eventbus_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Shopify/sarama"
	eventbus "github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/eventbus"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/eventbus/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm"
)

func TestNewSaramaKafkaConsumerGroupFullConfigAdapter_Success(t *testing.T) {
	subscriber, err := eventbus.NewSaramaKafkaConsumerGroupFullConfigAdapter(
		logrus.New(), []string{"kafka1.com", "kafka2.com", "kafka3.com"}, "test-group", []string{"test-topic"},
		eventbus.NewDefaultSaramaConsumerGroupHandler(apm.DefaultTracer, "service-test", &mocks.EventHandler{}, &mocks.DLQHandler{}),
		sarama.NewConfig(),
	)
	assert.Error(t, err)
	assert.Nil(t, subscriber)
}

func TestSaramaKafkaConsumserGroupAdapter_Success(t *testing.T) {
	cgHandler := eventbus.NewDefaultSaramaConsumerGroupHandler(apm.DefaultTracer, "service-test", &mocks.EventHandler{}, &mocks.DLQHandler{})
	topics := []string{"test-topic"}

	cg := new(mocks.SaramaConsumerGroup)
	cg.On("Consume", mock.Anything, mock.AnythingOfType("[]string"), mock.AnythingOfType("*eventbus.DefaultSaramaConsumerGroupHandler")).Return(nil)
	// cg.On("Consume", mock.Anything, mock.AnythingOfType("[]string"), mock.AnythingOfType("*eventbus.DefaultSaramaConsumerGroupHandler")).Return(sarama.ErrOutOfBrokers)
	cg.On("Close").Return(nil)

	subscriber := eventbus.NewSaramaKafkaConsumserGroupAdapter(logrus.New(), &eventbus.SaramaKafkaConsumserGroupAdapterConfig{
		ConsumerGroupClient:  cg,
		ConsumerGroupHandler: cgHandler,
		Topics:               topics,
	})

	subscriber.Subscribe()
	<-time.After(time.Millisecond * 10)
	subscriber.Close()

	cg.AssertExpectations(t)
}

func TestSaramaKafkaConsumserGroupAdapter_ConsumeError(t *testing.T) {
	cgHandler := eventbus.NewDefaultSaramaConsumerGroupHandler(apm.DefaultTracer, "service-test", &mocks.EventHandler{}, &mocks.DLQHandler{})
	topics := []string{"test-topic"}

	cg := new(mocks.SaramaConsumerGroup)
	// cg.On("Consume", mock.Anything, mock.AnythingOfType("[]string"), mock.AnythingOfType("*eventbus.DefaultSaramaConsumerGroupHandler")).Return(nil)
	cg.On("Consume", mock.Anything, mock.AnythingOfType("[]string"), mock.AnythingOfType("*eventbus.DefaultSaramaConsumerGroupHandler")).Return(sarama.ErrOutOfBrokers)
	cg.On("Close").Return(nil)

	subscriber := eventbus.NewSaramaKafkaConsumserGroupAdapter(logrus.New(), &eventbus.SaramaKafkaConsumserGroupAdapterConfig{
		ConsumerGroupClient:  cg,
		ConsumerGroupHandler: cgHandler,
		Topics:               topics,
	})

	subscriber.Subscribe()
	<-time.After(time.Millisecond * 10)
	subscriber.Close()

	cg.AssertExpectations(t)
}

func TestSaramaKafkaConsumserGroupAdapter_ClosingError(t *testing.T) {
	cgHandler := eventbus.NewDefaultSaramaConsumerGroupHandler(apm.DefaultTracer, "service-test", &mocks.EventHandler{}, &mocks.DLQHandler{})
	topics := []string{"test-topic"}

	cg := new(mocks.SaramaConsumerGroup)
	cg.On("Consume", mock.Anything, mock.AnythingOfType("[]string"), mock.AnythingOfType("*eventbus.DefaultSaramaConsumerGroupHandler")).Return(nil)
	cg.On("Close").Return(sarama.ErrBrokerNotAvailable)

	subscriber := eventbus.NewSaramaKafkaConsumserGroupAdapter(logrus.New(), &eventbus.SaramaKafkaConsumserGroupAdapterConfig{
		ConsumerGroupClient:  cg,
		ConsumerGroupHandler: cgHandler,
		Topics:               topics,
	})

	subscriber.Subscribe()
	<-time.After(time.Millisecond * 10)
	subscriber.Close()

	cg.AssertExpectations(t)
}

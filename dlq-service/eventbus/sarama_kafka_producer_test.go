package eventbus_test

import (
	"context"
	"testing"

	"github.com/Shopify/sarama"

	"github.com/Shopify/sarama/mocks"
	eventbus "github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/eventbus"
	"github.com/sirupsen/logrus"
)

func TestSaramaKafkaProducer(t *testing.T) {
	t.Run("when the message is successfully sent and delivered", func(t *testing.T) {
		saramaProducer := mocks.NewAsyncProducer(t, nil)
		saramaProducer.ExpectInputAndSucceed()

		publisher := eventbus.NewSaramaKafkaProducerAdapter(logrus.New(), &eventbus.SaramaKafkaProducerAdapterConfig{
			AsyncProducer: saramaProducer,
		})

		headers := eventbus.MessageHeaders{}
		headers.Add("test", "header")

		publisher.Send(context.TODO(), "test-topic", "test-key", headers, []byte("Hola"))
		publisher.Close()
	})

	t.Run("when the connection is timeout", func(t *testing.T) {
		saramaProducer := mocks.NewAsyncProducer(t, nil)
		saramaProducer.ExpectInputAndFail(sarama.ErrRequestTimedOut)

		publisher := eventbus.NewSaramaKafkaProducerAdapter(logrus.New(), &eventbus.SaramaKafkaProducerAdapterConfig{
			AsyncProducer: saramaProducer,
		})

		headers := eventbus.MessageHeaders{}
		headers.Add("test", "header")

		publisher.Send(context.TODO(), "test-topic", "test-key", headers, []byte("Hola"))
		publisher.Close()
	})

	t.Run("when the message channel is already closed", func(t *testing.T) {
		saramaProducer := mocks.NewAsyncProducer(t, nil)

		publisher := eventbus.NewSaramaKafkaProducerAdapter(logrus.New(), &eventbus.SaramaKafkaProducerAdapterConfig{
			AsyncProducer: saramaProducer,
		})

		headers := eventbus.MessageHeaders{}
		headers.Add("test", "header")

		publisher.Close()
		publisher.Send(context.TODO(), "test-topic", "test-key", headers, []byte("Hola"))
	})
}

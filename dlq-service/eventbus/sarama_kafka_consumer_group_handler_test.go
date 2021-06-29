package eventbus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Shopify/sarama"

	eventbus "github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/eventbus"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/eventbus/mocks"
	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm"
)

func geConsumerMessageMock() <-chan *sarama.ConsumerMessage {
	messageChan := make(chan *sarama.ConsumerMessage, 1)
	defer close(messageChan)

	message := &sarama.ConsumerMessage{
		Headers: []*sarama.RecordHeader{
			&sarama.RecordHeader{Key: []byte("test"), Value: []byte("header")},
		},
		Key:       []byte("test-key"),
		Value:     []byte("test-message"),
		Partition: int32(1),
		Offset:    int64(40),
		Topic:     "test-topic",
	}
	messageChan <- message
	return messageChan
}

func TestSaramaKafkaConsumerGroupHandler_SuccessProceedMessage_WithoutEventHandler_WithoutDLQHandler(t *testing.T) {

	cgSess := &mocks.SaramaConsumerGroupSession{}
	cgSess.On("MarkMessage", mock.AnythingOfType("*sarama.ConsumerMessage"), mock.AnythingOfType("string"))
	cgSess.On("Context").Return(context.TODO())

	cgClaim := &mocks.SaramaConsumerGroupClaim{}
	cgClaim.On("Messages").Return(geConsumerMessageMock())

	cgh := eventbus.NewDefaultSaramaConsumerGroupHandler(apm.DefaultTracer, "service-test", nil, nil)
	cgh.Setup(cgSess)
	cgh.ConsumeClaim(cgSess, cgClaim)
	cgh.Cleanup(cgSess)

	cgSess.AssertExpectations(t)
	cgClaim.AssertExpectations(t)
}

func TestSaramaKafkaConsumerGroupHandler_SuccessProceedMessage_WithEventHandler_WithoutDLQHandler(t *testing.T) {
	eventHandler := &mocks.EventHandler{}
	eventHandler.On("Handle", mock.Anything, mock.Anything).Return(nil)

	cgSess := &mocks.SaramaConsumerGroupSession{}
	cgSess.On("MarkMessage", mock.AnythingOfType("*sarama.ConsumerMessage"), mock.AnythingOfType("string"))
	cgSess.On("Context").Return(context.TODO())

	cgClaim := &mocks.SaramaConsumerGroupClaim{}
	cgClaim.On("Messages").Return(geConsumerMessageMock())

	cgh := eventbus.NewDefaultSaramaConsumerGroupHandler(apm.DefaultTracer, "service-test", eventHandler, nil)
	cgh.Setup(cgSess)
	cgh.ConsumeClaim(cgSess, cgClaim)
	cgh.Cleanup(cgSess)

	cgSess.AssertExpectations(t)
	cgClaim.AssertExpectations(t)
	eventHandler.AssertExpectations(t)
}

func TestSaramaKafkaConsumerGroupHandler_ErrorProceedMessage_WithEventHandler_WithoutDLQHandler(t *testing.T) {
	eventHandler := &mocks.EventHandler{}
	eventHandler.On("Handle", mock.Anything, mock.Anything).Return(fmt.Errorf("error"))

	cgSess := &mocks.SaramaConsumerGroupSession{}
	cgSess.On("MarkMessage", mock.AnythingOfType("*sarama.ConsumerMessage"), mock.AnythingOfType("string"))
	cgSess.On("Context").Return(context.TODO())

	cgClaim := &mocks.SaramaConsumerGroupClaim{}
	cgClaim.On("Messages").Return(geConsumerMessageMock())

	cgh := eventbus.NewDefaultSaramaConsumerGroupHandler(apm.DefaultTracer, "service-test", eventHandler, nil)
	cgh.Setup(cgSess)
	cgh.ConsumeClaim(cgSess, cgClaim)
	cgh.Cleanup(cgSess)

	cgSess.AssertExpectations(t)
	cgClaim.AssertExpectations(t)
	eventHandler.AssertExpectations(t)
}

func TestSaramaKafkaConsumerGroupHandler_ErrorProceedMessage_WithEventHandler_WithDLQHandler(t *testing.T) {
	eventHandler := &mocks.EventHandler{}
	eventHandler.On("Handle", mock.Anything, mock.Anything).Return(fmt.Errorf("error"))

	// Send(ctx context.Context, topic string, key string, headers MessageHeaders, message []byte) (err error)
	publisher := &mocks.Publisher{}
	publisher.On("Send", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("[]uint8")).Return(nil)

	dlqHandler := eventbus.NewDLQHandlerAdapter("dlq-test-topic", publisher)

	cgSess := &mocks.SaramaConsumerGroupSession{}
	cgSess.On("MarkMessage", mock.AnythingOfType("*sarama.ConsumerMessage"), mock.AnythingOfType("string"))
	cgSess.On("Context").Return(context.TODO())

	cgClaim := &mocks.SaramaConsumerGroupClaim{}
	cgClaim.On("Messages").Return(geConsumerMessageMock())

	cgh := eventbus.NewDefaultSaramaConsumerGroupHandler(apm.DefaultTracer, "service-test", eventHandler, dlqHandler)
	cgh.Setup(cgSess)
	cgh.ConsumeClaim(cgSess, cgClaim)
	cgh.Cleanup(cgSess)

	cgSess.AssertExpectations(t)
	cgClaim.AssertExpectations(t)
	eventHandler.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

package eventbus

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"go.elastic.co/apm"
)

// DefaultSaramaConsumerGroupHandler is a default consumer group handler for sarama kafka client.
// Create your own to address some customization.
// It is the implementation of `sarama.KafkaConsumerGroupHandler`
type DefaultSaramaConsumerGroupHandler struct {
	utcTZ        *time.Location
	tracer       *apm.Tracer
	serviceName  string
	eventHandler EventHandler
	dlqHandler   DLQHandler
}

// NewDefaultSaramaConsumerGroupHandler is a constructor.
func NewDefaultSaramaConsumerGroupHandler(tracer *apm.Tracer, serviceName string, eventHandler EventHandler, dlqHandler DLQHandler) *DefaultSaramaConsumerGroupHandler {
	utcTz, _ := time.LoadLocation("UTC")
	return &DefaultSaramaConsumerGroupHandler{
		utcTZ:        utcTz,
		tracer:       tracer,
		eventHandler: eventHandler,
		dlqHandler:   dlqHandler,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *DefaultSaramaConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	// close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *DefaultSaramaConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *DefaultSaramaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		consumer.claim(session.Context(), message)
		session.MarkMessage(message, "")
	}

	return nil
}

func (consumer *DefaultSaramaConsumerGroupHandler) claim(ctx context.Context, message *sarama.ConsumerMessage) {
	txName := fmt.Sprintf("On Event: %s", message.Topic)
	txType := "Kafka Consumer"
	txSuccess := "Success"

	tx := consumer.tracer.StartTransaction(txName, txType)
	defer tx.End()

	ctx = apm.ContextWithTransaction(ctx, tx)

	if consumer.eventHandler == nil {
		consumer.printMessage(message)
		tx.Result = txSuccess
		return
	}

	if err := consumer.eventHandler.Handle(ctx, message); err != nil {
		consumer.sendToDLQ(ctx, message, err)
		tx.Result = err.Error()
	}
	return
}

func (consumer *DefaultSaramaConsumerGroupHandler) printMessage(message *sarama.ConsumerMessage) {
	log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, partition = %d", string(message.Value), message.Timestamp, message.Topic, message.Partition)
}

func (consumer *DefaultSaramaConsumerGroupHandler) sendToDLQ(ctx context.Context, message *sarama.ConsumerMessage, err error) {
	if consumer.dlqHandler == nil {
		return
	}

	originalHeaders := MessageHeaders{}

	for _, header := range message.Headers {
		originalHeaders.Add(string(header.Key), string(header.Value))
	}

	dlqMessage := &DeadLetterQueueMessage{
		Channel:           message.Topic,
		Publisher:         originalHeaders["origin"],
		Consumer:          consumer.serviceName,
		Key:               string(message.Key),
		Headers:           originalHeaders,
		Message:           string(message.Value),
		CausedBy:          err.Error(),
		FailedConsumeDate: message.Timestamp.In(consumer.utcTZ).Format(time.RFC3339Nano),
	}

	consumer.dlqHandler.Send(ctx, dlqMessage)
}

// SaramaConsumerGroupSession is an interface that purposed for mock creation for unit testing.
// Do not use this for an implementation.
type SaramaConsumerGroupSession interface {
	Claims() map[string][]int32
	MemberID() string
	GenerationID() int32
	MarkOffset(topic string, partition int32, offset int64, metadata string)
	Commit()
	ResetOffset(topic string, partition int32, offset int64, metadata string)
	MarkMessage(msg *sarama.ConsumerMessage, metadata string)
	Context() context.Context
}

// SaramaConsumerGroupClaim is an interface that purposed for mock creation for unit testing.
// Do not use this for an implementation.
type SaramaConsumerGroupClaim interface {
	Topic() string
	Partition() int32
	InitialOffset() int64
	HighWaterMarkOffset() int64
	Messages() <-chan *sarama.ConsumerMessage
}

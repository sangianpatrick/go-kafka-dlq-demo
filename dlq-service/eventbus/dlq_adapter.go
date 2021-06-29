package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// DeadLetterQueueMessage is an entity.
type DeadLetterQueueMessage struct {
	Channel           string         `json:"channel"`
	Publisher         string         `json:"publisher"`
	Consumer          string         `json:"consumer"`
	Key               string         `json:"key"`
	Headers           MessageHeaders `json:"headers"`
	Message           string         `json:"message"`
	CausedBy          string         `json:"causedBy"`
	FailedConsumeDate string         `json:"failedConsumeDate"`
}

// DLQHandlerAdapter is an dead letter queue adapter.
type DLQHandlerAdapter struct {
	topic     string
	publisher Publisher
}

// NewDLQHandlerAdapter is a constructor.
func NewDLQHandlerAdapter(topic string, publisher Publisher) *DLQHandlerAdapter {
	return &DLQHandlerAdapter{topic, publisher}
}

// Send will publish the dlq message to the assigned topic.
func (dlqHandlerAdapter *DLQHandlerAdapter) Send(ctx context.Context, dlqMessage *DeadLetterQueueMessage) (err error) {
	headers := MessageHeaders{}
	headers.Add("dlq", "true")

	key := fmt.Sprintf("%s:%s:%s:%d",
		dlqMessage.Consumer,
		dlqMessage.Channel,
		dlqMessage.Key,
		time.Now().UnixNano(),
	)

	messageByte, _ := json.Marshal(dlqMessage)
	topics := dlqHandlerAdapter.topic

	err = dlqHandlerAdapter.publisher.Send(
		ctx,
		topics,
		key,
		headers,
		messageByte,
	)

	return
}

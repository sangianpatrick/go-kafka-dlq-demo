package eventbus

import "context"

// DLQHandler is an handler to handler dead letter queue or an unprocessed message
type DLQHandler interface {
	Send(ctx context.Context, dlqMessage *DeadLetterQueueMessage) (err error)
}

// EventHandler is an event handler. It will be called after message is arrived to consumer
type EventHandler interface {
	Handle(ctx context.Context, message interface{}) (err error)
}

// Publisher is a collection of behavior of a publisher
type Publisher interface {
	// Will send the message to the assigned topic.
	Send(ctx context.Context, topic string, key string, headers MessageHeaders, message []byte) (err error)
	Close() (err error)
}

// Subscriber is a collection of behavior of a subscriber
type Subscriber interface {
	Subscribe()
	Close() (err error)
}

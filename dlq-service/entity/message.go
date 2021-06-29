package entity

// MessageHeaders is an entity.
type MessageHeaders map[string]string

// Message is an entity.
type Message struct {
	Channel           string         `json:"channel" bson:"channel"`
	Publisher         string         `json:"publisher" bson:"publisher"`
	Consumer          string         `json:"consumer" bson:"consumer"`
	Key               string         `json:"key" bson:"key"`
	Headers           MessageHeaders `json:"headers" bson:"headers"`
	Message           string         `json:"message" bson:"message"`
	CausedBy          string         `json:"causedBy" bson:"causedBy"`
	FailedConsumeDate string         `json:"failedConsumeDate" bson:"failedConsumeDate"`
}

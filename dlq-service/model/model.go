package model

import "time"

// MessageHeadersParams is a model.
type MessageHeadersParams map[string]string

// MessageParams is a model.
type MessageParams struct {
	Channel           string               `json:"channel" validate:"required"`
	Publisher         string               `json:"publisher" validate:"required"`
	Consumer          string               `json:"consumer" validate:"required"`
	Key               string               `json:"key"  validate:"required"`
	Headers           MessageHeadersParams `json:"headers" validate:"required"`
	Message           string               `json:"message" validate:"required"`
	CausedBy          string               `json:"causedBy" validate:"required"`
	FailedConsumeDate time.Time            `json:"failedConsumeDate" validate:"required"`
}

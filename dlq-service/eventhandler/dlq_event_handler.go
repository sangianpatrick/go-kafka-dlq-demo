package eventhandler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Shopify/sarama"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/eventbus"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/model"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/usecase"
	"github.com/sirupsen/logrus"
)

type dlqEventHandler struct {
	logger  *logrus.Logger
	usecase usecase.DLQUsecase
}

func NewDLQEventHandler(logger *logrus.Logger, usecase usecase.DLQUsecase) eventbus.EventHandler {
	return &dlqEventHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (handler *dlqEventHandler) Handle(ctx context.Context, eventMessage interface{}) (err error) {
	kafkaMessage, ok := eventMessage.(*sarama.ConsumerMessage)
	if !ok {
		err = errors.New("not a kafka message")
		handler.logger.Error(err)

		return
	}

	var payload model.MessageParams

	if err = json.Unmarshal(kafkaMessage.Value, &payload); err != nil {
		handler.logger.Error(err)
		return
	}

	response := handler.usecase.Add(ctx, payload)

	if !response.IsSuccess {
		err = errors.New(response.Error)
		return
	}

	return
}

package usecase

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/entity"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/eventbus"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/model"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/repository"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

type meta struct {
	Page            int `json:"page"`
	TotalPage       int `json:"totalPage"`
	TotalDataOnPage int `json:"totalDataOnPage"`
	TotalData       int `json:"totalData"`
}

type DLQUsecase interface {
	Add(ctx context.Context, payload model.MessageParams) (response model.Response)
	GetMany(ctx context.Context, page, size int64) (response model.Response)
	Republish(ctx context.Context, ID string) (response model.Response)
}

type dlqUsecase struct {
	logger     *logrus.Logger
	publisher  eventbus.Publisher
	repository repository.DLQRepository
}

func NewDLQUsecase(logger *logrus.Logger, repository repository.DLQRepository) DLQUsecase {
	return &dlqUsecase{
		logger:     logger,
		repository: repository,
	}
}

func (u *dlqUsecase) Add(ctx context.Context, payload model.MessageParams) (response model.Response) {
	dlqMessage := entity.Message{}

	headers := entity.MessageHeaders{}
	for hk, hv := range payload.Headers {
		headers[hk] = hv
	}

	dlqMessage.ID = uuid.New().String()
	dlqMessage.Key = payload.Key
	dlqMessage.Channel = payload.Channel
	dlqMessage.Consumer = payload.Consumer
	dlqMessage.Publisher = payload.Publisher
	dlqMessage.Headers = headers
	dlqMessage.Message = payload.Message
	dlqMessage.CausedBy = payload.CausedBy
	dlqMessage.FailedConsumeDate = payload.FailedConsumeDate

	if err := u.repository.InsertOne(ctx, dlqMessage); err != nil {
		response.Status = model.StatusInternalServerError
		response.Error = err.Error()

		return
	}

	response.IsSuccess = true
	response.Status = model.StatusCreated

	return
}

func (u *dlqUsecase) GetMany(ctx context.Context, page, size int64) (response model.Response) {
	skip := (page - 1) * size
	limit := size

	bunchOfMessage, totalCounted, err := u.getMany(ctx, limit, skip)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response.Status = model.StatusNotFoundError
			response.Error = err.Error()

			return
		}

		response.Status = model.StatusInternalServerError
		response.Error = err.Error()

		return
	}

	lengthOfMessages := len(bunchOfMessage)

	response.Meta = meta{
		Page:            int(page),
		TotalPage:       int(math.Ceil(float64(totalCounted) / float64(size))),
		TotalDataOnPage: lengthOfMessages,
		TotalData:       int(totalCounted),
	}

	response.IsSuccess = true
	response.Status = model.StatusOK
	response.Data = bunchOfMessage

	return
}

func (u *dlqUsecase) Republish(ctx context.Context, ID string) (response model.Response) {
	dlqMessage, err := u.repository.FindByID(ctx, ID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			response.Status = model.StatusNotFoundError
			response.Error = err.Error()

			return
		}

		response.Status = model.StatusInternalServerError
		response.Error = err.Error()

		return
	}

	headers := eventbus.MessageHeaders{}
	for hk, hv := range dlqMessage.Headers {
		headers.Add(hk, hv)
	}

	u.publisher.Send(ctx, dlqMessage.Channel, dlqMessage.ID, headers, []byte(dlqMessage.Message))

	if err := u.repository.DeleteByID(ctx, ID); err != nil {
		response.Status = model.StatusInternalServerError
		response.Error = err.Error()

		return
	}

	response.Status = model.StatusOK
	response.Data = dlqMessage

	return
}

func (u *dlqUsecase) getMany(ctx context.Context, limit, skip int64) (bunchOfDLQMessage []entity.Message, totalCounted int64, err error) {
	resultAggregation := make(map[string]interface{})

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		counted, err := u.repository.CountDocuments(ctx)
		if err != nil {
			return err
		}

		resultAggregation["counted"] = counted

		return nil
	})

	g.Go(func() error {
		bunchOfDLQMessage, err := u.repository.FindMany(ctx, limit, skip)
		if err != nil {
			return err
		}

		resultAggregation["bunchOfDLQMessage"] = bunchOfDLQMessage

		return nil
	})

	if err = g.Wait(); err != nil {
		return
	}

	totalCounted = resultAggregation["counted"].(int64)
	bunchOfDLQMessage = resultAggregation["bunchOfDLQMessage"].([]entity.Message)

	return
}

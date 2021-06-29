package usecase

import (
	"context"
	"math"

	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/entity"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/model"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type meta struct {
	Page            int `json:"page"`
	TotalPage       int `json:"totalPage"`
	TotalDataOnPage int `json:"totalDataOnPage"`
	TotalData       int `json:"totalData"`
}

type DLQUsecase interface {
	GetMany(ctx context.Context, page, size int64) (response model.Response)
	Republish(ctx context.Context, ID string) (response model.Response)
}

type dlqUsecase struct {
	logger     *logrus.Logger
	repository repository.DLQRepository
}

func NewDLQUsecase(logger *logrus.Logger, repository repository.DLQRepository) DLQUsecase {
	return &dlqUsecase{
		logger:     logger,
		repository: repository,
	}
}

func (u *dlqUsecase) GetMany(ctx context.Context, page, size int64) (response model.Response) {
	skip := (page - 1) * size
	limit := size

	bunchOfMessage, totalCounted, err := u.getMany(ctx, limit, skip)
	if err != nil {
		response.Status = model.StatusInternalServerError
		response.Error = err.Error()

		return
	}

	response.Meta = meta{
		Page:            int(page),
		TotalPage:       int(math.Ceil(float64(totalCounted) / float64(size))),
		TotalDataOnPage: len(bunchOfMessage),
		TotalData:       int(totalCounted),
	}
	response.Status = model.StatusOK
	response.Data = bunchOfMessage

	return
}

func (u *dlqUsecase) Republish(ctx context.Context, ID string) (response model.Response) {
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

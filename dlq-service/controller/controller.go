package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/model"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/usecase"
	"github.com/sirupsen/logrus"
)

type DLQController struct {
	Logger  *logrus.Logger
	Usecase usecase.DLQUsecase
}

func InitDLQController(logger *logrus.Logger, router *mux.Router, usecase usecase.DLQUsecase) {
	controller := &DLQController{
		Logger:  logger,
		Usecase: usecase,
	}

	router.HandleFunc("/dlq-service/messages", controller.GetMany).Methods(http.MethodGet)
}

func (c *DLQController) GetMany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	queryString := r.URL.Query()
	page, _ := strconv.ParseInt(queryString.Get("page"), 10, 64)
	size, _ := strconv.ParseInt(queryString.Get("size"), 10, 64)

	response := c.Usecase.GetMany(ctx, page, size)

	httpStatusCode := model.GetHTTPStatusCodeByResponseStatus(response.Status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	json.NewEncoder(w).Encode(response)
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload" // for development
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/controller"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/eventbus"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/eventhandler"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/mongodb"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/repository"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/usecase"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmgorilla"
	"go.elastic.co/apm/module/apmlogrus"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	tracer *apm.Tracer
)

func init() {
	tracer = apm.DefaultTracer
}

func main() {
	serviceName := os.Getenv("SERVICE_NAME")
	servicePort, _ := strconv.Atoi(os.Getenv("SERVICE_PORT"))
	mongodbURL := os.Getenv("MONGODB_URL")
	mongodbDatabase := os.Getenv("MONGODB_DATABASE")
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	})
	logger.SetReportCaller(true)
	logger.AddHook(&apmlogrus.Hook{
		LogLevels: logrus.AllLevels,
	})

	mongodbClient, err := mongo.NewClient(
		options.Client().
			SetAppName(serviceName).
			ApplyURI(mongodbURL).
			SetConnectTimeout(time.Second * 10).
			SetMonitor(apmmongo.CommandMonitor()),
	)

	if err != nil {
		logger.Fatal(err)
	}

	dbClient := mongodb.NewClientAdapter(mongodbClient)
	if dbClient.Connect(context.Background()); err != nil {
		logger.Fatal((err))
	}

	db := dbClient.Database(mongodbDatabase)

	router := mux.NewRouter()
	apmgorilla.Instrument(router)

	dlqRepository := repository.NewDLQRepository(logger, db)
	dlqUsecase := usecase.NewDLQUsecase(logger, dlqRepository)
	controller.InitDLQController(logger, router, dlqUsecase)

	consumerGroupClient, err := sarama.NewConsumerGroup(kafkaBrokers, serviceName, sarama.NewConfig())
	if err != nil {
		logger.Fatal(err)
	}
	consumerGroupHandler := eventbus.NewDefaultSaramaConsumerGroupHandler(tracer, serviceName, eventhandler.NewDLQEventHandler(logger, dlqUsecase), nil)
	subscriber := eventbus.NewSaramaKafkaConsumserGroupAdapter(
		logger, &eventbus.SaramaKafkaConsumserGroupAdapterConfig{
			ConsumerGroupClient:  consumerGroupClient,
			ConsumerGroupHandler: consumerGroupHandler,
			Topics:               []string{"dead-letter-queue"},
		})
	subscriber.Subscribe()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", servicePort),
		Handler: router,
	}

	go func() {
		logger.Infof("server is running on port %d", servicePort)
		httpServer.ListenAndServe()
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	httpServer.Shutdown(context.Background())
	subscriber.Close()
	dbClient.Disconnect(context.Background())
}

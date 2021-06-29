module github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service

go 1.15

require (
	github.com/Shopify/sarama v1.29.1
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.3.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	go.elastic.co/apm v1.12.0
	go.elastic.co/apm/module/apmgorilla v1.12.0
	go.elastic.co/apm/module/apmlogrus v1.12.0
	go.mongodb.org/mongo-driver v1.5.3
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
)

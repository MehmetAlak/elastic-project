package main

import (
	"context"
	"elastic-project/application/elastic_operation"
	"elastic-project/client/elasticsearch"
	"elastic-project/interface/rest"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	gracefulShutdown := createGracefulShutdownChannel()

	elastic, err := elasticsearch.New([]string{"http://0.0.0.0:9200"})
	if err != nil {
		log.Fatalln(err)
	}
	if err := elastic.CreateIndex("user"); err != nil {
		log.Fatalln(err)
	}

	storage := elasticsearch.NewUserInfoStorage(*elastic)

	elasticsearchService := elastic_operation.NewElasticsearchService(storage)
	elasticsearchEndpoint := rest.NewElasticsearchEndpoint(elasticsearchService)

	server := rest.NewServer(elasticsearchEndpoint)

	router := server.SetupRouter()
	_ = router.Run(":8084")

	<-gracefulShutdown
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

}

func createGracefulShutdownChannel() chan os.Signal {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	return gracefulShutdown
}

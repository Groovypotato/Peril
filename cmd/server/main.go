package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	cstring := "amqp://guest:guest@localhost:5672/"
	rconnect, err := amqp.Dial(cstring)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rconnect.Close()
	fmt.Println("Connection was successfull!")
	newChan, err := rconnect.Channel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pubsub.PublishJSON(newChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
	<-sigChan
	fmt.Println("\nReceived interrupt signal. Performing graceful shutdown...")

}

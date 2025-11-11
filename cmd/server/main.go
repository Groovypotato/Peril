package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	sigServerChan := make(chan os.Signal, 1)
	signal.Notify(sigServerChan, syscall.SIGINT, syscall.SIGTERM)
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
	_, _, err = pubsub.DeclareAndBind(rconnect, "peril_topic", "game_logs", "game_logs.*", pubsub.DurableQueue)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	gamelogic.PrintServerHelp()
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		} else {
			if strings.ToLower(input[0]) == "pause" {
				fmt.Println("sending 'pause' command")
				pubsub.PublishJSON(newChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			} else if strings.ToLower(input[0]) == "resume" {
				fmt.Println("sending 'resume' command")
				pubsub.PublishJSON(newChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			} else if strings.ToLower(input[0]) == "quit" {
				fmt.Println("sending 'quit' command")
				break
			} else {
				fmt.Println("¿Que?")
			}
		}
	}

	<-sigServerChan
	fmt.Println("\nReceived interrupt signal. Performing graceful shutdown...")
}

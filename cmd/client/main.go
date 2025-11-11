package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
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

	uName, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	queueName := routing.PauseKey + "." + uName
	var queueType pubsub.SimpleQueueType = "transient"

	_, _, err = pubsub.DeclareAndBind(rconnect, routing.ExchangePerilDirect, queueName, routing.PauseKey, queueType)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gameState := gamelogic.NewGameState(uName)
	err = pubsub.SubscribeJSON(rconnect, routing.ExchangePerilDirect, queueName, routing.PauseKey, queueType, handlerPause(gameState))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
Loop:
	for {
		uInput := gamelogic.GetInput()

		command := strings.ToLower(uInput[0])

		switch command {
		case "spawn":
			if len(uInput) != 3 {
				fmt.Println("usage: <command> <location> <unit> ")
				continue
			}
			switch strings.ToLower(uInput[2]) {
			case "infantry":
			case "cavalry":
			case "artillery":
			default:
				fmt.Println("unknown unit")
				continue
			}

			switch uInput[1] {
			case "americas":
			case "europe":
			case "africa":
			case "asia":
			case "antarctica":
			case "australia":
			default:
				fmt.Println("unknown location")
				continue
			}

			err = gameState.CommandSpawn(uInput)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
		case "move":
			if len(uInput) != 3 {
				fmt.Println("usage: <command> <location> <unit id>")
				continue
			}
			switch uInput[1] {
			case "americas":
			case "europe":
			case "africa":
			case "asia":
			case "antarctica":
			case "australia":
			default:
				fmt.Println("unknown location")
				continue
			}

			_, err = strconv.Atoi(uInput[2])
			if err != nil {
				fmt.Println("unit id is not a number")
				continue
			}

			move, err := gameState.CommandMove(uInput)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf("%s moved the following units to %s\n", move.Player.Username, move.ToLocation)
			for _, unit := range move.Units {
				fmt.Printf("ID:%d    Type:%s\n", unit.ID, unit.Rank)
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			break Loop
		default:
			println("unknown command")
		}

	}
	fmt.Println("quitting....")
	os.Exit(0)
}

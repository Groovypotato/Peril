package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)

func handlerMove(gs *gamelogic.GameState) func(m gamelogic.ArmyMove) {
	return func(m gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(m)
	}
}

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Minimax algorithm for Magic the Gathering
// plies are turn segments where the player holds priority

var pass = "PASS"

type action struct {
	card         string
	playerSelf   int
	playerTarget int
}

func (a action) execute(p, opp *player) {
	if a.card == pass {
		return
	}
	p.hand[a.card] -= 1
	if p.hand[a.card] == 0 {
		delete(p.hand, a.card)
	}
	card := cards[a.card]
	for _, effect := range card.effects {
		effect(p, opp)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	p1 := newPlayer("player1", deckList)
	p2 := newPlayer("player2", deckList)
	p1.drawN(7)
	p2.drawN(7)

	game := newGame([]*player{p1, p2})

	for {
		game.playUntilPriority()
		if gameEnds := game.checkStateBasedActions(); gameEnds {
			fmt.Println("End of game")
			break
		}
		action := game.getPlayerAction()
		if action.card == pass {
			game.nextStep()
		}
		action.execute(game.getPlayer(action.playerSelf), game.getPlayer(action.playerTarget))
		fmt.Printf("-> %s played %s \n", game.getPlayer(action.playerSelf).name, action.card)
		game.debug()
	}
}

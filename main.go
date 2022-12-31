package main

import (
	"math/rand"
	"time"
)

// Minimax algorithm for Magic the Gathering
// plies are turn segments where the player holds priority

func main() {
	rand.Seed(time.Now().UnixNano())

	p1 := newPlayer("player1", deckList)
	p2 := newPlayer("player2", deckList)
    p1.strategy = simpleStrategy{}
    p2.strategy = minmaxStrategy{}

	game := newGame(p1, p2)
	game.loop()
}

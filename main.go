package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	p1 := newPlayer("player1", deckList)
	p2 := newPlayer("player2", deckList)
    p1.strategy = simpleStrategy{}
    p2.strategy = minmaxStrategy{}

	startingPlayer := rand.Intn(2)

	game := newGame(startingPlayer, p1, p2)
	game.loop()
}

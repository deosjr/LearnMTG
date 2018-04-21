package main

import (
	"fmt"

	"github.com/MagicTheGathering/mtg-sdk-go"
)

var (
	mountain = card{
		name: "Mountain",
		prereqs: []prerequisiteFunc{
			func(p *player) bool { return !p.landPlayed },
		},
		effect: func(p *player) {
			// Maybe move all of this to 'execute play land action' ?
			p.landPlayed = true
			// TODO: land enters the battlefield
			p.manaTotal += 1
			p.manaAvailable += 1
		},
	}

	lavaSpike = card{
		name: "Lava Spike",
		prereqs: []prerequisiteFunc{
			func(p *player) bool {
				// TODO: target available (lets ignore hexproof players for now)
				return p.manaAvailable >= 1
			},
		},
		effect: func(p *player) {
			p.manaAvailable -= 1
			// TODO: lava spike can target self!
			p.opponent.lifeTotal -= 3
		},
	}

	cards = map[string]card{
		mountain.name:  mountain,
		lavaSpike.name: lavaSpike,
	}

	deckList = map[string]int{
		mountain.name:  10,
		lavaSpike.name: 20,
	}
)

// Eventually this will be a hierarchy
// i.e. lands all share some prerequisites and effects
// A type-token distinction needs to be made at some point
type card struct {
	name    string
	prereqs []prerequisiteFunc
	effect  effectFunc
}

type prerequisiteFunc func(*player) bool
type effectFunc func(*player)

type Permanent interface{}
type permanent struct {
	isTapped bool
}
type land struct {
	permanent
}

// TODO
func getCard(name string) card {
	cards, err := mtg.NewQuery().Where(mtg.CardName, name).All()
	if err != nil {
		panic(err)
	}
	fmt.Println(cards)
	return card{name: cards[0].Name}
}

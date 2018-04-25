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
		effects: []effectFunc{
			func(p, _ *player) {
				// Maybe move all of this to 'execute play land action' ?
				p.landPlayed = true
				// TODO: land enters the battlefield
				p.manaTotal += 1
				p.manaAvailable += 1
			},
		},
	}

	lavaSpike = card{
		name:     "Lava Spike",
		manacost: manacost{r: 1},
		prereqs: []prerequisiteFunc{
			func(p *player) bool {
				// TODO: target available (lets ignore hexproof players for now)
				return p.manaAvailable >= 1
			},
		},
		effects: []effectFunc{
			func(p, _ *player) {
				p.manaAvailable -= 1
			},
			func(_, opp *player) {
				// TODO: lava spike can target self!
				opp.lifeTotal -= 3
			},
		},
	}

	falkenrathReaver = card{
		name:     "Falkenrath Reaver",
		manacost: manacost{c: 1, r: 1},
		prereqs: []prerequisiteFunc{
			func(p *player) bool {
				return p.manaAvailable >= 2
			},
		},
		effects: []effectFunc{
			func(p, _ *player) {
				p.manaAvailable -= 2
			},
			// TODO: creature enters the battlefield
		},
		//power:     2,
		//toughness: 2,
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

type card struct {
	name     string
	manacost manacost
	prereqs  []prerequisiteFunc
	effects  []effectFunc
}

type manacost struct {
	c, w, u, b, r, g int
}

// TODO: on effect -> go on stack (except lands)
// onResolve -> type specific (go to battlefield/graveyard)
// and card specific (resolve sorcery/instant effect)

// TODO: stack -> first good reason for generic game state struct
// --> implies rewrite of minimax and main

// TODO: breakdown effectFunc into types by target

type prerequisiteFunc func(*player) bool
type effectFunc func(*player, *player)

// TODO: tap is for tokens, P/T on creature is about type!
type Permanent interface {
	tap() (success bool)
}

type permanent struct {
	isTapped bool
}

func (p *permanent) tap() (success bool) {
	if p.isTapped {
		return false
	}
	p.isTapped = true
	return true
}

type land struct {
	permanent
}
type creature struct {
	permanent
	power     int
	toughness int
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

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
		effects: []effect{
			selfEffect{
			effect: func(p *player) {
				// Maybe move all of this to 'execute play land action' ?
				p.landPlayed = true
				// TODO: land enters the battlefield
				p.manaTotal += 1
				p.manaAvailable += 1
				},
			},
		},
	}

	lavaSpike = card{
		name:     "Lava Spike",
		manacost: manacost{r: 1},
		prereqs: []prerequisiteFunc{
			func(p *player) bool {
				// TODO: target available (lets ignore hexproof players for now)
				// --> this is handled by possibleTargets returning 0 actions
				return p.manaAvailable >= 1
			},
		},
		effects: []effect{
			selfEffect{
			effect: func(p *player) {
				p.manaAvailable -= 1
				},
			},
			playerEffect{
				effect:func(p *player) {
				p.lifeTotal -= 3
				},
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
		effects: []effect{
			selfEffect{
			effect: func(p *player) {
				p.manaAvailable -= 2
			},
			// TODO: creature enters the battlefield
			},
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
	effects  []effect
}

// TODO: generate manacost prereq funcs using closures?
type manacost struct {
	c, w, u, b, r, g int
}

// TODO: on effect -> go on stack (except lands)
// onResolve -> type specific (go to battlefield/graveyard)
// and card specific (resolve sorcery/instant effect)

// For permanent targets this might get hairy if they change
// but players never change so this is simple
type effect interface {
	possibleTargets(controllingPlayer int, game *game) []effect
	apply(*game)  
}

type selfEffect struct{
	target int
	effect func(*player)
}

func (e selfEffect) possibleTargets(controllingPlayer int, game *game) []effect {
	return []effect{
		selfEffect{
			target: controllingPlayer,
			effect: e.effect,
		},
	}
}

func (e selfEffect) apply(game *game) {
	p := game.getPlayer(e.target)
	e.effect(p)
}

type playerEffect struct {
	controller int
	target int
	effect func(*player)
}

func (e playerEffect) possibleTargets(controllingPlayer int, game *game) []effect {
	effects := []effect{}
	for i := 0; i < game.numPlayers; i++ {
		pe := playerEffect{
			target: i,
			effect: e.effect,
		}
		effects = append(effects, pe)
	}
	return effects
}

func (e playerEffect) apply(game *game) {
	p := game.getPlayer(e.target)
	e.effect(p)
}

type prerequisiteFunc func(*player) bool

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

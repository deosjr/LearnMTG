package main

import (
	"fmt"

	"github.com/MagicTheGathering/mtg-sdk-go"
)

type unorderedCards map[string]int // cardName : amount
type orderedCards []string         // cardName

var (
	mountain = card{
		name:     "Mountain",
		cardType: land{},
	}

	lavaSpike = card{
		name:     "Lava Spike",
		manacost: manacost{r: 1},
		cardType: sorcery{},
		effects: []effect{
			playerEffect{
				effect: func(p *player) {
					p.lifeTotal -= 3
				},
			},
		},
	}

	falkenrathReaver = card{
		name:     "Falkenrath Reaver",
		manacost: manacost{c: 1, r: 1},
		cardType: creature{
			power:     2,
			toughness: 2,
		},
	}

	cards = map[string]card{
		mountain.name:  mountain,
		lavaSpike.name: lavaSpike,
	}

	deckList = unorderedCards{
		mountain.name:  10,
		lavaSpike.name: 20,
	}
)

type card struct {
	name     string
	manacost manacost
	cardType cardType
	prereqs  []prerequisiteFunc
	effects  []effect
}

type manacost struct {
	c, w, u, b, r, g int
}

func (m manacost) converted() int {
	return m.c + m.w + m.u + m.b + m.r + m.g
}

// For permanent targets this might get hairy if they change
// but players never change so this is simple
type effect interface {
	possibleTargets(controllingPlayer int, game *game) []effect
	apply(*game)
}

type selfEffect struct {
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
	target     int
	effect     func(*player)
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

type cardType interface {
	prereq(*game, *player) bool
	// sorcery/instant goes to graveyard from stack
	// permanents enter the battlefield (unless countered ofc)
	// TODO: on effect -> go on stack (except lands)
	// onResolve -> type specific (go to battlefield/graveyard)
	// and card specific (resolve sorcery/instant effect, see game.resolve)
	resolve(*player)
}

type sorcery struct{}

func (s sorcery) prereq(g *game, p *player) bool {
	return g.isMainPhase()
}

func (s sorcery) resolve(p *player) {}

type land struct{}

func (l land) prereq(g *game, p *player) bool {
	return g.isMainPhase() && !p.landPlayed
}

func (l land) resolve(p *player) {
	p.landPlayed = true
	p.manaTotal += 1
	p.manaAvailable += 1
	// TODO: add to players battlefield
}

type creature struct {
	power     int
	toughness int
}

func (c creature) prereq(g *game, p *player) bool {
	return g.isMainPhase()
}

func (c creature) resolve(p *player) {
	// TODO: add to players battlefield
}

// TODO: generate card data from online database
// build once card structure has stabilised a bit more
func getCard(name string) card {
	cards, err := mtg.NewQuery().Where(mtg.CardName, name).All()
	if err != nil {
		panic(err)
	}
	fmt.Println(cards)
	return card{name: cards[0].Name}
}

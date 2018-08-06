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

// sorcery/instant goes to graveyard from stack
// permanents enter the battlefield (unless countered ofc)
// NOTE: on effect -> go on stack (except lands)
// onResolve -> type specific (go to battlefield/graveyard)
// and card specific (resolve sorcery/instant effect, see game.resolve)
func (c card) resolve(p *player) {
	switch c.cardType.(type) {
	case land:
		p.landPlayed = true
		p.manaTotal += 1
		p.manaAvailable += 1
		p.battlefield.lands = append(p.battlefield.lands, cardInstance{card: c})
	case creature:
		p.battlefield.creatures = append(p.battlefield.creatures, cardInstance{card: c})
	case sorcery:
		p.graveyard = append(p.graveyard, c.name)
	}
}

type manacost struct {
	c, w, u, b, r, g int
}

func (m manacost) converted() int {
	return m.c + m.w + m.u + m.b + m.r + m.g
}

type prerequisiteFunc func(*player) bool

type cardType interface {
	prereq(*game, int) bool
}

type sorcery struct{}

func (s sorcery) prereq(g *game, pindex int) bool {
	return sorcerySpeed(g, pindex)
}

type land struct{}

func (l land) prereq(g *game, pindex int) bool {
	if !sorcerySpeed(g, pindex) {
		return false
	}
	p := g.getPlayer(pindex)
	return !p.landPlayed
}

type creature struct {
	power     int
	toughness int
}

func (c creature) prereq(g *game, pindex int) bool {
	return sorcerySpeed(g, pindex)
}

func sorcerySpeed(g *game, pindex int) bool {
	return g.isMainPhase() && len(g.stack) == 0 && g.activePlayer == pindex
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

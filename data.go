package main

import (
	"fmt"
	"strings"

	"github.com/MagicTheGathering/mtg-sdk-go"
)

type unorderedCards map[Card]int // card : amount
type orderedCards []Card         // card

var (
	mountain = &land{
		card: card{
			name: "Mountain",
		},
	}

	lavaSpike = &sorcery{
		card: card{
			name:     "Lava Spike",
			manaCost: manaCost{r: 1},
			effects: []Effect{
				playerEffect{
					effect: effect{func(p *player) {
						p.lifeTotal -= 3
					}},
				},
			},
		},
	}

	falkenrathReaver = &creature{
		card: card{
			name:     "Falkenrath Reaver",
			manaCost: manaCost{c: 1, r: 1},
		},
		power:     2,
		toughness: 2,
	}

	cards = map[string]Card{
		mountain.name:         mountain,
		lavaSpike.name:        lavaSpike,
		falkenrathReaver.name: falkenrathReaver,
	}

	deckList = unorderedCards{
		mountain:         10,
		lavaSpike:        10,
		falkenrathReaver: 10,
	}
)

// sorcery/instant goes to graveyard from stack
// permanents enter the battlefield (unless countered ofc)
// NOTE: on effect -> go on stack (except lands)
// onResolve -> type specific (go to battlefield/graveyard)
// and card specific (resolve sorcery/instant effect, see game.resolve)
type Card interface {
	prereq(*game, int) bool
	resolve(p *player)
	getName() string
	getManaCost() manaCost
	getPrereqs() []prerequisiteFunc
	getEffects() []Effect
	apply(g *game, t target)
}

type card struct {
	name     string
	manaCost manaCost
	prereqs  []prerequisiteFunc
	effects  []Effect
}

type manaCost struct {
	c, w, u, b, r, g int
}

func (m manaCost) converted() int {
	return m.c + m.w + m.u + m.b + m.r + m.g
}

type prerequisiteFunc func(*player) bool

func (c card) getName() string {
	return c.name
}

func (c card) getManaCost() manaCost {
	return c.manaCost
}

func (c card) getPrereqs() []prerequisiteFunc {
	return c.prereqs
}

func (c card) getEffects() []Effect {
	return c.effects
}

type sorcery struct {
	card
}

func (s *sorcery) prereq(g *game, pindex int) bool {
	return sorcerySpeed(g, pindex)
}

func (s *sorcery) resolve(p *player) {
	p.graveyard = append(p.graveyard, s)
}

type land struct {
	card
}

func (l *land) prereq(g *game, pindex int) bool {
	if !sorcerySpeed(g, pindex) {
		return false
	}
	p := g.getPlayer(pindex)
	return !p.landPlayed
}

func (l *land) resolve(p *player) {
	p.landPlayed = true
	p.manaTotal += 1
	p.manaAvailable += 1
	p.battlefield.lands = append(p.battlefield.lands, cardInstance{card: l})
}

type creature struct {
	card
	power     int
	toughness int
}

func (c *creature) prereq(g *game, pindex int) bool {
	return sorcerySpeed(g, pindex)
}

func (c *creature) resolve(p *player) {
	p.battlefield.creatures = append(p.battlefield.creatures, cardInstance{card: c})
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

func (c unorderedCards) String() string {
	var ss []string
	for k, v := range c {
		if v == 1 {
			ss = append(ss, k.getName())
			continue
		}
		ss = append(ss, fmt.Sprintf("%s(%d)", k.getName(), v))
	}
	return strings.Join(ss, ",")
}

package main

import (
	"fmt"
	"strings"

	"github.com/MagicTheGathering/mtg-sdk-go"
)

type unorderedCards map[Card]int // card : amount
type orderedCards []Card         // card

// sorcery/instant goes to graveyard from stack
// permanents enter the battlefield (unless countered ofc)
// NOTE: on effect -> go on stack (except lands)
// onResolve -> type specific (go to battlefield/graveyard)
// and card specific (resolve sorcery/instant effect, see game.resolve)
type Card interface {
	prereq(*game, int) bool
	resolve(g *game, a cardAction)
	getName() string
	getManaCost() manaCost
	getPrereqs() []prerequisiteFunc
}

type card struct {
	name       string
	manaCost   manaCost
	prereqs    []prerequisiteFunc
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

type sorcery struct {
	card
    spellAbility SpellAbility
}

func (s *sorcery) prereq(g *game, pindex int) bool {
	return sorcerySpeed(g, pindex)
}

func (s *sorcery) resolve(g *game, a cardAction) {
	p := g.getPlayer(a.controller)
    f := s.spellAbility.getEffect()
    f.apply(g, a.targets)
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

func (l *land) resolve(g *game, a cardAction) {
	p := g.getPlayer(a.controller)
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

func (c *creature) resolve(g *game, a cardAction) {
	p := g.getPlayer(a.controller)
	instance := cardInstance{
		card:              c,
		attacking:         -1,
		summoningSickness: true,
	}
	p.battlefield.creatures = append(p.battlefield.creatures, instance)
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

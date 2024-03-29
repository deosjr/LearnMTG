package main

import (
	"fmt"
	"math/rand"
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
	getManaCost() mana
	getPrereqs() []prerequisiteFunc
	getActivatedAbilities() []ActivatedAbility
}

type card struct {
	name     string
	manaCost mana
	prereqs  []prerequisiteFunc
	// abilities
	activatedAbilities []ActivatedAbility
	triggeredAbilities []TriggeredAbility
	staticAbilities    []StaticAbility
}

// the token of type card, i.e. a specific Mountain
type cardInstance struct {
	// unique id for targetting etc
	id                uint64
	card              Card
	tapped            bool
	summoningSickness bool
	attacking         int
}

func instanceOf(c Card) cardInstance {
	return cardInstance{
		// TODO: better rand to prevent clashes?
		id:   rand.Uint64(),
		card: c,
	}
}

type mana struct {
	c, w, u, b, r, g int
}

func (m mana) converted() int {
	return m.c + m.w + m.u + m.b + m.r + m.g
}

func (m mana) add(n mana) mana {
	return mana{c: m.c + n.c, w: m.w + n.w, u: m.u + n.u, b: m.b + n.b, r: m.r + n.r, g: m.g + n.g}
}

func (m mana) sub(n mana) mana {
	return mana{c: m.c - n.c, w: m.w - n.w, u: m.u - n.u, b: m.b - n.b, r: m.r - n.r, g: m.g - n.g}
}

func (m mana) covers(n mana) bool {
	if m.w < n.w {
		return false
	}
	m.w -= n.w
	if m.u < n.u {
		return false
	}
	m.u -= n.u
	if m.b < n.b {
		return false
	}
	m.b -= n.b
	if m.r < n.r {
		return false
	}
	m.r -= n.r
	if m.g < n.g {
		return false
	}
	m.g -= n.g
	return m.converted() >= n.c
}

type cost struct {
	mana mana
	tap  bool
	// alternative costs
}

type prerequisiteFunc func(*player) bool

func (c card) getName() string {
	return c.name
}

func (c card) getManaCost() mana {
	return c.manaCost
}

func (c card) getPrereqs() []prerequisiteFunc {
	return c.prereqs
}

func (c card) getActivatedAbilities() []ActivatedAbility {
	return c.activatedAbilities
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
	p.battlefield.lands = append(p.battlefield.lands, instanceOf(l))
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
	instance := instanceOf(c)
	instance.attacking = -1
	instance.summoningSickness = true
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

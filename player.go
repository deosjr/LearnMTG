package main

import (
	"fmt"
	"math/rand"
)

type player struct {
	name     string
	deckList unorderedCards

	lifeTotal   int
	hand        unorderedCards
	library     orderedCards
	battlefield battlefield
	graveyard   orderedCards

	landPlayed bool
	decked     bool
	// simplification for now, think hearthstone
	manaTotal     int
	manaAvailable int
}

type battlefield struct {
	lands     []cardInstance
	creatures []cardInstance
	other     []cardInstance
}

// the token of type card, i.e. a specific Mountain
type cardInstance struct {
	card Card
}

// TODO: battlefield?
func (p *player) copy() *player {
	newP := &player{}
	*newP = *p
	if len(p.hand) == 0 {
		return newP
	}
	newP.hand = unorderedCards{}
	for k, v := range p.hand {
		newP.hand[k] = v
	}
	return newP
}

func (p *player) drawN(n int) {
	for i := 0; i < n; i++ {
		p.draw()
	}
}

func (p *player) draw() {
	if len(p.library) == 0 {
		p.decked = true
		return
	}
	if p.hand == nil {
		p.hand = unorderedCards{}
	}
	card := p.library[0]
	p.hand[card] += 1
	p.library = p.library[1:]
}

func (p *player) hasMana(m manaCost) bool {
	return p.manaAvailable >= m.converted()
}

func (p *player) payMana(m manaCost) {
	p.manaAvailable -= m.converted()
}

func (p *player) String() string {
	return fmt.Sprintf("life: %d, mana: %d/%d, hand: %s", p.lifeTotal, p.manaAvailable, p.manaTotal, p.hand.String())
}

func newPlayer(name string, deckList unorderedCards) *player {
	var list orderedCards
	for k, v := range deckList {
		for i := 0; i < v; i++ {
			list = append(list, k)
		}
	}
	shuffled := make(orderedCards, len(list))
	for i, v := range rand.Perm(len(list)) {
		shuffled[v] = list[i]
	}
	return &player{
		name:      name,
		lifeTotal: 20,
		library:   shuffled,
		hand:      unorderedCards{},
		deckList:  deckList,
	}
}

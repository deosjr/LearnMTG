package main

import (
	"fmt"
	"math/rand"
)

type player struct {
	name     string
	deckList map[string]int

	lifeTotal   int
	hand        map[string]int // cardName : amount in hand
	library     []string
	battlefield []Permanent
	graveyard   []string

	landPlayed bool
	decked     bool
	// simplification for now, think hearthstone
	manaTotal     int
	manaAvailable int
}

// TODO: battlefield?
func (p *player) copy() *player {
	newP := &player{}
	*newP = *p
	newP.hand = map[string]int{}
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
	card := p.library[0]
	p.hand[card] += 1
	p.library = p.library[1:]
}

func (p *player) String() string {
	return fmt.Sprintf("life: %d, mana: %d/%d", p.lifeTotal, p.manaAvailable, p.manaTotal)
}

func newPlayer(name string, deckList map[string]int) *player {
	var list []string
	for k, v := range deckList {
		for i := 0; i < v; i++ {
			list = append(list, k)
		}
	}
	shuffled := make([]string, len(list))
	for i, v := range rand.Perm(len(list)) {
		shuffled[v] = list[i]
	}
	return &player{
		name:      name,
		lifeTotal: 20,
		library:   shuffled,
		hand:      map[string]int{},
		deckList:  deckList,
	}
}

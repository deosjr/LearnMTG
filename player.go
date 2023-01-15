package main

import (
	"fmt"
	"math/rand"
)

type player struct {
	name     string
    idx      int
	deckList unorderedCards

	lifeTotal   int
	hand        unorderedCards
	library     orderedCards
	battlefield battlefield
	graveyard   orderedCards
    manaPool    mana

	landPlayed bool
	decked     bool

    strategy Strategy
}

type battlefield struct {
	lands     []cardInstance
	creatures []cardInstance
	other     []cardInstance
}

func (b battlefield) copy() battlefield {
	return battlefield{
		lands:     copyBattlefield(b.lands),
		creatures: copyBattlefield(b.creatures),
		other:     copyBattlefield(b.other),
	}
}

func copyBattlefield(list []cardInstance) []cardInstance {
	if len(list) == 0 {
		return nil
	}
	newlist := make([]cardInstance, len(list))
	for i, ci := range list {
		newlist[i] = ci
	}
	return newlist
}

// the token of type card, i.e. a specific Mountain
type cardInstance struct {
	card              Card
	tapped            bool
	summoningSickness bool
	attacking         int
}

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
	newP.battlefield = p.battlefield.copy()
    // TODO: deep copy strategy once you keep state on it
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

func (p *player) manaAvailable() mana {
    manaAvailable := mana{}
    for _, l := range p.battlefield.lands {
        if l.tapped {
            continue
        }
        // TODO: assume land only has one activated ability
        // and that is a mana ability
        a := l.card.getActivatedAbilities()[0]
        if !a.isManaAbility() {
            panic("broken assumption on land abilities")
        }
        m := a.getEffect().(addMana).amount
        manaAvailable = manaAvailable.add(m)
    }
    return manaAvailable
}

func (p *player) hasMana(m mana) bool {
	return p.manaAvailable().covers(m)
}

func (p *player) String() string {
	return fmt.Sprintf("life: %d, mana: %d/%d, hand: %s", p.lifeTotal, p.manaAvailable(), len(p.battlefield.lands), p.hand.String())
}

func newPlayer(idx int, name string, deckList unorderedCards) *player {
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
        idx:       idx,
		lifeTotal: 20,
		library:   shuffled,
		hand:      unorderedCards{},
		deckList:  deckList,
	}
}

func (p *player) creaturesThatCanAttack() []int {
	creatures := []int{}
	for i, c := range p.battlefield.creatures {
		if c.tapped || c.summoningSickness {
			continue
		}
		creatures = append(creatures, i)
	}
    return creatures
}

// this means we only check prereqs against what we know
// may have to change that to a probability prereq is met
func (p *player) canPlayCard(g *game, card Card) bool {
	// prerequisites given by card type
	if !card.prereq(g, p.idx) {
		return false
	}

	// can player pay for the card?
	if !p.hasMana(card.getManaCost()) {
		return false
	}
	// other prerequisites such as paying life
	// NOTE: prereq is target available?
	// --> this is handled by possibleTargets returning 0 actions
	for _, prereq := range card.getPrereqs() {
		if !prereq(p) {
			return false
		}
	}
	return true
}

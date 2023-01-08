package main

// a player has a strategy they follow, their AI (or human-controlled) behaviour

type Strategy interface {
    NextAction(*player, *game) Action
    Attacks(*player, *game) attackAction
}

// your goldfish can't play magic, so it always just passes
type goldfish struct {}

func (goldfish) NextAction(_ *player, g *game) Action {
    pIndex := g.priorityPlayer
    return passAction{action{controller: pIndex}}
}

func (goldfish) Attacks(_ *player, g *game) attackAction {
    pIndex := g.priorityPlayer
	return attackAction{action: action{controller: pIndex}, attackers: nil}
}

// TODO: a simpler strategy hardcoding the simple deck we have
// never do anything first main phase.
// always attack with everything, never block
// second main phase, always play a land first
// always play a creature if you can
// otherwise, always play lava spike face
// pass in every other step ever
type simpleStrategy struct {}

func (simpleStrategy) NextAction(p *player, g *game) Action {
    pIndex := g.priorityPlayer
    if g.currentStep != postcombatMainPhase {
        return passAction{action{controller: pIndex}}
    }
    for c := range p.hand {
        if _, ok := c.(*land); !ok {
            continue
        }
        if !g.canPlayCard(pIndex, c) {
            continue
        }
        return cardAction{card: c, action: action{controller: pIndex}}
    }
    for c := range p.hand {
        if _, ok := c.(*creature); !ok {
            continue
        }
        if !g.canPlayCard(pIndex, c) {
            continue
        }
        return cardAction{card: c, action: action{controller: pIndex}}
    }
    for c := range p.hand {
        s, ok := c.(*sorcery)
        if !ok {
            continue
        }
        if s.name != "Lava Spike" {
            continue
        }
        if !g.canPlayCard(pIndex, c) {
            continue
        }
        return cardAction{card: c, action: action{controller: pIndex}, targets: []effectTarget{{index:target((pIndex+1)%2),ttype:targetPlayer}}}
    }
    for c := range p.hand {
        s, ok := c.(*sorcery)
        if !ok {
            continue
        }
        if s.name != "Flame Rift" {
            continue
        }
        if !g.canPlayCard(pIndex, c) {
            continue
        }
        if p.lifeTotal <= 4 || p.lifeTotal < g.getOpponent(pIndex).lifeTotal {
            continue
        }
        return cardAction{card: c, action: action{controller: pIndex}, targets: []effectTarget{{ttype:eachPlayer}}}
    }
    return passAction{action{controller: pIndex}}
}

func (simpleStrategy) Attacks(p *player, g *game) attackAction {
    pIndex := g.priorityPlayer
    return attackWithAll(p, pIndex)
}

func attackWithAll(p *player, index int) attackAction {
    creatures := p.creaturesThatCanAttack()
	// two player assumption
	opp := (index + 1) % 2
	attackers := []combatTarget{}
	for _, c := range creatures {
		attackers = append(attackers, combatTarget{index: c, target: opp})
	}
	return attackAction{action: action{controller: index}, attackers: attackers}
}

package main

// a player has a strategy they follow, their AI (or human-controlled) behaviour

type Strategy interface {
    NextAction(*player, *game) Action
    Attacks(*player, *game) attackAction
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
        return cardAction{card: c, action: action{controller: pIndex}, targets: []target{{index:0, target:(pIndex+1)%2}}}
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
	attackers := []target{}
	for _, c := range creatures {
		attackers = append(attackers, target{index: c, target: opp})
	}
	return attackAction{action: action{controller: index}, attackers: attackers}
}

package main

// a player has a strategy they follow, their AI (or human-controlled) behaviour

type Strategy interface {
    NextAction(*player, *game) Action
    Attackers(*player, *game) []cardInstance
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
    if g.currentStep == declareAttackersStep && g.declarations == 0 {
        return g.getAttacks(pIndex)[0]
    }
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

func (simpleStrategy) Attackers(p *player, _ *game) []cardInstance {
    return p.battlefield.creatures
}

// legacy: minmax is probably not feasible to use
type minmaxStrategy struct {}

func (minmaxStrategy) NextAction(_ *player, g *game) Action {
    return startMinimax(g)
}

// attack with everything that can attack
func (minmaxStrategy) Attackers(p *player, _ *game) []cardInstance {
    return p.battlefield.creatures
}

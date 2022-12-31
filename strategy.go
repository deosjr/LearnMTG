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
    if g.currentStep == declareAttackersStep && g.declarations == 0 {
        return g.getAttacks(g.activePlayer)[0]
    }
    return passAction{action{controller: g.activePlayer}}
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

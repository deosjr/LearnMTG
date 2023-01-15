package main

// a player has a strategy they follow, their AI (or human-controlled) behaviour

type Strategy interface {
    NextAction(*player, *game) Action
    Attacks(*player, *game) attackAction
    PayManaCost(p *player, cost mana)
}

// your goldfish can't play magic, so it always just passes
type goldfish struct {}

func (goldfish) NextAction(p *player, g *game) Action {
    return passAction{action{controller: p.idx}}
}

func (goldfish) Attacks(p *player, g *game) attackAction {
	return attackAction{action: action{controller: p.idx}, attackers: nil}
}

func (goldfish) PayManaCost(p *player, cost mana) {
    payNaive(p, cost)
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
    if g.currentStep != postcombatMainPhase {
        return passAction{action{controller: p.idx}}
    }
    for c := range p.hand {
        if _, ok := c.(*land); !ok {
            continue
        }
        if !p.canPlayCard(g, c) {
            continue
        }
        return cardAction{card: c, action: action{controller: p.idx}}
    }
    for c := range p.hand {
        if _, ok := c.(*creature); !ok {
            continue
        }
        if !p.canPlayCard(g, c) {
            continue
        }
        return cardAction{card: c, action: action{controller: p.idx}}
    }
    for c := range p.hand {
        s, ok := c.(*sorcery)
        if !ok {
            continue
        }
        if s.name != "Lava Spike" {
            continue
        }
        if !p.canPlayCard(g, c) {
            continue
        }
        return cardAction{card: c, action: action{controller: p.idx}, targets: []effectTarget{{index:target((p.idx+1)%2),ttype:targetPlayer}}}
    }
    for c := range p.hand {
        s, ok := c.(*sorcery)
        if !ok {
            continue
        }
        if s.name != "Flame Rift" {
            continue
        }
        if !p.canPlayCard(g, c) {
            continue
        }
        if p.lifeTotal <= 4 || p.lifeTotal < g.getOpponent(p.idx).lifeTotal {
            continue
        }
        return cardAction{card: c, action: action{controller: p.idx}, targets: []effectTarget{{ttype:eachPlayer}}}
    }
    return passAction{action{controller: p.idx}}
}

func (simpleStrategy) Attacks(p *player, g *game) attackAction {
    return attackWithAll(p, p.idx)
}

func (simpleStrategy) PayManaCost(p *player, cost mana) {
    payNaive(p, cost)
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

// assumption: player has the mana to pay
func payNaive(p *player, cost mana) {
    needed := cost.converted()
    for i, l := range p.battlefield.lands {
        if needed == 0 {
            break
        }
        if l.tapped {
            continue
        }
        // TODO: assume land only has one activated ability
        // and that is a mana ability
        a := l.card.getActivatedAbilities()[0]
        if !a.isManaAbility() {
            panic("broken assumption on land abilities")
        }
        // TODO: find/pay colored mana first, then spend rest to pay colorless
        //m := a.getEffect().(addMana).amount
        l.tapped = true
        p.battlefield.lands[i] = l
        needed--
    }
}

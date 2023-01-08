package main

import (
	"math"
)

// Minimax algorithm for Magic the Gathering
// plies are turn segments where the player holds priority

// legacy: minmax is probably not feasible to use
type minmaxStrategy struct {}

func (minmaxStrategy) NextAction(_ *player, g *game) Action {
    return startMinimax(g)
}

func (minmaxStrategy) Attacks(p *player, g *game) attackAction {
    return startMinimax(g).(attackAction)
}

var maxDepth = 20

type node struct {
	game        *game
	pointOfView int
}

func startMinimax(g *game) Action {
	root := node{game: g, pointOfView: g.priorityPlayer}
	var a Action
	bestValue := -math.MaxFloat64
	for _, childAction := range root.getActionsSelf() {
		child := root.getChild(childAction)
		v := minimax(child, maxDepth)
		if v > bestValue {
			bestValue = v
			a = childAction
		}
	}
	return a
}

func minimax(node node, depth int) float64 {
	if depth == 0 || node.isTerminal() {
		score := node.evaluate(depth)
		return score
	}

	if node.maximizing() {
		bestValue := -math.MaxFloat64
		for _, childAction := range node.getActionsSelf() {
			child := node.getChild(childAction)
			v := minimax(child, depth-1)
			bestValue = math.Max(bestValue, v)
		}
		return bestValue
	}

	bestValue := math.MaxFloat64
	for _, childAction := range node.getActionsOpponent() {
		child := node.getChild(childAction)
		v := minimax(child, depth-1)
		bestValue = math.Min(bestValue, v)
	}
	return bestValue
}

func (n node) maximizing() bool {
	return n.pointOfView == n.game.priorityPlayer
}

func (n node) getChild(action Action) node {
	g := n.game.copy()
	g.resolveAction(action)
	return node{
		game:        g,
		pointOfView: n.pointOfView,
	}
}

func (n node) getActionsSelf() []Action {
	return getActions(n.game, n.pointOfView)
}

func (n node) getActionsOpponent() []Action {
	// again, two player game assumption for now
	oppIndex := (n.pointOfView + 1) % 2
	return getActions(n.game, oppIndex)
}

// use an arbitrarily large number because I
// don't want to calculate using actual MaxFloat64
const infinity float64 = 1000000.0

// evaluate payoff function
func (n node) evaluate(depth int) float64 {
	p := n.game.getPlayer(n.pointOfView)
	opp := n.game.getOpponent(n.pointOfView)
	if opp.lifeTotal <= 0 || opp.decked {
		// penalise long term plans: winning earlier is better!
		return infinity - float64(-depth)
	}
	if p.lifeTotal <= 0 || p.decked {
		return -infinity
	}

	// TODO: weights per feature
	lifeDiff := float64(p.lifeTotal - opp.lifeTotal)
	return lifeDiff + float64(p.manaTotal) - float64(-depth)
}

//does the game end in this configuration next state-based check?
func (n node) isTerminal() bool {
	// NOTE: needs to only consider game ending state based actions
	return n.game.checkStateBasedActions()
}

func getActions(g *game, index int) []Action {
	actions := []Action{passAction{action{controller: index}}}
	if g.currentStep == declareAttackersStep && g.declarations == 0 {
		return getAttacks(g, index)
	}
	p := g.getPlayer(index)
	for card, _ := range p.hand {
		if !g.canPlayCard(index, card) {
			continue
		}
        switch c := card.(type) {
        case *sorcery:
            for _, t := range getTargets(g, c.spellAbility, index) {
		        actions = append(actions, cardAction{card: card, action: action{controller: index}, targets: t})
            }
        default:
            actions = append(actions, cardAction{card: card, action: action{controller: index}})
        }
	}
	return actions
}

func getAttacks(g *game, index int) []Action {
	// TODO: first attempt, always attack with everything
    // for minimax, this should return the superset of attackers instead
    p := g.getPlayer(index)
    return []Action{attackWithAll(p, index)}
}

func possibleTargets(g *game, t targettype, controller int) []target {
    switch t {
    case selfTarget:
        return []target{ target(controller) }
    case playerTarget:
        ts := []target{}
        for i:=0; i < g.numPlayers; i++ {
            ts = append(ts, target(i))
        }
        return ts
    }
    return nil
}

func getTargets(g *game, a Ability, controller int) [][]target {
    // TODO: multiple targets
    targets := [][]target{possibleTargets(g, a.getTargets()[0], controller)}
	// Multi-target combinatorics (TODO: with constraints such as no same target!)
	if len(targets) == 0 {
		return nil
	}
	oldArr := [][]target{}
	for _, t := range targets[0] {
		oldArr = append(oldArr, []target{t})
	}
	newArr := [][]target{}
	for _, ta := range targets[1:] {
		for _, t := range ta {
			for _, ot := range oldArr {
				newArr = append(newArr, append(ot, t))
			}
		}
		oldArr = newArr
		newArr = [][]target{}
	}
    return oldArr
}

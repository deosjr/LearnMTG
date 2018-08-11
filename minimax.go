package main

import (
	"math"
)

var maxDepth = 10

type node struct {
	game        *game
	pointOfView int
}

func startMinimax(g *game) action {
	root := node{game: g, pointOfView: g.priorityPlayer}
	var a action
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

func (n node) getChild(action action) node {
	g := n.game.copy()
	g.resolveAction(action)
	return node{
		game:        g,
		pointOfView: n.pointOfView,
	}
}

var pass = "PASS"

// TODO: split out in play, ability, special, attack and block actions
type action struct {
	card       string
	controller int
	effects    []effect
}

func (n node) getActionsSelf() []action {
	return n.game.getActions(n.pointOfView)
}

func (n node) getActionsOpponent() []action {
	// again, two player game assumption for now
	oppIndex := (n.pointOfView + 1) % 2
	return n.game.getActions(oppIndex)
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

package main

import (
	"math"
)

type node struct {
	pointOfView  *player
	isActive     bool
	actionsTaken []action
}

func minimax(node node, depth int, maximizingPlayer bool) float64 {
	p := node.simulate()
	if depth == 0 || node.isTerminal(p) {
		score := node.evaluate(p)
		return score
	}

	if maximizingPlayer {
		bestValue := -math.MaxFloat64
		for _, child := range node.getChildren(node.pointOfView) {
			v := minimax(child, depth-1, false)
			bestValue = math.Max(bestValue, v)
		}
		return bestValue
	}
	// minimizing player
	bestValue := math.MaxFloat64
	for _, child := range node.getChildren(node.pointOfView.opponent) {
		v := minimax(child, depth-1, true)
		bestValue = math.Min(bestValue, v)
	}
	return bestValue
}

// we evaluate a node with children, not a complete tree
// per child we need to build the tree one layer further
func (n node) getChildren(p *player) []node {
	var actions []action
	if p == n.pointOfView {
		actions = n.getActionsSelf()
	} else {
		actions = n.getActionsOpponent()
	}

	children := make([]node, len(actions))
	for i, a := range actions {
		newActions := append(n.actionsTaken, a)
		children[i] = node{
			pointOfView:  n.pointOfView,
			actionsTaken: newActions,
			isActive:     !n.isActive,
		}
	}
	return children
}

func (n node) getActionsSelf() []action {
	passAction := action{card: pass, player: n.pointOfView.name}
	actions := []action{passAction}
	if !n.isActive {
		return actions
	}
Loop:
	for k, _ := range n.pointOfView.hand {
		card := cards[k]
		for _, prereq := range card.prereqs {
			if !prereq(n.pointOfView) {
				continue Loop
			}
		}
		action := action{card: card.name, player: n.pointOfView.name}
		actions = append(actions, action)
	}
	return actions
}

// TODO: deal with incomplete information!
func (n node) getActionsOpponent() []action {
	passAction := action{card: pass, player: n.pointOfView.opponent.name}
	actions := []action{passAction}
	if !n.isActive || len(n.pointOfView.opponent.hand) == 0 {
		return actions
	}
	// worst-case assumption: player always has all possible cards
	// extra prereq: they have at least 1 card in hand
Loop:
	for k, _ := range n.pointOfView.opponent.deckList {
		card := cards[k]
		for _, prereq := range card.prereqs {
			// this means we only check prereqs against what we know
			// may have to change that to a probability prereq is met
			if !prereq(n.pointOfView.opponent) {
				continue Loop
			}
		}
		action := action{card: card.name, player: n.pointOfView.opponent.name}
		actions = append(actions, action)
	}
	return actions
}

// reproduce the current board state by replaying actions
// TODO: simulate going through phases/steps
func (n node) simulate() *player {
	// copy player and opponent
	p, opp := n.pointOfView.copy(), n.pointOfView.opponent.copy()
	p.opponent, opp.opponent = opp, p

	for _, action := range n.actionsTaken {
		if action.card == pass {
			continue
		}
		actingPlayer := p
		if action.player == opp.name {
			actingPlayer = opp
		}
		action.execute(actingPlayer)
	}
	return p
}

// use an arbitrarily large number because I
// don't want to calculate using actual MaxFloat64
const infinity float64 = 1000000.0

// evaluate payoff function
func (n node) evaluate(p *player) float64 {
	opp := p.opponent
	if opp.lifeTotal <= 0 || opp.decked {
		// penalise lang term plans: winning earlier is better!
		return infinity - float64(len(n.actionsTaken))
	}
	if p.lifeTotal <= 0 || p.decked {
		return -infinity
	}

	// TODO: weights per feature
	lifeDiff := float64(p.lifeTotal - opp.lifeTotal)
	return lifeDiff + float64(p.manaTotal) - float64(len(n.actionsTaken))
}

//does the game end in this configuration next state-based check?
func (n node) isTerminal(p *player) bool {
	// NOTE: needs to only consider game ending state based actions
	return checkStateBasedActions(p)
}

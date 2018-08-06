package main

import (
	"math"
)

type node struct {
	game        *game
	pointOfView int
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

// TODO: split out in play, ability, special, attack and defend actions
type action struct {
	card       string
	controller int
	effects    []effect
}

func (n node) getActionsSelf() []action {
	actions := []action{action{card: pass, controller: n.pointOfView}}
	p := n.game.getPlayer(n.pointOfView)
	for k, _ := range p.hand {
		card := cards[k]
		if !n.game.canPlayCard(n.pointOfView, card) {
			continue
		}
		actions = append(actions, getCardActions(card, n.pointOfView, n.game)...)
	}
	return actions
}

// TODO: deal with incomplete information!
func (n node) getActionsOpponent() []action {
	// again, two player game assumption for now
	oppIndex := (n.pointOfView + 1) % 2
	opp := n.game.getPlayer(oppIndex)
	actions := []action{action{card: pass, controller: oppIndex}}
	// worst-case assumption: player always has all possible cards
	// extra prereq: they have at least 1 card in hand
	if len(opp.hand) == 0 {
		return actions
	}
	for k, _ := range opp.deckList {
		card := cards[k]
		if !n.game.canPlayCard(oppIndex, card) {
			continue
		}
		actions = append(actions, getCardActions(card, oppIndex, n.game)...)
	}
	return actions
}

func getCardActions(card card, controller int, game *game) []action {
	if card.effects == nil {
		return []action{action{card: card.name, controller: controller}}
	}
	actions := []action{}
	effects := [][]effect{}
	for _, e := range card.effects {
		effects = append(effects, e.possibleTargets(controller, game))
	}
	// Multi-target combinatorics (TODO: with constraints such as no same target!)
	if len(effects) == 0 {
		return actions
	}
	oldArr := [][]effect{}
	for _, e := range effects[0] {
		oldArr = append(oldArr, []effect{e})
	}
	newArr := [][]effect{}
	for _, ea := range effects[1:] {
		for _, e := range ea {
			for _, oe := range oldArr {
				newArr = append(newArr, append(oe, e))
			}
		}
		oldArr = newArr
		newArr = [][]effect{}
	}
	for _, e := range oldArr {
		action := action{card: card.name, controller: controller, effects: e}
		actions = append(actions, action)
	}
	return actions
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

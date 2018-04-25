package main

import (
	"fmt"
	"math"
	"math/rand"
)

type step int

const (
	// beginning phase
	untapStep step = iota
	upkeepStep
	drawStep

	precombatMainPhase

	// combat phase
	beginningOfCombatStep
	declareAttackersStep
	declareBlockersStep
	combatDamageFirstStrikeStep
	combatDamageStep
	endOfCombatStep

	postcombatMainPhase

	// ending phase
	endStep
	cleanupStep

	numSteps // so we can do step + 1 % numSteps for turn wrap
)

type game struct {
	players        []*player
	currentStep    step
	turn           int
	activePlayer   int
	priorityPlayer int
	startingPlayer int
	numPlayers     int
}

func newGame(players []*player) *game {
	numPlayers := len(players)
	startingPlayer := rand.Intn(numPlayers)
	return &game{
		players:        players,
		currentStep:    precombatMainPhase,
		turn:           1,
		activePlayer:   startingPlayer,
		priorityPlayer: startingPlayer,
		startingPlayer: startingPlayer,
		numPlayers:     numPlayers,
	}
}

func (g *game) getPlayer(i int) *player {
	if len(g.players) <= i {
		panic(fmt.Sprintf("invalid index %d", i))
	}
	return g.players[i]
}

// Two player game for now, getOpponents() later
func (g *game) getOpponent(i int) *player {
	return g.getPlayer((i + 1) % 2)
}

func (g *game) getActivePlayer() *player {
	return g.getPlayer(g.activePlayer)
}

func (g *game) getPriorityPlayer() *player {
	return g.getPlayer(g.priorityPlayer)
}

func (g *game) playUntilPriority() {
	activePlayer := g.getActivePlayer()
	for {
		switch g.currentStep {
		case untapStep:
			activePlayer.manaAvailable = activePlayer.manaTotal
		case upkeepStep:
			break //skip
		case drawStep:
			activePlayer.draw()
		case precombatMainPhase:
			return
		case beginningOfCombatStep:
			break //skip
		case declareAttackersStep:
			break //skip
		case declareBlockersStep:
			break //skip
		case combatDamageFirstStrikeStep:
			break //skip
		case combatDamageStep:
			break //skip
		case endOfCombatStep:
			break //skip
		case postcombatMainPhase:
			break //skip
		case endStep:
			break //skip
		case cleanupStep:
			break //skip
		}
		// just passed past the cleanup into next turn
		if g.currentStep == cleanupStep {
			activePlayer.landPlayed = false
			// NOTE: priority only switches on turn (i.e. no instants)
			g.activePlayer = (g.activePlayer + 1) % g.numPlayers
			g.priorityPlayer = (g.priorityPlayer + 1) % g.numPlayers
			if g.activePlayer == g.startingPlayer {
				g.turn++
			}
			activePlayer = g.getActivePlayer()
		}
		g.nextStep()
	}
	return
}

func (g *game) nextStep() {
	g.currentStep = (g.currentStep + 1) % numSteps
}

func (g *game) checkStateBasedActions() (gameEnds bool) {
	for _, p := range g.players {
		if p.lifeTotal <= 0 || p.decked {
			return true
		}
	}
	return false
}

func (g *game) copy() *game {
	newG := &game{}
	*newG = *g
	newG.players = make([]*player, len(g.players))
	for i, p := range g.players {
		newG.players[i] = p.copy()
	}
	return newG
}

func (g *game) debug() {
	activePlayer := g.getActivePlayer()
	opp := g.getOpponent(g.activePlayer)
	fmt.Printf("%s turn %d: %s VS %s \n", activePlayer.name, g.turn, activePlayer.String(), opp.String())
}

var maxDepth = 10

func (g *game) getPlayerAction() action {
	index := g.priorityPlayer
	root := node{game: g, pointOfView: index}
	childActions := root.getActionsSelf()
	var a action
	bestValue := -math.MaxFloat64
	for _, childAction := range childActions {
		child := root.getChild(childAction)
		v := minimax(child, maxDepth, false)
		if v > bestValue {
			bestValue = v
			a = childAction
		}
	}
	if a.card == "" {
		// TODO: all options lead to certain death
		return action{card: pass, playerSelf: index}
	}
	return a
}

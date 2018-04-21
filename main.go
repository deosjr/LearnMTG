package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Minimax algorithm for Magic the Gathering
// plies are turn segments where the player holds priority

var pass = "PASS"

// divide per target type?
type action struct {
	card   string
	player string
}

func (a action) execute(p *player) {
	if a.card == pass {
		return
	}
	p.hand[a.card] -= 1
	if p.hand[a.card] == 0 {
		delete(p.hand, a.card)
	}
	card := cards[a.card]
	card.effect(p)
}

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

// lets start with a very simple form of magic:
// only stop is on first main phase
// there are only sorcery and basic land cards

func main() {
	rand.Seed(time.Now().UnixNano())

	p1 := newPlayer("player1", deckList)
	p2 := newPlayer("player2", deckList)
	p1.opponent, p2.opponent = p2, p1
	p1.drawN(7)
	p2.drawN(7)

	activePlayer := p1
	if rand.Intn(2) == 0 {
		activePlayer = p2
	}

	currentStep := precombatMainPhase

	i := 0
	turn := 1
	//opponentPassed := false

	for {
		playerToAct, newStep, newTurn := playUntilPriority(activePlayer, currentStep)
		currentStep = newStep
		if newTurn {
			i++
			fmt.Printf("%s turn %d: %s VS %s \n", activePlayer.name, turn, activePlayer.String(), activePlayer.opponent.String())
			turn = i/2 + 1
			activePlayer = activePlayer.opponent
		}
		gameEnds := checkStateBasedActions(activePlayer)
		if gameEnds {
			fmt.Println("End of game")
			fmt.Printf("%s turn %d: %s VS %s \n", activePlayer.name, turn, activePlayer.String(), activePlayer.opponent.String())
			break
		}
		action := playerToAct.act(activePlayer == playerToAct, currentStep)
		if action.card == pass {
			// if opponentPassed {
			currentStep = (currentStep + 1) % numSteps
			// }
			//opponentPassed = true
			//} else {
			//	opponentPassed = false
		}
		action.execute(playerToAct)
		fmt.Printf("-> %s played %s \n", action.player, action.card)
	}
}

func playUntilPriority(activePlayer *player, currentStep step) (playerToAct *player, newStep step, newTurn bool) {
	for {
		switch currentStep {
		case untapStep:
			activePlayer.manaAvailable = activePlayer.manaTotal
		case upkeepStep:
			break //skip
		case drawStep:
			activePlayer.draw()
		case precombatMainPhase:
			// notice there is no priority passing here yet
			return activePlayer, currentStep, newTurn
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
		if currentStep == cleanupStep {
			activePlayer.landPlayed = false
			activePlayer = activePlayer.opponent
			newTurn = true
		}
		currentStep = (currentStep + 1) % numSteps
	}
	return activePlayer, currentStep, newTurn
}

func checkStateBasedActions(p *player) (gameEnds bool) {
	return p.lifeTotal <= 0 || p.opponent.lifeTotal <= 0 || p.decked || p.opponent.decked
}

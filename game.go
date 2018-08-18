package main

import (
	"fmt"
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
	stack          []cardAction
	currentStep    step
	turn           int
	activePlayer   int
	priorityPlayer int
	startingPlayer int
	numPlayers     int
	numPasses      int
	// some phases have a number of action points for
	// one or both players involved in combat for example.
	// we need to track how far along the phase we are here.
	// decided NOT to split in subphases for clarity later on.
	declarations int
}

func newGame(players ...*player) *game {
	numPlayers := len(players)
	startingPlayer := rand.Intn(numPlayers)
	fmt.Printf("Starting player: %s\n", players[startingPlayer].name)
	for _, p := range players {
		p.drawN(7)
	}
	g := &game{
		players:        players,
		currentStep:    precombatMainPhase,
		turn:           1,
		activePlayer:   startingPlayer,
		priorityPlayer: startingPlayer,
		startingPlayer: startingPlayer,
		numPlayers:     numPlayers,
	}
	g.nextDecisionPoint()
	return g
}

func (g *game) loop() {
	for {
		g.debug()
		a := g.getPlayerAction()
		var ac cardAction
		var stacklength int
		if len(g.stack) != 0 {
			stacklength = len(g.stack)
			ac = g.stack[stacklength-1]
		}
		g.resolveAction(a)
		switch at := a.(type) {
		case passAction:
			fmt.Printf("-> %s passes\n", g.getPlayer(a.getController()).name)
			if len(g.stack) < stacklength {
				// ac resolved
				for _, target := range ac.targets {
					fmt.Printf("%s resolved by %s targeting %s \n", ac.card.getName(), g.getPlayer(ac.controller).name, g.getPlayer(target.target).name)
				}
			}
		case cardAction:
			fmt.Printf("-> %s played %s", g.getPlayer(at.controller).name, at.card.getName())
			if len(at.targets) > 0 {
				fmt.Printf(" targeting %s", g.getPlayer(at.targets[0].target).name)
			}
			fmt.Println()
		}
		if gameEnds := g.checkStateBasedActions(); gameEnds {
			g.debug()
			fmt.Println("End of game")
			return
		}
	}
}

func (g *game) resolveAction(action Action) {
	switch a := action.(type) {
	case passAction:
		g.numPasses++
		// 116.3d If a player has priority and chooses not to take any actions,
		// that player passes. [...] Then the next player in turn order receives priority.
		g.advancePriority()
		// 116.4. If all players pass in succession
		// (that is, if all players pass without taking any actions in between passing),
		// the spell or ability on top of the stack resolves or, if the stack is empty,
		// the phase or step ends.
		if g.numPasses == g.numPlayers {
			g.numPasses = 0
			if len(g.stack) != 0 {
				g.resolve()
			} else {
				g.nextStep()
				g.nextDecisionPoint()
			}
			// 116.3a The active player receives priority at the beginning of most steps and phases [...]
			// 116.3b The active player receives priority after a spell or ability (other than a mana ability) resolves.
			g.priorityPlayer = g.activePlayer
		}
	case cardAction:
		// 116.3c If a player has priority when they cast a spell,
		// activate an ability, or take a special action, that player receives priority afterward.
		g.numPasses = 0
		g.play(a)
		// TODO (currently a hack): special actions (such as playing land)
		// do not always pass priority to the other player
		if a.card == mountain {
			g.resolve()
		}
	case attackAction:
		g.declarations += 1
	case blockAction:
		g.declarations += 1
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

func (g *game) advancePriority() {
	g.priorityPlayer = (g.priorityPlayer + 1) % g.numPlayers
}

func (g *game) nextDecisionPoint() {
	for {
		switch g.currentStep {
		case untapStep:
			activePlayer := g.getActivePlayer()
			activePlayer.manaAvailable = activePlayer.manaTotal
		case upkeepStep:
			break //skip
		case drawStep:
			g.getActivePlayer().draw()
		case precombatMainPhase:
			return
		case beginningOfCombatStep:
			break //skip
		case declareAttackersStep:
			switch g.declarations {
			// active player declares attackers
			case 0:
				return
			// attackers have already been declared, continue this step
			case 1:
				break
			}
			// triggered abilities that trigger off attackers being declared trigger
			return
		case declareBlockersStep:
			// if no attackers, skip
			return
		case combatDamageFirstStrikeStep:
			// if no attackers, skip
			break //skip
		case combatDamageStep:
			// if no attackers, skip
			return
		case endOfCombatStep:
			return
		case postcombatMainPhase:
			return
		case endStep:
			break //skip
		case cleanupStep:
			break //skip
		}
		// just passed past the cleanup into next turn
		if g.currentStep == cleanupStep {
			g.nextTurn()
		}
		g.nextStep()
	}
	return
}

// this means we only check prereqs against what we know
// may have to change that to a probability prereq is met
func (g *game) canPlayCard(pindex int, card Card) bool {
	// prerequisites given by card type
	if !card.prereq(g, pindex) {
		return false
	}

	p := g.getPlayer(pindex)
	// can player pay for the card?
	if !p.hasMana(card.getManaCost()) {
		return false
	}
	// other prerequisites such as paying life
	// NOTE: prereq is target available?
	// --> this is handled by possibleTargets returning 0 actions
	for _, prereq := range card.getPrereqs() {
		if !prereq(p) {
			return false
		}
	}
	return true
}

func (g *game) nextStep() {
	g.declarations = 0
	g.currentStep = (g.currentStep + 1) % numSteps
}

func (g *game) nextTurn() {
	g.getActivePlayer().landPlayed = false
	g.activePlayer = (g.activePlayer + 1) % g.numPlayers
	// TODO: should check statebased actions here too!
	g.advancePriority()
	if g.activePlayer == g.startingPlayer {
		g.turn++
	}
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
	if len(g.stack) == 0 {
		return newG
	}
	newG.stack = make([]cardAction, len(g.stack))
	for i, a := range g.stack {
		newG.stack[i] = a
	}
	return newG
}

func (g *game) debug() {
	activePlayer := g.getActivePlayer()
	opp := g.getOpponent(g.activePlayer)
	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("%s turn %d step %d: %s \n", activePlayer.name, g.turn, g.currentStep, activePlayer.String())
	fmt.Printf("           VS %s: %s \n", opp.name, opp.String())
}

func (g *game) play(a cardAction) {
	p := g.getPlayer(a.controller)

	// remove card from players hand
	p.hand[a.card] -= 1
	if p.hand[a.card] == 0 {
		delete(p.hand, a.card)
	}

	p.payMana(a.card.getManaCost())

	g.stack = append(g.stack, a)
}

func (g *game) resolve() {
	if len(g.stack) == 0 {
		panic("no stack to resolve")
	}
	a := g.stack[len(g.stack)-1]
	g.stack = g.stack[:len(g.stack)-1]

	p := g.getPlayer(a.controller)
	a.card.resolve(p)

	// card specific effects
	for _, t := range a.targets {
		a.card.apply(g, t)
	}
}

func (g *game) isMainPhase() bool {
	return g.currentStep == precombatMainPhase || g.currentStep == postcombatMainPhase
}

func (g *game) getPlayerAction() Action {
	return startMinimax(g)
}

func (g *game) getActions(index int) []Action {
	p := g.getPlayer(index)
	actions := []Action{passAction{action{controller: index}}}
	for card, _ := range p.hand {
		if !g.canPlayCard(index, card) {
			continue
		}
		actions = append(actions, g.getCardActions(card, index)...)
	}
	return actions
}

func (g *game) getCardActions(card Card, controller int) []Action {
	if card.getEffects() == nil {
		return []Action{cardAction{card: card, action: action{controller: controller}}}
	}
	actions := []Action{}
	targets := [][]target{}
	for i, e := range card.getEffects() {
		targets = append(targets, e.possibleTargets(i, controller, g))
	}
	// Multi-target combinatorics (TODO: with constraints such as no same target!)
	if len(targets) == 0 {
		return actions
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
	for _, t := range oldArr {
		action := cardAction{card: card, action: action{controller: controller}, targets: t}
		actions = append(actions, action)
	}
	return actions
}

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
	stack          []action
	currentStep    step
	turn           int
	activePlayer   int
	priorityPlayer int
	startingPlayer int
	numPlayers     int
	numPasses      int
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
	g.advanceUntilPriority()
	return g
}

func (g *game) loop() {
	for {
		g.debug()
		a := g.getPlayerAction()
		var ac action
		var stacklength int
		if len(g.stack) != 0 {
			stacklength = len(g.stack)
			ac = g.stack[stacklength-1]
		}
		g.resolveAction(a)
		if a.card == pass {
			fmt.Printf("-> %s passes\n", g.getPlayer(a.controller).name)
			if len(g.stack) < stacklength {
				// ac resolved
				for _, eff := range ac.effects {
					e := eff.(playerEffect)
					fmt.Printf("%s resolved by %s targeting %s \n", ac.card, g.getPlayer(ac.controller).name, g.getPlayer(e.target).name)
				}
			}
		} else {
			fmt.Printf("-> %s played %s", g.getPlayer(a.controller).name, a.card)
			if len(a.effects) > 0 {
				e := a.effects[0].(playerEffect)
				fmt.Printf(" targeting %s", g.getPlayer(e.target).name)
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

func (g *game) resolveAction(action action) {
	if action.card == pass {
		g.numPasses++
		// 116.3d If a player has priority and chooses not to take any actions,
		// that player passes. [...] Then the next player in turn order receives priority.
		g.advancePriority()
	} else {
		// 116.3c If a player has priority when they cast a spell,
		// activate an ability, or take a special action, that player receives priority afterward.
		g.numPasses = 0
		g.play(action)
		// TODO (currently a hack): special actions (such as playing land)
		// do not always pass priority to the other player
		if action.card == "Mountain" {
			g.resolve()
		}
	}

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
			g.advanceUntilPriority()
		}
		// 116.3a The active player receives priority at the beginning of most steps and phases [...]
		// 116.3b The active player receives priority after a spell or ability (other than a mana ability) resolves.
		g.priorityPlayer = g.activePlayer
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

func (g *game) advanceUntilPriority() {
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
			g.nextTurn()
		}
		g.nextStep()
	}
	return
}

// this means we only check prereqs against what we know
// may have to change that to a probability prereq is met
func (g *game) canPlayCard(pindex int, card card) bool {
	// prerequisites given by card type
	if !card.cardType.prereq(g, pindex) {
		return false
	}

	p := g.getPlayer(pindex)
	// can player pay for the card?
	if !p.hasMana(card.manacost) {
		return false
	}
	// other prerequisites such as paying life
	// NOTE: prereq is target available?
	// --> this is handled by possibleTargets returning 0 actions
	for _, prereq := range card.prereqs {
		if !prereq(p) {
			return false
		}
	}
	return true
}

func (g *game) nextStep() {
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
	newG.stack = make([]action, len(g.stack))
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

func (g *game) play(a action) {
	p := g.getPlayer(a.controller)

	// remove card from players hand
	p.hand[a.card] -= 1
	if p.hand[a.card] == 0 {
		delete(p.hand, a.card)
	}

	c := cards[a.card]
	p.payMana(c.manacost)

	g.stack = append(g.stack, a)
}

func (g *game) resolve() {
	if len(g.stack) == 0 {
		panic("no stack to resolve")
	}
	a := g.stack[len(g.stack)-1]
	g.stack = g.stack[:len(g.stack)-1]

	c := cards[a.card]
	p := g.getPlayer(a.controller)
	c.resolve(p)

	// card specific effects
	for _, effect := range a.effects {
		effect.apply(g)
	}
}

func (g *game) isMainPhase() bool {
	return g.currentStep == precombatMainPhase || g.currentStep == postcombatMainPhase
}

func (g *game) getPlayerAction() action {
	return startMinimax(g)
}

func (g *game) getActions(index int) []action {
	p := g.getPlayer(index)
	actions := []action{action{card: pass, controller: index}}
	for k, _ := range p.hand {
		card := cards[k]
		if !g.canPlayCard(index, card) {
			continue
		}
		actions = append(actions, g.getCardActions(card, index)...)
	}
	return actions
}

func (g *game) getCardActions(card card, controller int) []action {
	if card.effects == nil {
		return []action{action{card: card.name, controller: controller}}
	}
	actions := []action{}
	effects := [][]effect{}
	for _, e := range card.effects {
		effects = append(effects, e.possibleTargets(controller, g))
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

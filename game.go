package main

import (
	"fmt"
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

var steps = []string{
	"untap step",
	"upkeep step",
	"draw step",
	"precombat main phase",
	"beginning of combat step",
	"declare attackers step",
	"declare blockers step",
	"first strike damage step",
	"combat damage step",
	"end of combat step",
	"postcombat main phase",
	"end step",
	"cleanup step",
}

func stepName(step step) string {
	return steps[step]
}

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
	numAttackers int
}

func newGame(startingPlayer int, players ...*player) *game {
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
		numPlayers:     len(players),
	}
	g.nextDecisionPoint()
	return g
}

// getPlayerAction -> resolveAction -> check gameEnds -> repeat
// rest is debugging print statements
func (g *game) loop() {
	for {
		a := g.getPlayerAction()
		if _, ok := a.(passAction); !ok {
			if at, ok := a.(attackAction); !ok || len(at.attackers) > 0 {
				g.debug()
			}
		}
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
					if target.ttype.isUntargeted() {
						fmt.Printf("%s resolves by %s \n", ac.card.getName(), g.getPlayer(ac.controller).name)
					} else {
						fmt.Printf("%s resolves by %s targeting %s \n", ac.card.getName(), g.getPlayer(ac.controller).name, g.getPlayer(int(target.index)).name)
					}
				}
			}
		case cardAction:
			fmt.Printf("-> %s plays %s", g.getPlayer(at.controller).name, at.card.getName())
			if len(at.targets) > 0 {
				fmt.Printf(" targeting %s", g.getPlayer(int(at.targets[0].index)).name)
			}
			fmt.Println()
		case attackAction:
			attackers := []string{}
			for _, c := range g.getPlayer(at.controller).battlefield.creatures {
				if c.attacking == -1 {
					continue
				}
				attackers = append(attackers, c.card.getName())
			}
			fmt.Printf("-> %s attacks with %s \n", g.getPlayer(at.controller).name, attackers)
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
		g.declareAttackers(a)
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
			g.untapStep()
			return
		case upkeepStep:
			break //skip
		case drawStep:
			g.drawStep()
			return
		case precombatMainPhase:
			return
		case beginningOfCombatStep:
			break //skip
		case declareAttackersStep:
			if g.declarations == 0 {
				// active player declares attackers
				return
			}
			// attackers have already been declared, continue this step
			// triggered abilities that trigger off attackers being declared trigger
			return
		case declareBlockersStep:
			// if no attackers, skip
			if g.numAttackers == 0 {
				break
			}
			break //skip
		case combatDamageFirstStrikeStep:
			// if no attackers, skip
			if g.numAttackers == 0 {
				break
			}
			break //skip
		case combatDamageStep:
			// if no attackers, skip
			if g.numAttackers == 0 {
				break
			}
			g.combatDamageStep()
			return
		case endOfCombatStep:
			g.endOfCombatStep()
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

func (g *game) untapStep() {
	activePlayer := g.getActivePlayer()
	for i, l := range activePlayer.battlefield.lands {
		l.tapped = false
		activePlayer.battlefield.lands[i] = l
	}
	for i, c := range activePlayer.battlefield.creatures {
		// TODO: edgecase: flash in creature start of turn
		c.summoningSickness = false
		c.tapped = false
		activePlayer.battlefield.creatures[i] = c
	}
}

func (g *game) drawStep() {
	g.getActivePlayer().draw()
}

func (g *game) combatDamageStep() {
	activePlayer := g.getActivePlayer()
	for i, c := range activePlayer.battlefield.creatures {
		if c.attacking == -1 {
			continue
		}
		defendingPlayer := g.getPlayer(c.attacking)
		defendingPlayer.lifeTotal -= c.card.(*creature).power
		c.attacking = -1
		activePlayer.battlefield.creatures[i] = c
	}
}

func (g *game) endOfCombatStep() {
	g.numAttackers = 0
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
	fmt.Printf("%s turn %d step %s: %s \n", activePlayer.name, g.turn, stepName(g.currentStep), activePlayer.String())
	fmt.Printf("           VS %s: %s \n", opp.name, opp.String())
}

func (g *game) play(a cardAction) {
	p := g.getPlayer(a.controller)

	// remove card from players hand
	p.hand[a.card] -= 1
	if p.hand[a.card] == 0 {
		delete(p.hand, a.card)
	}

	p.strategy.PayManaCost(p, a.card.getManaCost())

	g.stack = append(g.stack, a)
}

func (g *game) resolve() {
	if len(g.stack) == 0 {
		panic("no stack to resolve")
	}
	a := g.stack[len(g.stack)-1]
	g.stack = g.stack[:len(g.stack)-1]
	a.card.resolve(g, a)
}

func (g *game) declareAttackers(a attackAction) {
	p := g.getPlayer(a.getController())
	for _, att := range a.attackers {
		attacker := p.battlefield.creatures[att.index]
		attacker.attacking = att.target
		attacker.tapped = true
		p.battlefield.creatures[att.index] = attacker
	}
	g.numAttackers = len(a.attackers)
}

func (g *game) isMainPhase() bool {
	return g.currentStep == precombatMainPhase || g.currentStep == postcombatMainPhase
}

func (g *game) getPlayerAction() Action {
	p := g.players[g.priorityPlayer]
	if g.currentStep == declareAttackersStep && g.declarations == 0 {
		return p.strategy.Attacks(p, g)
	}
	return p.strategy.NextAction(p, g)
}

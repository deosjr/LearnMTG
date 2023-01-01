package main

type Action interface {
	getController() int
}

type action struct {
	controller int
}

func (a action) getController() int {
	return a.controller
}

type passAction struct {
	action
}

// TODO: similarly, activated abilities & triggers etc
type cardAction struct {
	action
	card    Card
    // index:  in card.effects
    // target: depends on effect type(?)
	targets []target
}

type attackAction struct {
	action
    // index:  in controller.battlefield
    // target: in game.players
	attackers []target
}

type blockAction struct {
	action
    // TODO: array of blocker/attacker pairs ?
}

// NOTE: the target struct is reused for multiple purposes right now
type target struct {
	index  int
	target int
}


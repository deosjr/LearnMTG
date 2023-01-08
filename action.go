package main

type Action interface {
	getController() int
}

type action struct {
	controller int
}

type target int

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
    // targets: used for casting spells with a target
    // i.e. instants and sorceries with spell abilities
    // index: in relevant zone(s), as per ability target type(s)
	targets []target
}

type attackAction struct {
	action
    // index:  in controller.battlefield
    // target: in game.players
	attackers []combattarget
}

type blockAction struct {
	action
    // index:  in controller.battlefield
    // target: in activeplayer.battlefield
    blockers []combattarget
}

type combattarget struct {
	index  int
	target int
}

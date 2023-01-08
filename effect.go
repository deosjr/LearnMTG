package main

// For permanent targets this might get hairy if their index changes (?)
// but players never change so this is simpler
// TODO: imagine targeting a creature by index, in response a creature with smaller
// index is killed, if battlefield array is reordered this goes very wrong
type Effect interface {
	apply(g *game, targets []target)
}

type damage struct {
    amount int
}

func (e damage) apply(g *game, targets []target) {
    // TODO: assumed player targets atm
    for _, p := range targets {
        g.getPlayer(int(p)).lifeTotal -= e.amount
    }
}

type lifegain struct {
    amount int
}

func (e lifegain) apply(g *game, targets []target) {
    for _, p := range targets {
        g.getPlayer(int(p)).lifeTotal += e.amount
    }
}

package main

// For permanent targets this might get hairy if their index changes (?)
// but players never change so this is simpler
// TODO: imagine targeting a creature by index, in response a creature with smaller
// index is killed, if battlefield array is reordered this goes very wrong
// TODO: use cardinstance.ID instead!
type Effect interface {
	apply(g *game, targets []effectTarget)
}

type draw struct {
	amount int
}

func (e draw) apply(g *game, targets []effectTarget) {
	for _, t := range targets {
		switch t.ttype {
		case you, targetPlayer:
			g.getPlayer(int(t.index)).drawN(e.amount)
		case eachPlayer:
			for _, p := range g.players {
				p.drawN(e.amount)
			}
		}
	}
}

type damage struct {
	amount int
}

func (e damage) apply(g *game, targets []effectTarget) {
	for _, t := range targets {
		switch t.ttype {
		case you, targetPlayer:
			g.getPlayer(int(t.index)).lifeTotal -= e.amount
		case eachPlayer:
			for _, p := range g.players {
				p.lifeTotal -= e.amount
			}
		}
	}
}

type lifegain struct {
	amount int
}

func (e lifegain) apply(g *game, targets []effectTarget) {
	for _, t := range targets {
		if !t.ttype.isPlayer() {
			panic("wrong target type")
		}
		g.getPlayer(int(t.index)).lifeTotal += e.amount
	}
}

type addMana struct {
	amount mana
}

func (e addMana) apply(g *game, targets []effectTarget) {
	if len(targets) != 1 && targets[0].ttype != you {
		panic("wrong target type")
	}
	p := g.getPlayer(int(targets[0].index))
	p.manaPool = p.manaPool.add(e.amount)
}

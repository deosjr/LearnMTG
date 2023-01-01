package main

// For permanent targets this might get hairy if their index changes (?)
// but players never change so this is simpler
// TODO: imagine targeting a creature by index, in response a creature with smaller
// index is killed, if battlefield array is reordered this goes very wrong
type Effect interface {
	possibleTargets(effectIndex, controllingPlayer int, game *game) []target
	getEffect() func(*player)
}

type effect struct {
	effect func(*player)
}

func (e effect) getEffect() func(*player) {
	return e.effect
}

func (c card) apply(g *game, target target) {
	e := c.getEffects()[target.index]
	f := e.getEffect()
	f(g.getPlayer(target.target))
}

type selfEffect struct {
	effect
}

type playerEffect struct {
	effect
}

func (e selfEffect) possibleTargets(effectIndex, controllingPlayer int, _ *game) []target {
	return []target{
		target{
			target: controllingPlayer,
			index:  effectIndex,
		},
	}
}

func (e playerEffect) possibleTargets(effectIndex, _ int, game *game) []target {
	effects := []target{}
	for i := 0; i < game.numPlayers; i++ {
		pe := target{
			target: i,
			index:  effectIndex,
		}
		effects = append(effects, pe)
	}
	return effects
}

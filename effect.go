package main

// For permanent targets this might get hairy if they change
// but players never change so this is simple
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

type target struct {
	index  int
	target int
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

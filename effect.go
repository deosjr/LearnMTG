package main

// For permanent targets this might get hairy if their index changes (?)
// but players never change so this is simpler
// TODO: imagine targeting a creature by index, in response a creature with smaller
// index is killed, if battlefield array is reordered this goes very wrong
type Effect interface {
	possibleTargets(controllingPlayer int, game *game) []target
	getEffect() func(*player)
}

type effect struct {
	effect func(*player)
}

func (e effect) getEffect() func(*player) {
	return e.effect
}

type selfEffect struct {
	effect
}

type playerEffect struct {
	effect
}

func (e selfEffect) possibleTargets(controllingPlayer int, _ *game) []target {
	return []target{
		target{
			target: controllingPlayer,
		},
	}
}

func (e playerEffect) possibleTargets(_ int, game *game) []target {
	effects := []target{}
	for i := 0; i < game.numPlayers; i++ {
		pe := target{
			target: i,
		}
		effects = append(effects, pe)
	}
	return effects
}

package main

// For permanent targets this might get hairy if they change
// but players never change so this is simple
type effect interface {
	possibleTargets(controllingPlayer int, game *game) []effect
	apply(*game)
}

type selfEffect struct {
	controller int
	target     int
	effect     func(*player)
}

func (e selfEffect) possibleTargets(controllingPlayer int, game *game) []effect {
	return []effect{
		selfEffect{
			controller: controllingPlayer,
			target:     controllingPlayer,
			effect:     e.effect,
		},
	}
}

func (e selfEffect) apply(game *game) {
	p := game.getPlayer(e.controller)
	e.effect(p)
}

type playerEffect struct {
	controller int
	target     int
	effect     func(*player)
}

func (e playerEffect) possibleTargets(controllingPlayer int, game *game) []effect {
	effects := []effect{}
	for i := 0; i < game.numPlayers; i++ {
		pe := playerEffect{
			controller: controllingPlayer,
			target:     i,
			effect:     e.effect,
		}
		effects = append(effects, pe)
	}
	return effects
}

func (e playerEffect) apply(game *game) {
	p := game.getPlayer(e.target)
	e.effect(p)
}

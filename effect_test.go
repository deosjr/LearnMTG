package main

import (
	"reflect"
	"testing"
)

func TestPossibleTargets(t *testing.T) {
	for i, tt := range []struct {
		effect     Effect
		controller int
		game       *game
		want       []target
	}{
		{
			effect:     selfEffect{},
			controller: SELF,
			game:       &game{numPlayers: 2},
			want:       []target{{target: SELF}},
		},
		{
			effect:     playerEffect{},
			controller: SELF,
			game:       &game{numPlayers: 2},
			want:       []target{{target: SELF}, {target: OPP}},
		},
		{
			effect:     playerEffect{},
			controller: OPP,
			game:       &game{numPlayers: 2},
			want:       []target{{target: SELF}, {target: OPP}},
		},
	} {
		got := tt.effect.possibleTargets(tt.controller, tt.game)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestApplySelfEffect(t *testing.T) {
	for i, tt := range []struct {
		effect func(*player)
		target target
		game   *game
		want   player
	}{
		{
			target: target{target: SELF},
			effect: func(p *player) { p.lifeTotal += 1 },
			game: &game{
				numPlayers: 2,
				players: []*player{
					SELF: &player{lifeTotal: 3},
					OPP:  &player{lifeTotal: 5},
				},
			},
			want: player{lifeTotal: 4},
		},
	} {
		spellAbility := selfEffect{effect: effect{effect: tt.effect}}
        action := cardAction{action: action{controller: SELF}, targets: []target{tt.target}}
        tt.game.apply(spellAbility, action)
		got := *tt.game.getPlayer(tt.target.target)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestApplyPlayerEffect(t *testing.T) {
	for i, tt := range []struct {
		effect func(*player)
		target target
		game   *game
		want   player
	}{
		{
			target: target{target: OPP},
			effect: func(p *player) { p.lifeTotal -= 3 },
			game: &game{
				numPlayers: 2,
				players: []*player{
					&player{lifeTotal: 4},
					&player{lifeTotal: 8},
				},
			},
			want: player{lifeTotal: 5},
		},
		{
			target: target{target: SELF},
			effect: func(p *player) { p.lifeTotal -= 3 },
			game: &game{
				numPlayers: 2,
				players: []*player{
					&player{lifeTotal: 4},
					&player{lifeTotal: 8},
				},
			},
			want: player{lifeTotal: 1},
		},
	} {
		spellAbility := playerEffect{effect: effect{effect: tt.effect}}
        action := cardAction{action: action{controller: SELF}, targets: []target{tt.target}}
        tt.game.apply(spellAbility, action)
		got := *tt.game.getPlayer(tt.target.target)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

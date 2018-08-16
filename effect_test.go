package main

import (
	"reflect"
	"testing"
)

func TestPossibleTargets(t *testing.T) {
	for i, tt := range []struct {
		effect     effect
		controller int
		game       *game
		want       []effect
	}{
		{
			effect:     selfEffect{},
			controller: SELF,
			game:       &game{numPlayers: 2},
			want:       []effect{selfEffect{target: SELF}},
		},
		{
			effect:     playerEffect{},
			controller: SELF,
			game:       &game{numPlayers: 2},
			want:       []effect{playerEffect{target: SELF}, playerEffect{target: OPP}},
		},
		{
			effect:     playerEffect{},
			controller: OPP,
			game:       &game{numPlayers: 2},
			want:       []effect{playerEffect{target: SELF}, playerEffect{target: OPP}},
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
		effect selfEffect
		game   *game
		want   player
	}{
		{
			effect: selfEffect{
				target: SELF,
				effect: func(p *player) { p.lifeTotal += 1 },
			},
			game: &game{
				numPlayers: 2,
				players: []*player{
					&player{lifeTotal: 3},
					&player{lifeTotal: 5},
				},
			},
			want: player{lifeTotal: 4},
		},
	} {
		tt.effect.apply(tt.game)
		got := *tt.game.getPlayer(tt.effect.target)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestApplyPlayerEffect(t *testing.T) {
	for i, tt := range []struct {
		effect playerEffect
		game   *game
		want   player
	}{
		{
			effect: playerEffect{
				target: OPP,
				effect: func(p *player) { p.lifeTotal -= 3 },
			},
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
			effect: playerEffect{
				target: SELF,
				effect: func(p *player) { p.lifeTotal -= 3 },
			},
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
		tt.effect.apply(tt.game)
		got := *tt.game.getPlayer(tt.effect.target)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

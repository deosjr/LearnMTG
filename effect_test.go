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
			controller: 0,
			game:       &game{numPlayers: 2},
			want:       []effect{selfEffect{target: 0}},
		},
		{
			effect:     playerEffect{},
			controller: 0,
			game:       &game{numPlayers: 2},
			want:       []effect{playerEffect{target: 0}, playerEffect{target: 1}},
		},
		{
			effect:     playerEffect{},
			controller: 1,
			game:       &game{numPlayers: 2},
			want:       []effect{playerEffect{controller: 1, target: 0}, playerEffect{controller: 1, target: 1}},
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
				controller: 0,
				target:     0,
				effect:     func(p *player) { p.lifeTotal += 1 },
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
				controller: 0,
				target:     1,
				effect:     func(p *player) { p.lifeTotal -= 3 },
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
				controller: 1,
				target:     0,
				effect:     func(p *player) { p.lifeTotal -= 3 },
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

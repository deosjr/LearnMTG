package main

import (
	"reflect"
	"testing"
)

func TestApplySelfEffect(t *testing.T) {
	for i, tt := range []struct {
		effect Effect
		target target
		game   *game
		want   player
	}{
		{
			target: target(SELF),
			effect: lifegain{1},
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
        tt.effect.apply(tt.game, []target{tt.target})
		got := *tt.game.getPlayer(int(tt.target))
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestApplyPlayerEffect(t *testing.T) {
	for i, tt := range []struct {
		effect Effect
		target target
		game   *game
		want   player
	}{
		{
			target: target(OPP),
            effect: damage{3},
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
			target: target(SELF),
            effect: damage{3},
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
        tt.effect.apply(tt.game, []target{tt.target})
		got := *tt.game.getPlayer(int(tt.target))
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

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
        tt.effect.apply(tt.game, []effectTarget{{index:tt.target, ttype:you}})
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
        tt.effect.apply(tt.game, []effectTarget{{index:tt.target, ttype:targetPlayer}})
		got := *tt.game.getPlayer(int(tt.target))
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestApplyEachPlayerEffect(t *testing.T) {
	for i, tt := range []struct {
		effect Effect
		game   *game
		wantP1 player
		wantP2 player
	}{
		{
            effect: damage{4},
			game: &game{
				numPlayers: 2,
				players: []*player{
					&player{lifeTotal: 5},
					&player{lifeTotal: 8},
				},
			},
			wantP1: player{lifeTotal: 1},
			wantP2: player{lifeTotal: 4},
		},
	} {
        tt.effect.apply(tt.game, []effectTarget{{ttype:eachPlayer}})
		gotP1 := *tt.game.getPlayer(int(SELF))
		if !reflect.DeepEqual(gotP1, tt.wantP1) {
			t.Errorf("%d) P1 got %v want %v", i, gotP1, tt.wantP1)
		}
		gotP2 := *tt.game.getPlayer(int(OPP))
		if !reflect.DeepEqual(gotP2, tt.wantP2) {
			t.Errorf("%d) P2 got %v want %v", i, gotP2, tt.wantP2)
		}
	}
}

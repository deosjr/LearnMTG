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
			want:       []effect{playerEffect{target: 0}, playerEffect{target: 1}},
		},
	} {
		got := tt.effect.possibleTargets(tt.controller, tt.game)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

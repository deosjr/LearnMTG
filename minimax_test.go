package main

import (
	"reflect"
	"testing"
)

func TestGetChild(t *testing.T) {
	for i, tt := range []struct {
		node   node
		action Action
		want   node
	}{
		{
			node: node{
				game: &game{
					players: []*player{
						SELF: &player{},
						OPP:  &player{},
					},
					priorityPlayer: SELF,
					activePlayer:   OPP,
					currentStep:    precombatMainPhase,
					numPasses:      0,
				},
				pointOfView: SELF,
			},
			action: passAction{action: action{controller: SELF}},
			want: node{
				game: &game{
					players: []*player{
						SELF: &player{},
						OPP:  &player{},
					},
					priorityPlayer: OPP,
					activePlayer:   OPP,
					currentStep:    precombatMainPhase,
					numPasses:      1,
				},
				pointOfView: SELF,
			},
		},
	} {
		tt.node.game.numPlayers = 2
		tt.want.game.numPlayers = 2
		got := tt.node.getChild(tt.action)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %#v want %#v", i, got.game, tt.want.game)
		}
	}
}

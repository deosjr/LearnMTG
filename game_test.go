package main

import (
	"reflect"
	"sort"
	"testing"
)

const (
	SELF = 0
	OPP  = 1
)

func TestGetPlayerAction(t *testing.T) {
	for i, tt := range []struct {
		g    *game
		want action
	}{
		{
			g: &game{
				players: []*player{
					SELF: &player{},
					OPP:  &player{},
				},
				priorityPlayer: SELF,
				activePlayer:   OPP,
				currentStep:    precombatMainPhase,
			},
			want: action{card: pass, controller: SELF},
		},
		{
			g: &game{
				players: []*player{
					SELF: &player{
						hand: map[string]int{
							mountain.name:  1,
							lavaSpike.name: 1,
						},
						manaAvailable: 1,
						lifeTotal:     20,
					},
					OPP: &player{
						lifeTotal: 3,
					},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
			want: action{card: lavaSpike.name, controller: SELF},
		},
		{
			g: &game{
				players: []*player{
					SELF: &player{
						lifeTotal: 3,
					},
					OPP: &player{
						hand: map[string]int{
							mountain.name:  1,
							lavaSpike.name: 1,
						},
						manaAvailable: 1,
						lifeTotal:     20,
					},
				},
				priorityPlayer: OPP,
				activePlayer:   OPP,
				currentStep:    precombatMainPhase,
			},
			want: action{card: lavaSpike.name, controller: OPP},
		},
	} {
		oldMax := maxDepth
		maxDepth = 5
		defer func() {
			maxDepth = oldMax
		}()
		tt.g.numPlayers = 2
		got := tt.g.getPlayerAction()
		got.effects = nil
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestGetActions(t *testing.T) {
	for i, tt := range []struct {
		game        *game
		pointOfView int
		want        []action
	}{
		// SELF MOVES
		{
			game: &game{
				players: []*player{
					SELF: &player{},
					OPP:  &player{},
				},
				activePlayer:   OPP,
				priorityPlayer: SELF,
				currentStep:    precombatMainPhase,
			},
			pointOfView: SELF,
			want: []action{
				action{card: pass},
			},
		},
		{
			game: &game{
				players: []*player{
					SELF: &player{
						hand: map[string]int{
							mountain.name:  2,
							lavaSpike.name: 3,
						},
						manaAvailable: 0,
						lifeTotal:     20,
					},
					OPP: &player{},
				},
				activePlayer:   SELF,
				priorityPlayer: SELF,
				currentStep:    precombatMainPhase,
			},
			pointOfView: SELF,
			want: []action{
				{card: mountain.name, controller: SELF},
				action{card: pass},
			},
		},
		{
			game: &game{
				players: []*player{
					SELF: &player{
						hand: map[string]int{
							mountain.name:  2,
							lavaSpike.name: 3,
						},
						manaAvailable: 1,
						lifeTotal:     20,
					},
					OPP: &player{},
				},

				activePlayer:   SELF,
				priorityPlayer: SELF,
				currentStep:    precombatMainPhase,
			},
			pointOfView: SELF,
			want: []action{
				// targetting self, targetting opp
				{card: lavaSpike.name, controller: SELF},
				{card: lavaSpike.name, controller: SELF},
				{card: mountain.name, controller: SELF},
				action{card: pass},
			},
		},

		// OPPONENT MOVES
		{
			game: &game{
				players: []*player{
					SELF: &player{},
					OPP:  &player{},
				},
				activePlayer:   SELF,
				priorityPlayer: OPP,
				currentStep:    precombatMainPhase,
			},
			pointOfView: SELF,
			want:        []action{action{card: pass, controller: OPP}},
		},
		{
			game: &game{
				players: []*player{
					SELF: &player{},
					OPP: &player{
						hand: map[string]int{
							mountain.name:  2,
							lavaSpike.name: 3,
						},
						manaAvailable: 0,
						lifeTotal:     20,
						deckList:      deckList,
					},
				},

				activePlayer:   OPP,
				priorityPlayer: OPP,
				currentStep:    precombatMainPhase,
			},
			pointOfView: SELF,
			want: []action{
				{card: mountain.name, controller: OPP},
				action{card: pass, controller: OPP},
			},
		},
		{
			game: &game{
				players: []*player{
					SELF: &player{},
					OPP: &player{
						hand: map[string]int{
							mountain.name:  2,
							lavaSpike.name: 3,
						},
						manaAvailable: 1,
						lifeTotal:     20,
						deckList:      deckList,
					},
				},

				activePlayer:   OPP,
				priorityPlayer: OPP,
				currentStep:    precombatMainPhase,
			},
			pointOfView: SELF,
			want: []action{
				// targetting self, targetting opp
				{card: lavaSpike.name, controller: OPP},
				{card: lavaSpike.name, controller: OPP},
				{card: mountain.name, controller: OPP},
				action{card: pass, controller: OPP},
			},
		},
		{
			game: &game{
				players: []*player{
					SELF: &player{},
					OPP: &player{
						hand:          map[string]int{},
						manaAvailable: 1,
						lifeTotal:     20,
						deckList:      deckList,
					},
				},

				activePlayer:   OPP,
				priorityPlayer: OPP,
				currentStep:    precombatMainPhase,
			},
			pointOfView: SELF,
			want:        []action{action{card: pass, controller: OPP}},
		},
	} {
		tt.game.numPlayers = 2
		got := tt.game.getActions(tt.game.priorityPlayer)
		for i, a := range got {
			a.effects = nil
			got[i] = a
		}
		sort.Slice(got, func(i, j int) bool {
			return got[i].card < got[j].card
		})
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

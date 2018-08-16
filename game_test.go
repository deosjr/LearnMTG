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
		name string
		game *game
		want action
	}{
		{
			name: "no options means pass",
			game: &game{
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
			name: "lava spike for the win",
			game: &game{
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
						library:   []string{mountain.name},
					},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
			want: action{card: lavaSpike.name, controller: SELF},
		},
		{
			name: "opp lava spike for the win",
			game: &game{
				players: []*player{
					SELF: &player{
						lifeTotal: 3,
						library:   []string{mountain.name},
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
		tt.game.numPlayers = 2
		got := tt.game.getPlayerAction()
		got.effects = nil
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d: %s) got %v want %v", i, tt.name, got, tt.want)
		}
	}
}

func TestGetActions(t *testing.T) {
	for i, tt := range []struct {
		name        string
		game        *game
		pointOfView int
		want        []action
	}{
		// SELF MOVES
		{
			name: "no options means pass",
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
			name: "no mana no lava spike",
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
			name: "all the options",
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
					OPP: &player{
						lifeTotal: 3,
					},
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
		{
			name: "card on the stack -> no sorceries",
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
					OPP: &player{
						lifeTotal: 3,
					},
				},

				activePlayer:   SELF,
				priorityPlayer: SELF,
				currentStep:    precombatMainPhase,
				stack:          []action{{card: lavaSpike.name, controller: SELF}},
			},
			pointOfView: SELF,
			want:        []action{action{card: pass}},
		},

		// OPPONENT MOVES
		{
			name: "opp pass",
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
			name: "opp no mana",
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
			name: "opp all the options",
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
			name: "opp no cards in hand",
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
			t.Errorf("%d: %s) got %v want %v", i, tt.name, got, tt.want)
		}
	}
}

func TestResolveAction(t *testing.T) {
	for i, tt := range []struct {
		name   string
		game   *game
		action action
		want   *game
	}{
		{
			name: "pass action",
			game: &game{
				players: []*player{
					SELF: &player{},
					OPP:  &player{},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
			action: action{card: pass, controller: SELF},
			want: &game{
				players: []*player{
					SELF: &player{},
					OPP:  &player{},
				},
				priorityPlayer: OPP,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
				numPasses:      1,
			},
		},
		{
			name: "play lava spike",
			game: &game{
				players: []*player{
					SELF: &player{
						hand: map[string]int{
							lavaSpike.name: 1,
						},
						manaAvailable: 1,
					},
					OPP: &player{},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
			action: action{
				card:       lavaSpike.name,
				controller: SELF,
				effects: []effect{
					playerEffect{
						target: OPP,
					},
				},
			},
			want: &game{
				players: []*player{
					SELF: &player{
						hand:          map[string]int{},
						manaAvailable: 0,
					},
					OPP: &player{},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
				stack: []action{
					action{
						card:       lavaSpike.name,
						controller: SELF,
						effects: []effect{
							playerEffect{
								target: OPP,
							},
						},
					},
				},
			},
		},
	} {
		tt.game.numPlayers = 2
		tt.want.numPlayers = 2
		got := tt.game.copy()
		got.resolveAction(tt.action)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d: %s) got %#v want %#v", i, tt.name, got, tt.want)
		}
	}
}

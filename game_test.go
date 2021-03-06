package main

import (
	"fmt"
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
		want Action
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
			want: passAction{action{controller: SELF}},
		},
		{
			name: "lava spike for the win",
			game: &game{
				players: []*player{
					SELF: &player{
						hand: map[Card]int{
							mountain:  1,
							lavaSpike: 1,
						},
						manaAvailable: 1,
						lifeTotal:     20,
					},
					OPP: &player{
						lifeTotal: 3,
						library:   []Card{mountain},
					},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
			want: cardAction{
				card:   lavaSpike,
				action: action{controller: SELF},
				targets: []target{
					{target: OPP},
				},
			},
		},
		{
			name: "opp lava spike for the win",
			game: &game{
				players: []*player{
					SELF: &player{
						lifeTotal: 3,
						library:   []Card{mountain},
					},
					OPP: &player{
						hand: map[Card]int{
							mountain:  1,
							lavaSpike: 1,
						},
						manaAvailable: 1,
						lifeTotal:     20,
					},
				},
				priorityPlayer: OPP,
				activePlayer:   OPP,
				currentStep:    precombatMainPhase,
			},
			want: cardAction{
				card:   lavaSpike,
				action: action{controller: OPP},
				targets: []target{
					{target: SELF},
				},
			},
		},
	} {
		oldMax := maxDepth
		maxDepth = 5
		defer func() {
			maxDepth = oldMax
		}()
		tt.game.numPlayers = 2
		got := tt.game.getPlayerAction()
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
		want        []Action
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
			want: []Action{
				passAction{action{controller: SELF}},
			},
		},
		{
			name: "no mana no lava spike",
			game: &game{
				players: []*player{
					SELF: &player{
						hand: map[Card]int{
							mountain:  2,
							lavaSpike: 3,
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
			want: []Action{
				cardAction{card: mountain, action: action{controller: SELF}},
				passAction{action{controller: SELF}},
			},
		},
		{
			name: "all the options",
			game: &game{
				players: []*player{
					SELF: &player{
						hand: map[Card]int{
							mountain:  2,
							lavaSpike: 3,
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
			want: []Action{
				cardAction{
					card:   lavaSpike,
					action: action{controller: SELF},
					targets: []target{
						{target: SELF},
					},
				},
				cardAction{
					card:   lavaSpike,
					action: action{controller: SELF},
					targets: []target{
						{target: OPP},
					},
				},
				cardAction{card: mountain, action: action{controller: SELF}},
				passAction{action{controller: SELF}},
			},
		},
		{
			name: "card on the stack -> no sorceries",
			game: &game{
				players: []*player{
					SELF: &player{
						hand: map[Card]int{
							mountain:  2,
							lavaSpike: 3,
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
				stack:          []cardAction{cardAction{card: lavaSpike, action: action{controller: SELF}}},
			},
			pointOfView: SELF,
			want:        []Action{passAction{action{controller: SELF}}},
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
			want:        []Action{passAction{action{controller: OPP}}},
		},
		{
			name: "opp no mana",
			game: &game{
				players: []*player{
					SELF: &player{},
					OPP: &player{
						hand: map[Card]int{
							mountain:  2,
							lavaSpike: 3,
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
			want: []Action{
				cardAction{card: mountain, action: action{controller: OPP}},
				passAction{action{controller: OPP}},
			},
		},
		{
			name: "opp all the options",
			game: &game{
				players: []*player{
					SELF: &player{},
					OPP: &player{
						hand: map[Card]int{
							mountain:  2,
							lavaSpike: 3,
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
			want: []Action{
				cardAction{
					card:   lavaSpike,
					action: action{controller: OPP},
					targets: []target{
						{target: OPP},
					},
				},
				cardAction{
					card:   lavaSpike,
					action: action{controller: OPP},
					targets: []target{
						{target: SELF},
					},
				},
				cardAction{card: mountain, action: action{controller: OPP}},
				passAction{action{controller: OPP}},
			},
		},
		{
			name: "opp no cards in hand",
			game: &game{
				players: []*player{
					SELF: &player{},
					OPP: &player{
						hand:          map[Card]int{},
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
			want:        []Action{passAction{action{controller: OPP}}},
		},
	} {
		tt.game.numPlayers = 2
		got := tt.game.getActions(tt.game.priorityPlayer)
		sort.Slice(got, func(i, j int) bool {
			return fmt.Sprintf("%v", got[i]) < fmt.Sprintf("%v", got[j])
		})
		sort.Slice(tt.want, func(i, j int) bool {
			return fmt.Sprintf("%v", tt.want[i]) < fmt.Sprintf("%v", tt.want[j])
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
		action Action
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
			action: passAction{action{controller: SELF}},
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
						hand: map[Card]int{
							lavaSpike: 1,
						},
						manaAvailable: 1,
					},
					OPP: &player{},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
			action: cardAction{
				card:   lavaSpike,
				action: action{controller: SELF},
				targets: []target{
					target{
						target: OPP,
					},
				},
			},
			want: &game{
				players: []*player{
					SELF: &player{
						hand:          map[Card]int{},
						manaAvailable: 0,
					},
					OPP: &player{},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
				stack: []cardAction{
					cardAction{
						card:   lavaSpike,
						action: action{controller: SELF},
						targets: []target{
							target{
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

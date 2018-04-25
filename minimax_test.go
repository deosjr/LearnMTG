package main

import (
	"reflect"
	"sort"
	"testing"
)

var fullMountainLibrary = []string{
	mountain.name, mountain.name, mountain.name, mountain.name, mountain.name,
	mountain.name, mountain.name, mountain.name, mountain.name, mountain.name,
	mountain.name, mountain.name, mountain.name, mountain.name, mountain.name,
	mountain.name, mountain.name, mountain.name, mountain.name, mountain.name,
	mountain.name, mountain.name, mountain.name, mountain.name, mountain.name,
	mountain.name, mountain.name, mountain.name, mountain.name, mountain.name,
}

func TestGetPlayerAction(t *testing.T) {
	for i, tt := range []struct {
		g    *game
		want action
	}{
		{
			g: &game{
				players: []*player{
					0: &player{
						name:    "self",
						library: fullMountainLibrary,
					},
					1: &player{
						name:    "opp",
						library: fullMountainLibrary,
					},
				},
				priorityPlayer: 0,
				activePlayer:   1,
				currentStep:    precombatMainPhase,
			},
			want: action{card: pass, playerSelf: 0},
		},
		{
			g: &game{
				players: []*player{
					0: &player{
						name: "self",
						hand: map[string]int{
							mountain.name:  1,
							lavaSpike.name: 1,
						},
						library:       fullMountainLibrary,
						manaAvailable: 1,
						lifeTotal:     20,
					},
					1: &player{
						name:      "opp",
						lifeTotal: 3,
						library:   fullMountainLibrary,
					},
				},
				priorityPlayer: 0,
				activePlayer:   0,
				currentStep:    precombatMainPhase,
			},
			want: action{card: lavaSpike.name, playerSelf: 0, playerTarget: 1},
		},
	} {
		oldMax := maxDepth
		maxDepth = 10
		defer func() {
			maxDepth = oldMax
		}()
		tt.g.numPlayers = 2
		got := tt.g.getPlayerAction()
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestGetActionsSelf(t *testing.T) {
	for i, tt := range []struct {
		node node
		want []action
	}{
		{
			node: node{
				game: &game{
					players: []*player{
						0: &player{name: "self"},
						1: &player{name: "opp"},
					},
					activePlayer: 1,
					currentStep:  precombatMainPhase,
				},
				pointOfView: 0,
			},
			want: []action{
				{card: pass, playerSelf: 0},
			},
		},
		{
			node: node{
				game: &game{
					players: []*player{
						0: &player{
							name: "self",
							hand: map[string]int{
								mountain.name:  2,
								lavaSpike.name: 3,
							},
							manaAvailable: 0,
							lifeTotal:     20,
						},
						1: &player{name: "opp"},
					},

					activePlayer: 0,
					currentStep:  precombatMainPhase,
				},
				pointOfView: 0,
			},
			want: []action{
				{card: mountain.name, playerSelf: 0},
				{card: pass, playerSelf: 0},
			},
		},
		{
			node: node{
				game: &game{
					players: []*player{
						0: &player{
							name: "self",
							hand: map[string]int{
								mountain.name:  2,
								lavaSpike.name: 3,
							},
							manaAvailable: 1,
							lifeTotal:     20,
						},
						1: &player{name: "opp"},
					},

					activePlayer: 0,
					currentStep:  precombatMainPhase,
				},
				pointOfView: 0,
			},
			want: []action{
				{card: lavaSpike.name, playerSelf: 0, playerTarget: 1},
				{card: mountain.name, playerSelf: 0},
				{card: pass, playerSelf: 0},
			},
		},
	} {
		got := tt.node.getActionsSelf()
		sort.Slice(got, func(i, j int) bool {
			return got[i].card < got[j].card
		})
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestGetActionsOpponent(t *testing.T) {
	for i, tt := range []struct {
		node node
		want []action
	}{
		{
			node: node{
				game: &game{
					players: []*player{
						0: &player{name: "self"},
						1: &player{name: "opp"},
					},
					activePlayer: 0,
					currentStep:  precombatMainPhase,
				},
				pointOfView: 0,
			},
			want: []action{
				{card: pass, playerSelf: 1},
			},
		},
		{
			node: node{
				game: &game{
					players: []*player{
						0: &player{
							name: "self",
						},
						1: &player{
							name: "opp",
							hand: map[string]int{
								mountain.name:  2,
								lavaSpike.name: 3,
							},
							manaAvailable: 0,
							lifeTotal:     20,
							deckList:      deckList,
						},
					},

					activePlayer: 1,
					currentStep:  precombatMainPhase,
				},
				pointOfView: 0,
			},
			want: []action{
				{card: mountain.name, playerSelf: 1},
				{card: pass, playerSelf: 1},
			},
		},
		{
			node: node{
				game: &game{
					players: []*player{
						0: &player{
							name: "self",
						},
						1: &player{
							name: "opp",
							hand: map[string]int{
								mountain.name:  2,
								lavaSpike.name: 3,
							},
							manaAvailable: 1,
							lifeTotal:     20,
							deckList:      deckList,
						},
					},

					activePlayer: 1,
					currentStep:  precombatMainPhase,
				},
				pointOfView: 0,
			},
			want: []action{
				{card: lavaSpike.name, playerSelf: 1, playerTarget: 0},
				{card: mountain.name, playerSelf: 1},
				{card: pass, playerSelf: 1},
			},
		},
		{
			node: node{
				game: &game{
					players: []*player{
						0: &player{
							name: "self",
						},
						1: &player{
							name:          "opp",
							hand:          map[string]int{},
							manaAvailable: 1,
							lifeTotal:     20,
							deckList:      deckList,
						},
					},

					activePlayer: 1,
					currentStep:  precombatMainPhase,
				},
				pointOfView: 0,
			},
			want: []action{
				{card: pass, playerSelf: 1},
			},
		},
	} {
		got := tt.node.getActionsOpponent()
		sort.Slice(got, func(i, j int) bool {
			return got[i].card < got[j].card
		})
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

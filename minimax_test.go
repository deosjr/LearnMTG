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

func TestPlayerAct(t *testing.T) {
	for i, tt := range []struct {
		p      *player
		active bool
		want   action
	}{
		{
			p: &player{
				name:     "self",
				opponent: &player{name: "opp", library: fullMountainLibrary},
				library:  fullMountainLibrary,
			},
			active: false,
			want:   action{card: pass, player: "self"},
		},
		{
			p: &player{
				name: "self",
				opponent: &player{
					name:      "opp",
					lifeTotal: 3,
					library:   fullMountainLibrary,
				},
				hand: map[string]int{
					mountain.name:  1,
					lavaSpike.name: 1,
				},
				library:       fullMountainLibrary,
				manaAvailable: 1,
				lifeTotal:     20,
			},
			active: true,
			want:   action{card: lavaSpike.name, player: "self"},
		},
	} {
		oldMax := maxDepth
		maxDepth = 10
		defer func() {
			maxDepth = oldMax
		}()
		got := tt.p.act(tt.active, precombatMainPhase)
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
				pointOfView: &player{name: "self"},
				isActive:    false,
				currentStep: precombatMainPhase,
			},
			want: []action{
				{card: pass, player: "self"},
			},
		},
		{
			node: node{
				pointOfView: &player{
					name: "self",
					hand: map[string]int{
						mountain.name:  2,
						lavaSpike.name: 3,
					},
					manaAvailable: 0,
					lifeTotal:     20,
				},
				isActive:    true,
				currentStep: precombatMainPhase,
			},
			want: []action{
				{card: mountain.name, player: "self"},
				{card: pass, player: "self"},
			},
		},
		{
			node: node{
				pointOfView: &player{
					name: "self",
					hand: map[string]int{
						mountain.name:  2,
						lavaSpike.name: 3,
					},
					manaAvailable: 1,
					lifeTotal:     20,
				},
				isActive:    true,
				currentStep: precombatMainPhase,
			},
			want: []action{
				{card: lavaSpike.name, player: "self"},
				{card: mountain.name, player: "self"},
				{card: pass, player: "self"},
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
				pointOfView: &player{name: "self", opponent: &player{name: "opp"}},
				isActive:    false,
				currentStep: precombatMainPhase,
			},
			want: []action{
				{card: pass, player: "opp"},
			},
		},
		{
			node: node{
				pointOfView: &player{
					name: "self",
					opponent: &player{
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
				isActive:    true,
				currentStep: precombatMainPhase,
			},
			want: []action{
				{card: mountain.name, player: "opp"},
				{card: pass, player: "opp"},
			},
		},
		{
			node: node{
				pointOfView: &player{
					name: "self",
					opponent: &player{
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
				isActive:    true,
				currentStep: precombatMainPhase,
			},
			want: []action{
				{card: lavaSpike.name, player: "opp"},
				{card: mountain.name, player: "opp"},
				{card: pass, player: "opp"},
			},
		},
		// opponent on all lands but we dont know?
		// guess based on their decklist
		{
			node: node{
				pointOfView: &player{
					name: "self",
					opponent: &player{
						name: "opp",
						hand: map[string]int{
							mountain.name: 2,
						},
						manaAvailable: 1,
						lifeTotal:     20,
						deckList:      deckList,
					},
				},
				isActive:    true,
				currentStep: precombatMainPhase,
			},
			want: []action{
				{card: lavaSpike.name, player: "opp"},
				{card: mountain.name, player: "opp"},
				{card: pass, player: "opp"},
			},
		},
		{
			node: node{
				pointOfView: &player{
					name: "self",
					opponent: &player{
						name:          "opp",
						hand:          map[string]int{},
						manaAvailable: 1,
						lifeTotal:     20,
						deckList:      deckList,
					},
				},
				isActive:    true,
				currentStep: precombatMainPhase,
			},
			want: []action{
				{card: pass, player: "opp"},
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

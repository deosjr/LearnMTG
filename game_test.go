package main

import (
	"reflect"
	"testing"
)

const (
	SELF = 0
	OPP  = 1
)

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

func TestResolveVSGoldfish(t *testing.T) {
	for i, tt := range []struct {
		name    string
		game    *game
		actions []Action
		want    *game
	}{
		{
			name: "pass action",
			game: &game{
				players: []*player{
					SELF: &player{},
					OPP:  &player{strategy:goldfish{}},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
			actions: []Action{passAction{action{controller: SELF}}},
			want: &game{
				players: []*player{
					SELF: &player{},
					OPP:  &player{strategy:goldfish{}},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
                // TODO: currently skipping over beginningOfCombat
				currentStep:    declareAttackersStep,
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
					OPP: &player{strategy:goldfish{}},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
			actions: []Action{cardAction{
				    card:   lavaSpike,
				    action: action{controller: SELF},
				    targets: []target{
					    target{
						    target: OPP,
					    },
				    },
                },
                passAction{action{controller: SELF}},
			},
			want: &game{
				players: []*player{
					SELF: &player{
						hand:          map[Card]int{},
						manaAvailable: 0,
                        graveyard:     orderedCards{lavaSpike},
					},
					OPP: &player{
                        strategy:goldfish{},
                        lifeTotal: -3,
                    },
				},
                stack:          []cardAction{},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
		},
	} {
		tt.game.numPlayers = 2
		tt.want.numPlayers = 2
		got := tt.game.copy()
        for _, a := range tt.actions {
		    got.resolveAction(a)
        }
        opp := got.players[OPP]
        got.resolveAction(opp.strategy.NextAction(opp, got))
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d: %s) got %#v want %#v", i, tt.name, got, tt.want)
		}
	}
}

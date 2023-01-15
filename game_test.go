package main

import (
	"reflect"
	"testing"
)

const (
	SELF = 0
	OPP  = 1
)

func testManaAvailable(n int) battlefield {
    lands := make([]cardInstance, n)
    for i := range lands {
        lands[i] = cardInstance{card:mountain}
    }
    return battlefield{lands: lands}
}

func testManaTapUntap(t, u int) battlefield {
    lands := make([]cardInstance, t+u)
    for i:=0; i<u; i++ {
        lands[i] = cardInstance{card:mountain}
    }
    for i:=u; i<t+u; i++ {
        lands[i] = cardInstance{card:mountain, tapped:true}
    }
    return battlefield{lands: lands}
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
					SELF: &player{strategy:goldfish{}},
					OPP:  &player{strategy:goldfish{}},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
			action: passAction{action{controller: SELF}},
			want: &game{
				players: []*player{
					SELF: &player{strategy:goldfish{}},
					OPP:  &player{strategy:goldfish{}},
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
                        battlefield: testManaAvailable(1),
                        strategy: goldfish{},
					},
					OPP: &player{strategy:goldfish{}},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
			action: cardAction{
				card:   lavaSpike,
				action: action{controller: SELF},
				targets: []effectTarget{ {index:target(OPP), ttype: targetPlayer} },
			},
			want: &game{
				players: []*player{
					SELF: &player{
						hand:          map[Card]int{},
                        battlefield:   testManaTapUntap(1, 0),
                        strategy:      goldfish{},
					},
					OPP: &player{strategy:goldfish{}},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
				stack: []cardAction{
					cardAction{
						card:   lavaSpike,
						action: action{controller: SELF},
				        targets: []effectTarget{ {index:target(OPP), ttype: targetPlayer} },
					},
				},
			},
		},
		{
			name: "play mountain",
			game: &game{
				players: []*player{
					SELF: &player{
						hand: map[Card]int{
							mountain: 1,
						},
                        battlefield: testManaAvailable(1),
                        strategy:    goldfish{},
					},
					OPP: &player{strategy:goldfish{}},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
			action: cardAction{
				card:   mountain,
				action: action{controller: SELF},
			},
			want: &game{
				players: []*player{
					SELF: &player{
						hand:          map[Card]int{},
                        battlefield:   testManaAvailable(2),
                        landPlayed:    true,
                        strategy:      goldfish{},
					},
					OPP: &player{strategy:goldfish{}},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
                stack:          []cardAction{},
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
					SELF: &player{strategy:goldfish{}},
					OPP:  &player{strategy:goldfish{}},
				},
				priorityPlayer: SELF,
				activePlayer:   SELF,
				currentStep:    precombatMainPhase,
			},
			actions: []Action{passAction{action{controller: SELF}}},
			want: &game{
				players: []*player{
					SELF: &player{strategy:goldfish{}},
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
                        battlefield:   testManaAvailable(1),
                        strategy:      goldfish{},
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
				    targets: []effectTarget{ {index:target(OPP), ttype: targetPlayer} },
                },
                passAction{action{controller: SELF}},
			},
			want: &game{
				players: []*player{
					SELF: &player{
						hand:          map[Card]int{},
                        battlefield:   testManaTapUntap(1, 0),
                        graveyard:     orderedCards{lavaSpike},
                        strategy:      goldfish{},
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

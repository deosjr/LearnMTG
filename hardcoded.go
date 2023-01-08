package main

var (
	mountain = &land{
		card: card{
			name: "Mountain",
		},
	}

	lavaSpike = &sorcery{
		card: card{
			name:     "Lava Spike",
			manaCost: manaCost{r: 1},
		},
        spellAbility: playerEffect{
		    effect: effect{func(p *player) {
				p.lifeTotal -= 3
			}},
        },
	}

	falkenrathReaver = &creature{
		card: card{
			name:     "Falkenrath Reaver",
			manaCost: manaCost{c: 1, r: 1},
		},
		power:     2,
		toughness: 2,
	}

	cards = map[string]Card{
		mountain.name:         mountain,
		lavaSpike.name:        lavaSpike,
		falkenrathReaver.name: falkenrathReaver,
	}

	deckList = unorderedCards{
		mountain:         10,
		lavaSpike:        10,
		falkenrathReaver: 10,
	}
)


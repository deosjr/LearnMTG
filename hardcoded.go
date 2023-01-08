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
        spellAbility: SpellAbility{
            ability{
                targets: []targetType{targetPlayer},
                effect:  damage{3},
            },
        },
	}

    flameRift = &sorcery{
        card: card{
            name:     "Flame Rift",
            manaCost: manaCost{c: 1, r: 1},
        },
        spellAbility: SpellAbility{
            ability{
                targets: []targetType{eachPlayer},
                effect:  damage{4},
            },
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
        flameRift.name:        flameRift,
		falkenrathReaver.name: falkenrathReaver,
	}

	deckList = unorderedCards{
		mountain:         12,
		lavaSpike:        9,
        flameRift:        3,
		falkenrathReaver: 6,
	}
)


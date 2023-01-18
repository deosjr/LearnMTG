package main

var (
	mountain = &land{
		card: card{
			name: "Mountain",
            activatedAbilities: []ActivatedAbility{
                {
                    cost: cost{tap:true},
                    ability: ability{
                        targets: []targetType{you},
                        effect:  addMana{amount: mana{r: 1}},
                    },
                },
            },
		},
	}

	island = &land{
		card: card{
			name: "Island",
            activatedAbilities: []ActivatedAbility{
                {
                    cost: cost{tap:true},
                    ability: ability{
                        targets: []targetType{you},
                        effect:  addMana{amount: mana{u: 1}},
                    },
                },
            },
		},
	}

	lavaSpike = &sorcery{
		card: card{
			name:     "Lava Spike",
			manaCost: mana{r: 1},
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
            manaCost: mana{c: 1, r: 1},
        },
        spellAbility: SpellAbility{
            ability{
                targets: []targetType{eachPlayer},
                effect:  damage{4},
            },
        },
    }

    divination = &sorcery{
        card: card{
            name:     "Divination",
            manaCost: mana{c: 2, u: 1},
        },
        spellAbility: SpellAbility{
            ability{
                targets: []targetType{you},
                effect:  draw{2},
            },
        },
    }

	falkenrathReaver = &creature{
		card: card{
			name:     "Falkenrath Reaver",
			manaCost: mana{c: 1, r: 1},
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
		mountain:         7,
        island:           7,
		lavaSpike:        4,
        flameRift:        4,
		falkenrathReaver: 4,
        divination:       4,
	}
)


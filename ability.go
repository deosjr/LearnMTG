package main

// 113.3 There are four general categories of abilities:
// 113.3a Spell abilities [...] followed as instructions while an instant or sorcery spell is resolving
// 113.3b Activated abilities have a cost and an effect.
// 113.3c Triggered abilities have a trigger condition and an effect.
// 113.3d Static abilities are written as statements. They're simply true.
// Keyword abilities are simply shortcuts for one of the above

type Ability interface {
    getTargets() []targetType
    getEffect() Effect
}

type ability struct {
    targets []targetType
    effect Effect
}

func (a ability) getTargets() []targetType {
    return a.targets
}

func (a ability) getEffect() Effect {
    return a.effect
}

type SpellAbility struct {
    ability
}

type ActivatedAbility struct {
    ability
    cost cost
}

// TODO: 605.1a An activated ability is a mana ability if it meets the following criteria:
// - it doesn't require a target
// - it could add mana to a player's mana pool when it resolves
// - it is not a loyalty ability
func (aa ActivatedAbility) isManaAbility() bool {
    if len(aa.targets) == 0 {
        return true
    }
    noTarget := len(aa.targets) == 1 && aa.targets[0].isUntargeted()
    return noTarget // && addsMana && noLoyalty
}

type TriggeredAbility struct {
    ability
    //trigger trigger
}

// TODO 605.1b similar for triggered abilities

type StaticAbility struct {
    ability
    // ???
}

type targetType int

const (
    you targetType = iota
    targetPlayer
    eachPlayer
)

func (t targetType) isPlayer() bool {
    return t == you || t == targetPlayer || t == eachPlayer
}

func (t targetType) isUntargeted() bool {
    return t == you || t == eachPlayer
}

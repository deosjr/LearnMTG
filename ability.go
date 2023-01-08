package main

// 113.3 There are four general categories of abilities:
// 113.3a Spell abilities [...] followed as instructions while an instant or sorcery spell is resolving
// 113.3b Activated abilities have a cost and an effect.
// 113.3c Triggered abilities have a trigger condition and an effect.
// 113.3d Static abilities are written as statements. They're simply true.
// Keyword abilities are simply shortcuts for one of the above

type Ability interface {
    getTargets() []targettype
    getEffect() Effect
}

type ability struct {
    targets []targettype
    effect Effect
}

func (a ability) getTargets() []targettype {
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
    //cost cost
}

type TriggeredAbility struct {
    ability
    //trigger trigger
}

type StaticAbility struct {
    ability
    // ???
}

type targettype int

const (
    selfTarget targettype = iota
    playerTarget
)

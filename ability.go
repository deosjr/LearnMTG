package main

// 113.3 There are four general categories of abilities:
// 113.3a Spell abilities [...] followed as instructions while an instant or sorcery spell is resolving
// 113.3b Activated abilities have a cost and an effect.
// 113.3c Triggered abilities have a trigger condition and an effect.
// 113.3d Static abilities are written as statements. They're simply true.
// Keyword abilities are simply shortcuts for one of the above
// I will name the first (spell ability) an 'effect', since typically 
// activated and triggered abilities result in effects equivalent to some spells

type ActivatedAbility struct {
    //cost cost
    effect Effect
}

type TriggeredAbility struct {
    //trigger trigger
    effect Effect
}

type StaticAbility struct {
    // ???
}

package game

// Implementation is what players should implement
type Implementation func() Player

type Player interface {
	// Name returns the name of the player
	Name() string

	// Start is called for the starting player to begin the game
	Start() Pokémon

	// Play receives the Pokémon the other player picked. Returning a non-existing pokémon id (such as -1) is considered a forfeit.
	Play(Pokémon) Pokémon
}

type Pokémon int

func (p Pokémon) Valid() bool {
	_, valid := names[p]
	return valid
}

func (p Pokémon) Name() string {
	return names[p]
}

func (p Pokémon) Start() rune {
	return []rune(names[p])[0]
}

func (p Pokémon) End() rune {
	runes := []rune(names[p])
	return runes[len(runes)-1]
}

func (p Pokémon) String() string {
	if p.Name() != `` {
		return p.Name()
	}

	return `[INVALID]`
}

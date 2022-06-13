package dummy

import (
	"math/rand"
	"strconv"

	"github.com/PotatoesFall/pokegame/game"
)

type dummy struct {
	name string
}

func (dummy) Play(game.Pokémon) game.Pokémon {
	return game.Pokémon(rand.Intn(9) + 1)
}

func (dummy) Start() game.Pokémon {
	return game.Pokémon(1)
}

func (d dummy) Name() string {
	return d.name
}

var count = 0

func New() game.Player {
	count++
	return dummy{
		name: `dummy ` + strconv.Itoa(count),
	}
}

package remote

import "github.com/PotatoesFall/pokegame/game"

const (
	messageTypeNewGame     = `new-game`
	messageTypeNameRequest = `name-request`
	messageTypeName        = `name`
	messageTypeYourTurn    = `your-turn`
	messageTypeMyTurn      = `my-turn`
	messageTypeGameOver    = `game-over`
)

type yourTurnMessage struct {
	Prev game.Pokémon
}

type myTurnMessage struct {
	Next game.Pokémon
}

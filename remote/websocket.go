package remote

import (
	"encoding/json"

	"github.com/PotatoesFall/pokegame/game"
	"github.com/gorilla/websocket"
)

func send(conn *websocket.Conn, typ messageType, msg any) {
	var data []byte
	if msg != nil {
		jsn, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}
		data = jsn
	}
	fullMsg := message{
		Typ: typ,
		Msg: json.RawMessage(data),
	}
	if err := conn.WriteJSON(fullMsg); err != nil {
		panic(err)
	}
}

func readJSON[T any](msg json.RawMessage) T {
	data, _ := msg.MarshalJSON()
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		panic(err)
	}
	return t
}

type messageType string

const (
	messageTypeNewGame     messageType = `new-game`
	messageTypeNameRequest messageType = `name-request`
	messageTypeName        messageType = `name`
	messageTypeYourTurn    messageType = `your-turn`
	messageTypeMyTurn      messageType = `my-turn`
	messageTypeGameOver    messageType = `game-over`
)

type message struct {
	Typ messageType
	Msg json.RawMessage
}

type yourTurnMessage struct {
	Prev game.Pokémon
}

type myTurnMessage struct {
	Next game.Pokémon
}

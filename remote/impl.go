package remote

import (
	"encoding/json"
	"os"

	"github.com/PotatoesFall/pokegame/game"
	"github.com/gorilla/websocket"
)

func StartImplementation(impl game.Implementation, serverEndpoint string) {
	conn, _, err := websocket.DefaultDialer.Dial(serverEndpoint, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	h := handler{
		impl: impl,
		conn: conn,
	}

	for {
		var msg message
		err := conn.ReadJSON(&msg)
		if err != nil {
			panic(err)
		}

		go h.process(msg)
	}
}

type handler struct {
	impl   game.Implementation
	player game.Player
	conn   *websocket.Conn
}

func (h handler) process(msg message) {
	switch msg.typ {
	case messageTypeGameOver:
		os.Exit(0)

	case messageTypeMyTurn:
		yourTurn := readJSON[yourTurnMessage](msg.msg)
		go h.handleYourTurn(yourTurn)

	case messageTypeNameRequest:
		go h.handleNameRequest()

	case messageTypeNewGame:
		go h.handleNewGame()

	default:
		panic(msg.typ)
	}
}

func (h handler) handleYourTurn(msg yourTurnMessage) {
	var resp myTurnMessage
	if !msg.Prev.Valid() {
		resp.Next = h.player.Start()
	} else {
		resp.Next = h.player.Play(msg.Prev)
	}

	send(h.conn, messageTypeMyTurn, resp)
}

func (h *handler) handleNewGame() {
	h.player = h.impl()
}

func (h handler) handleNameRequest() {
	name, _ := json.Marshal(h.player.Name())

	send(h.conn, messageTypeName, json.RawMessage(name))
}

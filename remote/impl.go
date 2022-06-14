package remote

import (
	"encoding/json"
	"os"
	"time"

	"github.com/PotatoesFall/pokegame/game"
	"github.com/gorilla/websocket"
)

func StartImplementation(impl game.Implementation, serverEndpoint string) {
	conn, _, err := websocket.DefaultDialer.Dial(serverEndpoint, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	h := client{
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

type client struct {
	impl   game.Implementation
	player game.Player
	conn   *websocket.Conn
}

func (h *client) process(msg message) {
	switch msg.Typ {
	case messageTypeGameOver:
		os.Exit(0)

	case messageTypeYourTurn:
		yourTurn := readJSON[yourTurnMessage](msg.Msg)
		go h.handleYourTurn(yourTurn)

	case messageTypeNameRequest:
		go h.handleNameRequest()

	case messageTypeNewGame:
		go h.handleNewGame()

	default:
		panic(msg.Typ)
	}
}

func (h *client) handleYourTurn(msg yourTurnMessage) {
	time.Sleep(10 * time.Millisecond) // TODO FIX THIS

	var resp myTurnMessage
	if !msg.Prev.Valid() {
		resp.Next = h.player.Start()
	} else {
		resp.Next = h.player.Play(msg.Prev)
	}

	send(h.conn, messageTypeMyTurn, resp)
}

func (h *client) handleNewGame() {
	h.player = h.impl()
}

func (h *client) handleNameRequest() {
	name, _ := json.Marshal(h.player.Name())

	send(h.conn, messageTypeName, name)
}

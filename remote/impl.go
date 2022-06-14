package remote

import (
	"fmt"
	"os"

	"github.com/PotatoesFall/pokegame/game"
	"github.com/PotatoesFall/pokegame/remote/socket"
)

func StartImplementation(impl game.Implementation, serverEndpoint string) {
	c := client{
		impl: impl,
	}

	handlers := socket.Handlers{}
	handlers.Register(messageTypeNewGame, socket.NewHandler(c.handleNewGame))
	handlers.Register(messageTypeYourTurn, socket.NewHandler(c.handleYourTurn))
	handlers.Register(messageTypeNameRequest, socket.NewHandler(c.handleNameRequest))
	handlers.Register(messageTypeNewGame, socket.NewHandler(c.handleNewGame))

	conn, err := socket.NewConn(serverEndpoint, handlers, func() {
		fmt.Println(`connection closed.`)
		os.Exit(0)
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(`Connection to server established.`)

	c.conn = conn
}

type client struct {
	impl   game.Implementation
	player game.Player
	conn   socket.Conn
}

func (c *client) handleYourTurn(msg yourTurnMessage) {
	if msg.Prev.Valid() {
		fmt.Println("opponent:\t" + msg.Prev.Name())
	}

	var resp myTurnMessage
	if !msg.Prev.Valid() {
		resp.Next = c.player.Start()
	} else {
		resp.Next = c.player.Play(msg.Prev)
	}

	fmt.Println("you:\t\t" + resp.Next.Name())
	if err := c.conn.Send(messageTypeMyTurn, resp); err != nil {
		fmt.Println(`error sending turn`, err)
	}
}

func (c *client) handleNewGame(any) {
	c.player = c.impl()
}

func (c *client) handleNameRequest(any) {
	if err := c.conn.Send(messageTypeName, c.player.Name()); err != nil {
		fmt.Println(`error sending name`, err)
	}
}

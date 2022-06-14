package remote

import (
	"fmt"

	"github.com/PotatoesFall/pokegame/game"
	"github.com/PotatoesFall/pokegame/remote/socket"
)

func AwaitImplementation(s socket.Server) (impl game.Implementation, close func()) {
	h := socket.Handlers{}
	conn := s.AwaitConnection(h, func() {})

	// TODO: currently an implementation will not work with multiple concurrent players
	return func() game.Player {
		p := player{
			nameChan: make(chan string),
			turnChan: make(chan myTurnMessage),
		}

		h.Register(messageTypeMyTurn, socket.NewHandler(p.handleMyTurn))
		h.Register(messageTypeName, socket.NewHandler(p.handleName))

		p.conn = conn
		if err := p.conn.Send(messageTypeNewGame, nil); err != nil {
			fmt.Println(`error starting new game`, err)
		}

		return p
	}, conn.Close
}

type player struct {
	conn socket.Conn
	name string

	nameChan chan string
	turnChan chan myTurnMessage
}

func (p player) handleMyTurn(msg myTurnMessage) {
	p.turnChan <- msg
}

func (p player) handleName(name string) {
	p.nameChan <- name
}

func (p player) Name() string {
	if p.name == `` {
		p.name = p.getName()
	}

	return p.name
}

func (p player) getName() string {
	if err := p.conn.Send(messageTypeNameRequest, nil); err != nil {
		fmt.Println(`error sending name request`, err)
	}
	name := <-p.nameChan
	return name
}

func (p player) Start() game.Pokémon {
	if err := p.conn.Send(messageTypeYourTurn, yourTurnMessage{
		Prev: game.Pokémon(-1),
	}); err != nil {
		fmt.Println(`error sending your turn`, err)
	}

	resp := <-p.turnChan
	return resp.Next
}

func (p player) Play(prev game.Pokémon) game.Pokémon {
	if err := p.conn.Send(messageTypeYourTurn, yourTurnMessage{
		Prev: prev,
	}); err != nil {
		fmt.Println(`error sending your turn`, err)
	}

	resp := <-p.turnChan
	return resp.Next
}

func (p player) GameOver(won bool) {
	if err := p.conn.Send(messageTypeGameOver, gameOverMessage{won}); err != nil {
		fmt.Println(`error sending game over`, err)
	}
}

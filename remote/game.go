package remote

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/PotatoesFall/pokegame/game"
	"github.com/gorilla/websocket"
)

func WaitForConnections(port int) {
	listener, err := net.Listen(`tcp`, `:`+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	port = listener.Addr().(*net.TCPAddr).Port
	fmt.Println(`awaiting connections on port ` + strconv.Itoa(port))

	go func() {
		panic(http.Serve(listener, socketHandler{}))
	}()
}

type socketHandler struct{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (socketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	newConnections <- conn
}

var newConnections = make(chan *websocket.Conn)

func AwaitImplementation() game.Implementation {
	conn := <-newConnections

	return func() game.Player {
		p := player{
			conn:     conn,
			nameChan: make(chan string),
			turnChan: make(chan myTurnMessage),
		}

		return p
	}
}

type player struct {
	conn *websocket.Conn
	name string

	nameChan chan string
	turnChan chan myTurnMessage
}

func (p player) Name() string {
	if p.name == `` {
		p.name = p.getName()
	}

	return p.name
}

func (p player) getName() string {
	send(p.conn, messageTypeNameRequest, nil)
	name := <-p.nameChan
	return name
}

func (p player) Start() game.Pokémon {
	send(p.conn, messageTypeYourTurn, yourTurnMessage{
		Prev: game.Pokémon(-1),
	})

	resp := <-p.turnChan
	return resp.Next
}

func (p player) Play(prev game.Pokémon) game.Pokémon {
	send(p.conn, messageTypeYourTurn, yourTurnMessage{
		Prev: prev,
	})

	resp := <-p.turnChan
	return resp.Next
}

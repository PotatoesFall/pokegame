package socket

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type genericMessage[T any] struct {
	Type    string
	Content T
}

type Server interface {
	Start(listener net.Listener)
	AwaitConnection(h Handlers, onClose func()) Conn
}

func NewServer() Server {
	return &server{
		newConnections: make(chan *websocket.Conn),
	}
}

func (s *server) Start(listener net.Listener) {
	go func() {
		panic(http.Serve(listener, s))
	}()
}

func NewConn(endpoint string, h Handlers, onClose func()) (Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return Conn{}, err
	}

	go handle(conn, h, onClose)

	return Conn{
		conn: conn,
	}, nil
}

type Conn struct {
	conn *websocket.Conn
}

func (c Conn) Send(msgType string, msg any) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	genMsg := genericMessage[json.RawMessage]{
		Type:    msgType,
		Content: data,
	}

	// data, _ = json.Marshal(genMsg)
	// fmt.Println(`DEBUG-send`, msgType, string(data))
	// fmt.Printf("%s\n", msg)

	return c.conn.WriteJSON(genMsg)
}

func (c Conn) Close() {
	c.conn.Close()
}

type server struct {
	newConnections chan *websocket.Conn
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	s.newConnections <- conn
}

func (s *server) AwaitConnection(h Handlers, onClose func()) Conn {
	conn := <-s.newConnections

	go handle(conn, h, onClose)

	return Conn{
		conn: conn,
	}
}

type Handler func(json.RawMessage)

type Handlers map[string]Handler

func (h Handlers) Register(msgType string, handler Handler) {
	h[msgType] = handler
}

func NewHandler[T any](h func(T)) func(json.RawMessage) {
	return func(m json.RawMessage) {
		data, _ := m.MarshalJSON()
		var t T
		if err := json.Unmarshal(data, &t); err != nil {
			err = fmt.Errorf(`unable to unmarshal %q into type %T: %w`, string(data), t, err)
			fmt.Println(err.Error())
			return
		}
		h(t)
	}
}

func handle(conn *websocket.Conn, h Handlers, onClose func()) {
	for {
		var msg genericMessage[json.RawMessage]
		err := conn.ReadJSON(&msg)
		if err != nil {
			if onClose != nil {
				onClose()
			} else {
				panic(err)
			}
		}

		// fmt.Println(`DEBUG-receive`, msg.Type, string(msg.Content))
		handler, ok := h[msg.Type]
		if !ok {
			fmt.Printf("websocket: no handler registered for message of type %q\n", msg.Type)
			return
		}

		handler(msg.Content)
	}
}

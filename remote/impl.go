package remote

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/PotatoesFall/pokegame/game"
)

func StartImplementationServer(impl game.Implementation) {
	listener, err := net.Listen(`tcp`, `:0`)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting server on port %d\n", listener.Addr().(*net.TCPAddr).Port)
	go func() game.Player {
		panic(http.Serve(listener, handler{impl: impl, sessions: make(map[int]game.Player)}))
	}()
}

type handler struct {
	impl     game.Implementation
	sessions map[int]game.Player
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)

	switch r.URL.Path {
	case `/new`:
		h.handleNew(w, r)
	case `/name`:
		h.handleName(w, r)
	case `/start`:
		h.handleStart(w, r)
	case `/play`:
		h.handlePlay(w, r)
	}
}

func (h handler) handleNew(w http.ResponseWriter, r *http.Request) {
	req := scanJSON[newGameRequest](r)
	h.sessions[req.SessionID] = h.impl()
	respond(w, nil)
}

func (h handler) handleName(w http.ResponseWriter, r *http.Request) {
	sessionID := scanJSON[int](r)
	respond(w, h.sessions[sessionID].Name())
}

func (h handler) handleStart(w http.ResponseWriter, r *http.Request) {
	sessionID := scanJSON[int](r)
	respond(w, h.sessions[sessionID].Start())
}

func (h handler) handlePlay(w http.ResponseWriter, r *http.Request) {
	req := scanJSON[playRequest](r)
	respond(w, h.sessions[req.SessionID].Play(req.Prev))
}

func scanJSON[T any](r *http.Request) T {
	var t T
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		panic(err)
	}
	return t
}

func respond(w http.ResponseWriter, resp any) {
	if resp != nil {
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			panic(err)
		}
	}
}

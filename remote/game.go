package remote

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/PotatoesFall/pokegame/game"
)

type newGameRequest struct {
	SessionID int
}

type playRequest struct {
	SessionID int
	Prev      game.Pokémon
}

func NewImplementation(endpoint string) game.Implementation {
	return func() game.Player {
		sessionID := rand.Intn(999999) + 1

		p := player{
			endpoint:  endpoint,
			sessionID: sessionID,
		}

		p.post(`/new`, newGameRequest{
			SessionID: sessionID,
		}, nil)

		return p
	}
}

type player struct {
	endpoint  string
	sessionID int
	name      string
}

func (p player) Name() string {
	if p.name == `` {
		p.name = p.getName()
	}

	return p.name
}

func (p player) getName() string {
	var name string
	p.post(`/name`, p.sessionID, &name)
	return name
}

func (p player) Start() game.Pokémon {
	var pok game.Pokémon
	p.post(`/start`, p.sessionID, &pok)
	return pok
}

func (p player) Play(prev game.Pokémon) game.Pokémon {
	var pok game.Pokémon
	p.post(`/play`, playRequest{
		SessionID: p.sessionID,
		Prev:      prev,
	}, &pok)
	return pok
}

func (p player) post(path string, body any, dst any) {
	var data []byte
	if body != nil {
		j, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		data = j
	}

	resp, err := http.DefaultClient.Post(p.endpoint+path, `application/json`, bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}

	if dst != nil {
		err = json.NewDecoder(resp.Body).Decode(dst)
		if err != nil {
			panic(err)
		}
	}
}

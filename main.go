package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/PotatoesFall/pokegame/game"
	"github.com/PotatoesFall/pokegame/remote"
	"github.com/PotatoesFall/pokegame/remote/socket"
)

func main() {
	portEnv, _ := os.LookupEnv(`WS_PORT`)
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		panic(err)
	}
	listener, err := net.Listen(`tcp`, `:`+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	port = listener.Addr().(*net.TCPAddr).Port

	s := socket.NewServer()
	s.Start(listener)
	fmt.Println(`awaiting connections on port ` + strconv.Itoa(port))

	for {
		impl1, close1 := remote.AwaitImplementation(s)
		impl2, close2 := remote.AwaitImplementation(s)
		runGame(impl1, impl2)
		close1()
		close2()
	}
}

func runGame(impl1, impl2 game.Implementation) {
	fmt.Println(`GAME 1`)
	p1, p2 := impl1(), impl2()
	winner1 := playGame(p1, p2)

	fmt.Println()
	p1, p2 = impl1(), impl2()
	fmt.Println(`GAME 2`)
	winner2 := playGame(p2, p1)
	fmt.Println()

	if winner1 == winner2 {
		fmt.Println(`It's a tie! Both programs lost and won once.`)
	}

	if winner1 == 1 {
		fmt.Printf("%s won both games!\n", p1.Name())
	} else if winner1 == 2 {
		fmt.Printf("%s won both games!\n", p2.Name())
	} else {
		panic(`unable to figure out who won: ` + strconv.Itoa(winner1))
	}
}

func playGame(p1, p2 game.Player) int {
	p := p1.Start()
	fmt.Printf("%s started with %s\n", p1.Name(), p)
	if !p.Valid() {
		fmt.Printf("That's not a valid Pokémon! %s loses!\n", p1.Name())
		return 2
	}
	used := map[game.Pokémon]bool{
		p: true,
	}
	c := p.End()
	currentPlayer, otherPlayer := p2, p1
	winner := 0

	// allNames := game.AllNames()
	for turn(&p, used, currentPlayer, &c) {
		// if noMoreAnswers(allNames, used, c) {
		// 	fmt.Printf("There are no correct answers left! %s loses!\n", otherPlayer.Name())
		// 	break
		// }
		currentPlayer, otherPlayer = otherPlayer, currentPlayer
		winner = 1 - winner
	}

	p1.GameOver(winner == 0)
	p2.GameOver(winner == 1)

	return winner + 1
}

// func noMoreAnswers(allNames []game.Pokémon, used map[game.Pokémon]bool, c rune) bool {
// 	for _, pok := range game.AllNames() {
// 		if pok.Start() == c && !used[pok] {
// 			return false
// 		}
// 	}

// 	return true
// }

func turn(p *game.Pokémon, used map[game.Pokémon]bool, p1 game.Player, c *rune) bool {
	*p = p1.Play(*p)
	fmt.Printf("%s: %s\n", p1.Name(), p)
	if !p.Valid() {
		fmt.Printf("That's not a valid Pokémon! %s loses!\n", p1.Name())
		return false

	}
	if used[*p] {
		fmt.Printf("That Pokémon was already used! %s loses!\n", p1.Name())
		return false
	}
	used[*p] = true
	if p.Start() != *c {
		fmt.Printf("%s starts with %c, not %c! %s loses!\n", p, p.Start(), *c, p1.Name())
		return false
	}
	*c = p.End()

	return true
}

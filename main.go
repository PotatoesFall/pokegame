package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/PotatoesFall/pokegame/dummy"
	"github.com/PotatoesFall/pokegame/game"
)

func main() {
	rand.Seed(time.Now().UnixMicro())
	runGame(dummy.New, dummy.New)
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
		exit(`It's a tie! Both programs lost and won once.`)
	}

	if winner1 == 1 {
		exit(fmt.Sprintf(`%s won both games!`, p1.Name()))
	}

	if winner1 == 2 {
		exit(fmt.Sprintf(`%s won both games!`, p2.Name()))
	}

	panic(`unable to figure out who won`)
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(0)
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

	allNames := game.AllNames()
	for turn(&p, used, currentPlayer, &c) {
		if noMoreAnswers(allNames, used, c) {
			fmt.Printf("There are no correct answers left! %s loses!\n", otherPlayer.Name())
			break
		}
		currentPlayer, otherPlayer = otherPlayer, currentPlayer
		winner = 1 - winner
	}

	return winner + 1
}

func noMoreAnswers(allNames []game.Pokémon, used map[game.Pokémon]bool, c rune) bool {
	for _, pok := range game.AllNames() {
		if pok.Start() == c && !used[pok] {
			return false
		}
	}

	return true
}

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

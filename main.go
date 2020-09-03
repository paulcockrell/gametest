package main

import (
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Screen constants
const (
	screenWidth  = 240
	screenHeight = 240
)

type Game struct {
	runner *Runner
	level  *Level
}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.runner = NewRunner(screenWidth/2, screenHeight/2)
	g.level = NewLevel("resources/levels/level_one.json")
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.runner.update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.level.draw(screen)
	g.runner.draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Jon Wicker 3 - Parallelagram")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatalf("error starting game %v", err)
	}
}

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

type Sprite struct {
	image                   *ebiten.Image
	numFrames               int
	frameOX, frameOY        int
	frameHeight, frameWidth int
}

type Game struct {
	vaxerman *VaxerMan
	level    *Level
	enemies  []*Enemy
}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.vaxerman = NewVaxerMan(screenWidth/2, screenHeight/2)
	g.level = NewLevel("resources/levels/level_one.json")
	g.enemies = make([]*Enemy, 0)
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.vaxerman.update()
	g.updateEnemies()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.level.draw(screen)
	g.vaxerman.draw(screen)
	g.vaxerman.drawBullets(screen)
	for _, enemy := range g.enemies {
		enemy.draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) updateEnemies() {
	currentEnemies := make([]*Enemy, 0)
	for _, enemy := range g.enemies {
		if enemy.HasInfectedPlayer(g.vaxerman) {
			g.vaxerman.SetDead()
		}
		if g.vaxerman.hasShotEnemy(enemy) {
			enemy.status = EnemyHit
		} else {
			enemy.update()
		}
		if enemy.status != EnemyDead {
			currentEnemies = append(currentEnemies, enemy)
		}
	}
	g.enemies = currentEnemies

	if len(g.enemies) < MaxEnemies && (rand.Intn(5) == 1) {
		x, y, vx, vy := GenerateEnemyStartPos()
		newEnemy := NewEnemy(x, y, vx, vy)
		g.enemies = append(g.enemies, newEnemy)
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("VaxerMan - Corona Virus Killer")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatalf("error starting game %v", err)
	}
}

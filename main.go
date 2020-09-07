package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Screen constants
const (
	screenWidth   = 240
	screenHeight  = 240
	fontSize      = 12
	smallFontSize = fontSize / 2
)

var (
	arcadeFont      font.Face
	smallArcadeFont font.Face
)

func init() {
	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	smallArcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    smallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

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
	// If VaxerMan is infected, activate the "R" key to
	// reset the game
	if g.vaxerman.IsDead() && ebiten.IsKeyPressed(ebiten.KeyR) {
		g.init()
		return nil
	}

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
	g.drawInfo(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) updateEnemies() {
	currentEnemies := make([]*Enemy, 0)
	for _, enemy := range g.enemies {
		enemy.update()

		if enemy.HasInfectedPlayer(g.vaxerman) {
			g.vaxerman.Infect()
		}

		if g.vaxerman.hasShotEnemy(enemy) {
			enemy.status = EnemyHit
		}

		if !enemy.IsDead() {
			currentEnemies = append(currentEnemies, enemy)
		}
	}
	g.enemies = currentEnemies

	if len(g.enemies) < MaxEnemies && (rand.Intn(20) == 1) {
		x, y, vx, vy := GenerateEnemyStartPos()
		newEnemy := NewEnemy(x, y, vx, vy)
		g.enemies = append(g.enemies, newEnemy)
	}
}

func (g *Game) drawInfo(screen *ebiten.Image) {
	if g.vaxerman.IsDead() {
		texts := []string{"VaxerMan is too ill to fight on..", "", "", "", "Press R key to restart"}
		for i, l := range texts {
			x := (screenWidth - len(l)*fontSize) / 3
			text.Draw(screen, l, smallArcadeFont, x, (i+4)*fontSize, color.White)
		}
	}
	lives := fmt.Sprintf("Lives: %d", g.vaxerman.lives)
	text.Draw(screen, lives, smallArcadeFont, 150, 12, color.White)
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("VaxerMan - Corona Virus Killer")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatalf("error starting game %v", err)
	}
}

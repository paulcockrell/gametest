package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/paulcockrell/gametest/resources/images"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Screen constants
const (
	screenWidth  = 240
	screenHeight = 240
)

// Runner left/right constants
const (
	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameNum    = 8
)

var (
	runnerRightImage *ebiten.Image
	runnerLeftImage  *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(images.RunnerRight_png))
	if err != nil {
		log.Fatalf("error decoding image: %v", err)
	}
	runnerRightImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	img, _, err = image.Decode(bytes.NewReader(images.RunnerLeft_png))
	if err != nil {
		log.Fatalf("error decoding image: %v", err)
	}
	runnerLeftImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

type State int

const (
	StateIdle = iota
	StateRun
)

type SpriteSettings struct {
	frameCount       int
	frameOX, frameOY int
}

type Action int

const (
	ActionLeft Action = iota
	ActionRight
)

type runner struct {
	x, y       int
	vx, vy     int
	frameCount int
	action     Action
	state      State
}

func (r *runner) update() {
	r.x += r.vx
	r.y += r.vy
}

func (r *runner) draw(screen *ebiten.Image) {
	sprite := runnerLeftImage
	switch r.action {
	case ActionLeft:
		sprite = runnerLeftImage
	case ActionRight:
		sprite = runnerRightImage
	}

	op := &ebiten.DrawImageOptions{}
	w, h := sprite.Size()
	op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.GeoM.Translate(float64(r.x), float64(r.y))

	// Extract sprite frame

	i := (r.frameCount / 5) % 5       //frameNum
	sx, sy := frameOX+i*frameWidth, 0 //frameOY
	spriteSubImage := sprite.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image)

	screen.DrawImage(spriteSubImage, op)
}

// Tile constants
const (
	tileSize = 16
	tileXNum = 25
)

var (
	tilesImage *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatalf("error loading tiles images: %v", err)
	}
	tilesImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

type Game struct {
	// Main character
	runner *runner

	// Background layers
	layers [][]int
}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.runner = &runner{x: screenWidth / 2, y: screenHeight / 2, frameCount: 0}
	g.layers = [][]int{
		{
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 244, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 219, 243, 243, 243, 219, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		},
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 26, 27, 28, 29, 30, 31, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 51, 52, 53, 54, 55, 56, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 76, 77, 78, 79, 80, 81, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 101, 102, 103, 104, 105, 106, 0, 0, 0, 0,

			0, 0, 0, 0, 0, 126, 127, 128, 129, 130, 131, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 303, 303, 245, 242, 303, 303, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,

			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		},
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	const moveBy = 2

	g.runner.frameCount++
	g.runner.vx = 0
	g.runner.vy = 0

	// H - Left
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		if g.runner.x > 0 {
			g.runner.vx -= moveBy
		}
	}

	// L - Right
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		if g.runner.x < screenWidth-(frameWidth/2) {
			g.runner.vx += moveBy
		}
	}

	// K - Up
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		if g.runner.y > -(screenHeight/2)+frameHeight+10 {
			g.runner.vy -= moveBy
		}
	}

	// J - Down
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		if g.runner.y < screenHeight-(frameHeight*3) {
			g.runner.vy += moveBy
		}
	}

	g.runner.update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawTiles(screen)
	g.runner.draw(screen)
}

func (g *Game) drawTiles(screen *ebiten.Image) {
	const xNum = screenWidth / tileSize
	for _, l := range g.layers {
		for i, t := range l {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64((i%xNum)*tileSize), float64((i/xNum)*tileSize))

			sx := (t % tileXNum) * tileSize
			sy := (t / tileXNum) * tileSize
			tile := tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image)
			screen.DrawImage(tile, op)
		}
	}
}

/*
func (g *Game) drawRunner(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := runnerImage.Size()
	op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.GeoM.Translate(float64(g.x16), float64(g.y16))

	// Extract sprite frame
	i := (g.count / 5) % frameNum
	sx, sy := frameOX+i*frameWidth, frameOY
	runnerSprite := runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image)

	screen.DrawImage(runnerSprite, op)
}
*/

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

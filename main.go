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

type Sprite struct {
	image                   *ebiten.Image
	numFrames               int
	frameOX, frameOY        int
	frameHeight, frameWidth int
}

type Action int

const (
	ActionLeftIdle Action = iota
	ActionLeftRun
	ActionRightIdle
	ActionRightRun
)

type runner struct {
	x, y       int
	vx, vy     int
	frameCount int
	action     Action
	sprites    map[Action]Sprite
}

func newRunner(x, y int) *runner {
	r := &runner{
		x: x,
		y: y,
	}
	r.sprites = map[Action]Sprite{
		ActionLeftIdle: {
			image:       runnerLeftImage,
			numFrames:   5,
			frameOX:     0,
			frameOY:     0,
			frameHeight: 32,
			frameWidth:  32,
		},
		ActionLeftRun: {
			image:       runnerLeftImage,
			numFrames:   8,
			frameOX:     0,
			frameOY:     32,
			frameHeight: 32,
			frameWidth:  32,
		},
		ActionRightIdle: {
			image:       runnerRightImage,
			numFrames:   5,
			frameOX:     0,
			frameOY:     0,
			frameHeight: 32,
			frameWidth:  32,
		},
		ActionRightRun: {
			image:       runnerRightImage,
			numFrames:   8,
			frameOX:     0,
			frameOY:     32,
			frameHeight: 32,
			frameWidth:  32,
		},
	}

	return r
}

func (r *runner) update() {
	const moveBy = 2

	r.frameCount++
	r.vx = 0
	r.vy = 0

	// Reset action to idling for last direction incase no keypress detected
	if r.action == ActionLeftRun {
		r.action = ActionLeftIdle
	}
	if r.action == ActionRightRun {
		r.action = ActionRightIdle
	}

	// H - Left
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		r.action = ActionLeftRun
		if r.x > 0 {
			r.vx -= moveBy
		}
	}

	// L - Right
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		r.action = ActionRightRun
		sprite := r.sprites[r.action]
		if r.x < screenWidth-(sprite.frameWidth/2) {
			r.vx += moveBy
		}
	}

	// K - Up
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		sprite := r.sprites[r.action]
		if r.y > -(screenHeight/2)+sprite.frameHeight+10 {
			r.vy -= moveBy
		}
	}

	// J - Down
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		sprite := r.sprites[r.action]
		if r.y < screenHeight-(sprite.frameHeight*3) {
			r.vy += moveBy
		}
	}

	r.x += r.vx
	r.y += r.vy
}

func (r *runner) draw(screen *ebiten.Image) {
	sprite := r.sprites[r.action]

	op := &ebiten.DrawImageOptions{}
	w, h := sprite.image.Size()
	op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.GeoM.Translate(float64(r.x), float64(r.y))

	// Extract sprite frame
	i := (r.frameCount / sprite.numFrames) % sprite.numFrames
	sx, sy := sprite.frameOX+i*sprite.frameWidth, sprite.frameOY
	spriteSubImage := sprite.image.SubImage(image.Rect(sx, sy, sx+sprite.frameWidth, sy+sprite.frameHeight)).(*ebiten.Image)

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
	g.runner = newRunner(screenWidth/2, screenHeight/2)
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

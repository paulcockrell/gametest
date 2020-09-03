package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/paulcockrell/gametest/resources/images"
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

type Runner struct {
	x, y       int
	vx, vy     int
	frameCount int
	action     Action
	sprites    map[Action]Sprite
}

func NewRunner(x, y int) *Runner {
	r := &Runner{
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

func (r *Runner) update() {
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

func (r *Runner) draw(screen *ebiten.Image) {
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

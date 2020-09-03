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

type Actions uint8

const (
	Idle Actions = 1 << iota
	Run
	Shoot
	Left
	Right
	Up
	Down
)

func Set(a, flag Actions) Actions    { return a | flag }
func Clear(a, flag Actions) Actions  { return a &^ flag }
func Toggle(a, flag Actions) Actions { return a ^ flag }
func Has(a, flag Actions) bool       { return a&flag != 0 }

type Runner struct {
	x, y       int
	vx, vy     int
	frameCount int
	actions    Actions
	sprites    map[Actions]Sprite
}

func NewRunner(x, y int) *Runner {
	var a Actions = Idle | Right

	r := &Runner{
		x:       x,
		y:       y,
		actions: a,
	}

	r.sprites = map[Actions]Sprite{
		Left | Idle: {
			image:       runnerLeftImage,
			numFrames:   5,
			frameOX:     0,
			frameOY:     0,
			frameHeight: 32,
			frameWidth:  32,
		},
		Left | Run: {
			image:       runnerLeftImage,
			numFrames:   8,
			frameOX:     0,
			frameOY:     32,
			frameHeight: 32,
			frameWidth:  32,
		},
		Left | Shoot: {
			image:       runnerLeftImage,
			numFrames:   4,
			frameOX:     0,
			frameOY:     64,
			frameHeight: 32,
			frameWidth:  32,
		},
		Right | Idle: {
			image:       runnerRightImage,
			numFrames:   5,
			frameOX:     0,
			frameOY:     0,
			frameHeight: 32,
			frameWidth:  32,
		},
		Right | Run: {
			image:       runnerRightImage,
			numFrames:   8,
			frameOX:     0,
			frameOY:     32,
			frameHeight: 32,
			frameWidth:  32,
		},
		Right | Shoot: {
			image:       runnerRightImage,
			numFrames:   4,
			frameOX:     0,
			frameOY:     64,
			frameHeight: 32,
			frameWidth:  32,
		},
		Up | Idle: {
			image:       runnerRightImage,
			numFrames:   5,
			frameOX:     0,
			frameOY:     0,
			frameHeight: 32,
			frameWidth:  32,
		},
		Up | Run: {
			image:       runnerRightImage,
			numFrames:   8,
			frameOX:     0,
			frameOY:     32,
			frameHeight: 32,
			frameWidth:  32,
		},
		Up | Shoot: {
			image:       runnerRightImage,
			numFrames:   4,
			frameOX:     0,
			frameOY:     64,
			frameHeight: 32,
			frameWidth:  32,
		},
		Down | Idle: {
			image:       runnerLeftImage,
			numFrames:   5,
			frameOX:     0,
			frameOY:     0,
			frameHeight: 32,
			frameWidth:  32,
		},
		Down | Run: {
			image:       runnerLeftImage,
			numFrames:   8,
			frameOX:     0,
			frameOY:     32,
			frameHeight: 32,
			frameWidth:  32,
		},
		Down | Shoot: {
			image:       runnerLeftImage,
			numFrames:   4,
			frameOX:     0,
			frameOY:     64,
			frameHeight: 32,
			frameWidth:  32,
		},
	}

	return r
}

func (r *Runner) update() {
	const moveBy = 2

	// Reset velocity values
	r.vx = 0
	r.vy = 0

	// Reset movement to default idle
	if Has(r.actions, Left|Down) {
		r.actions = Left | Idle
	}
	if Has(r.actions, Right|Up) {
		r.actions = Right | Idle
	}

	// Update runner state based on keyboard input
	// H - Left
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		r.actions = Left | Run
		r.vx -= moveBy
	}

	// L - Right
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		r.actions = Right | Run
		r.vx += moveBy
	}

	// K - Up
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		r.actions = Up | Run
		r.vy -= moveBy
	}

	// J - Down
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		r.actions = Down | Run
		r.vy += moveBy
	}

	/* Wall collision detection */
	s := r.getSprite()
	if r.x < 0 {
		r.vx = 0
	}
	if r.x >= screenWidth-(s.frameWidth/2) {
		r.vx = screenWidth - (s.frameWidth / 2)
	}
	if r.y <= -(screenHeight/2)+s.frameHeight+10 {
		r.vy = -(screenHeight / 2) + s.frameHeight + 10
	}
	if r.y > screenHeight-(s.frameHeight*3) {
		r.vy = screenHeight - (s.frameHeight * 3)
	}

	// Update sprite's x & y positions based on velocity values and
	// frame counter used by animation
	r.frameCount++
	r.x += r.vx
	r.y += r.vy
}

func (r Runner) direction() Actions {
	switch {
	case Has(r.actions, Left):
		return Left
	case Has(r.actions, Right):
		return Right
	case Has(r.actions, Up):
		return Up
	case Has(r.actions, Down):
		return Down
	default:
		return Right
	}
}

func (r Runner) action() Actions {
	switch {
	case Has(r.actions, Idle):
		return Idle
	case Has(r.actions, Run):
		return Run
	case Has(r.actions, Shoot):
		return Shoot
	default:
		return Idle
	}
}

func (r Runner) getSprite() Sprite {
	direction := r.direction()
	action := r.action()

	return r.sprites[direction|action]
}

func (r *Runner) draw(screen *ebiten.Image) {
	sprite := r.getSprite()

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

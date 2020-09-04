package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/paulcockrell/gametest/resources/images"
)

var (
	runnerRightImage *ebiten.Image
	runnerLeftImage  *ebiten.Image
	bulletImage      *ebiten.Image
)

const (
	maxBullets = 3
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

	img, _, err = image.Decode(bytes.NewReader(images.Bullet_png))
	if err != nil {
		log.Fatalf("error decoding image: %v", err)
	}
	bulletImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
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

type Bullet struct {
	sprite  Sprite
	actions Actions
	x, y    int
}

func NewBullet(x, y int, a Actions) *Bullet {
	b := &Bullet{
		x:       x, // starting x
		y:       y, // starting y
		actions: a, // holds direction
	}
	b.sprite = Sprite{
		image:       bulletImage,
		numFrames:   1,
		frameOX:     0,
		frameOY:     0,
		frameHeight: 32,
		frameWidth:  32,
	}
	if Has(a, (Up | Down)) {
		b.sprite.frameOX = 1
	}

	return b
}

type Runner struct {
	x, y       int
	vx, vy     int
	frameCount int
	actions    Actions
	sprites    map[Actions]Sprite
	bullets    []*Bullet
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
	r.actions = r.actions &^ (Run | Shoot)
	r.actions = r.actions | Idle

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

	// SPACE - Spacebar
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		if len(r.bullets) < maxBullets {
			r.bullets = append(r.bullets, NewBullet(r.x, r.y, r.actions))

		}
		r.actions = r.actions &^ (Idle | Run)
		r.actions = r.actions | Shoot
	}

	// Update sprite's x & y positions based on velocity values and
	// frame counter used by animation
	r.frameCount++
	r.x += r.vx
	r.y += r.vy

	// Wall collision detection
	s := r.getSprite()
	if r.x < 0 {
		r.x = 0
	}
	//if r.x >= screenWidth-(s.frameWidth/2) {
	if r.x > screenWidth-s.frameWidth {
		r.x = screenWidth - s.frameWidth
	}
	if r.y < 0 {
		r.y = 0
	}
	if r.y > screenHeight-s.frameHeight {
		r.y = screenHeight - s.frameHeight
	}

	// Update bullets if any
	var activeBullets []*Bullet
	for _, bullet := range r.bullets {
		if Has(bullet.actions, Left) {
			bullet.x -= 5
		}
		if Has(bullet.actions, Right) {
			bullet.x += 5
		}
		if Has(bullet.actions, Up) {
			bullet.y -= 5
		}
		if Has(bullet.actions, Down) {
			bullet.y += 5
		}

		// If on screen keep bullet
		if bullet.x > 0 && bullet.x < screenWidth &&
			bullet.y > 0 && bullet.y < screenHeight {
			activeBullets = append(activeBullets, bullet)
		}
	}
	r.bullets = activeBullets
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
	op.GeoM.Translate(float64(r.x), float64(r.y))

	// Extract sprite frame
	i := (r.frameCount / sprite.numFrames) % sprite.numFrames
	sx, sy := sprite.frameOX+i*sprite.frameWidth, sprite.frameOY
	spriteSubImage := sprite.image.SubImage(image.Rect(sx, sy, sx+sprite.frameWidth, sy+sprite.frameHeight)).(*ebiten.Image)

	screen.DrawImage(spriteSubImage, op)
}

func (r *Runner) drawBullets(screen *ebiten.Image) {
	for _, bullet := range r.bullets {
		sprite := bullet.sprite

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(bullet.x), float64(bullet.y))

		// Extract sprite frame
		i := (1 / sprite.numFrames) % sprite.numFrames
		sx, sy := sprite.frameOX+i*sprite.frameWidth, sprite.frameOY
		spriteSubImage := sprite.image.SubImage(image.Rect(sx, sy, sx+sprite.frameWidth, sy+sprite.frameHeight)).(*ebiten.Image)

		screen.DrawImage(spriteSubImage, op)
	}
}

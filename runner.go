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

type RunnerActions uint16

const (
	RunnerIdle RunnerActions = 1 << iota
	RunnerRun
	RunnerShoot
	RunnerLeft
	RunnerRight
	RunnerUp
	RunnerDown
	RunnerHit
	RunnerDead
)

func (ra RunnerActions) Has(flags RunnerActions) bool {
	return ra&flags != 0
}

type BulletActions uint8

const (
	BulletLeft BulletActions = 1 << iota
	BulletRight
	BulletUp
	BulletDown
	BulletHit
)

func (ba BulletActions) Has(flags BulletActions) bool {
	return ba&flags != 0
}

type Bullet struct {
	sprite  Sprite
	actions BulletActions
	x, y    int
}

func NewBullet(x, y int, a BulletActions) *Bullet {
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
	if a.Has(BulletUp | BulletDown) {
		b.sprite.frameOX = 1
	}

	return b
}

type Runner struct {
	x, y       int
	vx, vy     int
	frameCount int
	actions    RunnerActions
	sprites    map[RunnerActions]Sprite
	bullets    []*Bullet
}

func NewRunner(x, y int) *Runner {
	var a RunnerActions = RunnerIdle | RunnerRight

	r := &Runner{
		x:       x,
		y:       y,
		actions: a,
	}

	r.sprites = map[RunnerActions]Sprite{
		RunnerLeft | RunnerIdle: {
			image:       runnerLeftImage,
			numFrames:   5,
			frameOX:     0,
			frameOY:     0,
			frameHeight: 32,
			frameWidth:  32,
		},
		RunnerLeft | RunnerRun: {
			image:       runnerLeftImage,
			numFrames:   8,
			frameOX:     0,
			frameOY:     32,
			frameHeight: 32,
			frameWidth:  32,
		},
		RunnerLeft | RunnerShoot: {
			image:       runnerLeftImage,
			numFrames:   4,
			frameOX:     0,
			frameOY:     64,
			frameHeight: 32,
			frameWidth:  32,
		},
		RunnerRight | RunnerIdle: {
			image:       runnerRightImage,
			numFrames:   5,
			frameOX:     0,
			frameOY:     0,
			frameHeight: 32,
			frameWidth:  32,
		},
		RunnerRight | RunnerRun: {
			image:       runnerRightImage,
			numFrames:   8,
			frameOX:     0,
			frameOY:     32,
			frameHeight: 32,
			frameWidth:  32,
		},
		RunnerRight | RunnerShoot: {
			image:       runnerRightImage,
			numFrames:   4,
			frameOX:     0,
			frameOY:     64,
			frameHeight: 32,
			frameWidth:  32,
		},
		RunnerUp | RunnerIdle: {
			image:       runnerRightImage,
			numFrames:   5,
			frameOX:     0,
			frameOY:     0,
			frameHeight: 32,
			frameWidth:  32,
		},
		RunnerUp | RunnerRun: {
			image:       runnerRightImage,
			numFrames:   8,
			frameOX:     0,
			frameOY:     32,
			frameHeight: 32,
			frameWidth:  32,
		},
		RunnerUp | RunnerShoot: {
			image:       runnerRightImage,
			numFrames:   4,
			frameOX:     0,
			frameOY:     64,
			frameHeight: 32,
			frameWidth:  32,
		},
		RunnerDown | RunnerIdle: {
			image:       runnerLeftImage,
			numFrames:   5,
			frameOX:     0,
			frameOY:     0,
			frameHeight: 32,
			frameWidth:  32,
		},
		RunnerDown | RunnerRun: {
			image:       runnerLeftImage,
			numFrames:   8,
			frameOX:     0,
			frameOY:     32,
			frameHeight: 32,
			frameWidth:  32,
		},
		RunnerDown | RunnerShoot: {
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
	r.actions = r.actions &^ (RunnerRun | RunnerShoot)
	r.actions = r.actions | RunnerIdle

	// RunnerUpdate runner state based on keyboard input
	// H - RunnerLeft
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		r.actions = RunnerLeft | RunnerRun
		r.vx -= moveBy
	}

	// L - RunnerRight
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		r.actions = RunnerRight | RunnerRun
		r.vx += moveBy
	}

	// K - RunnerUp
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		r.actions = RunnerUp | RunnerRun
		r.vy -= moveBy
	}

	// J - RunnerDown
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		r.actions = RunnerDown | RunnerRun
		r.vy += moveBy
	}

	// SPACE - Spacebar
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		if len(r.bullets) < maxBullets {
			direction := runnerDirToBulletDir(r)
			bullet := NewBullet(
				r.x,
				r.y,
				direction,
			)
			r.bullets = append(r.bullets, bullet)

		}
		r.actions = r.actions &^ (RunnerIdle | RunnerRun)
		r.actions = r.actions | RunnerShoot
	}

	// RunnerUpdate sprite's x & y positions based on velocity values and
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

	// RunnerUpdate bullets if any
	var activeBullets []*Bullet
	for _, bullet := range r.bullets {
		if bullet.actions.Has(BulletLeft) {
			bullet.x -= 5
		}
		if bullet.actions.Has(BulletRight) {
			bullet.x += 5
		}
		if bullet.actions.Has(BulletUp) {
			bullet.y -= 5
		}
		if bullet.actions.Has(BulletDown) {
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

func (r Runner) direction() RunnerActions {
	switch {
	case r.actions.Has(RunnerLeft):
		return RunnerLeft
	case r.actions.Has(RunnerRight):
		return RunnerRight
	case r.actions.Has(RunnerUp):
		return RunnerUp
	case r.actions.Has(RunnerDown):
		return RunnerDown
	default:
		return RunnerRight
	}
}

func (r Runner) action() RunnerActions {
	switch {
	case r.actions.Has(RunnerIdle):
		return RunnerIdle
	case r.actions.Has(RunnerRun):
		return RunnerRun
	case r.actions.Has(RunnerShoot):
		return RunnerShoot
	default:
		return RunnerIdle
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

func runnerDirToBulletDir(r *Runner) BulletActions {
	var direction BulletActions
	switch {
	case r.actions.Has(RunnerLeft):
		direction = BulletLeft
	case r.actions.Has(RunnerRight):
		direction = BulletRight
	case r.actions.Has(RunnerUp):
		direction = BulletUp
	case r.actions.Has(RunnerDown):
		direction = BulletDown
	}

	return direction
}

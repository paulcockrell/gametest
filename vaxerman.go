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
	vaxermanImage *ebiten.Image
	bulletImage   *ebiten.Image
)

const (
	maxBullets = 3
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(images.VaxerMan_png))
	if err != nil {
		log.Fatalf("error decoding image: %v", err)
	}
	vaxermanImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

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

type VaxerManActions uint16

const (
	VaxerManIdle VaxerManActions = 1 << iota
	VaxerManRun
	VaxerManShoot
	VaxerManLeft
	VaxerManRight
	VaxerManUp
	VaxerManDown
	VaxerManHit
	VaxerManDead
)

func (ra VaxerManActions) Has(flags VaxerManActions) bool {
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

type VaxerMan struct {
	x, y       int
	vx, vy     int
	frameCount int
	actions    VaxerManActions
	sprites    map[VaxerManActions]Sprite
	bullets    []*Bullet
}

func NewVaxerMan(x, y int) *VaxerMan {
	var a VaxerManActions = VaxerManIdle | VaxerManRight

	r := &VaxerMan{
		x:       x,
		y:       y,
		actions: a,
	}

	r.sprites = map[VaxerManActions]Sprite{
		VaxerManLeft | VaxerManIdle: {
			image:       vaxermanImage,
			numFrames:   4,
			frameOX:     32 * 0,
			frameOY:     32 * 0,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManLeft | VaxerManRun: {
			image:       vaxermanImage,
			numFrames:   6,
			frameOX:     32 * 0,
			frameOY:     32 * 4,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManLeft | VaxerManShoot: {
			image:       vaxermanImage,
			numFrames:   5,
			frameOX:     32 * 0,
			frameOY:     32 * 2,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManRight | VaxerManIdle: {
			image:       vaxermanImage,
			numFrames:   4,
			frameOX:     32 * 0,
			frameOY:     32 * 1,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManRight | VaxerManRun: {
			image:       vaxermanImage,
			numFrames:   6,
			frameOX:     32 * 0,
			frameOY:     32 * 5,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManRight | VaxerManShoot: {
			image:       vaxermanImage,
			numFrames:   5,
			frameOX:     32 * 0,
			frameOY:     32 * 3,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManUp | VaxerManIdle: {
			image:       vaxermanImage,
			numFrames:   4,
			frameOX:     32 * 0,
			frameOY:     32 * 1,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManUp | VaxerManRun: {
			image:       vaxermanImage,
			numFrames:   6,
			frameOX:     32 * 0,
			frameOY:     32 * 5,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManUp | VaxerManShoot: {
			image:       vaxermanImage,
			numFrames:   4,
			frameOX:     32 * 0,
			frameOY:     32 * 3,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManDown | VaxerManIdle: {
			image:       vaxermanImage,
			numFrames:   4,
			frameOX:     32 * 0,
			frameOY:     32 * 0,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManDown | VaxerManRun: {
			image:       vaxermanImage,
			numFrames:   6,
			frameOX:     32 * 0,
			frameOY:     32 * 4,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManDown | VaxerManShoot: {
			image:       vaxermanImage,
			numFrames:   5,
			frameOX:     32 * 0,
			frameOY:     32 * 2,
			frameHeight: 32,
			frameWidth:  32,
		},
	}

	return r
}

func (r *VaxerMan) update() {
	const moveBy = 2

	// Reset velocity values
	r.vx = 0
	r.vy = 0

	// Reset movement to default idle
	r.actions = r.actions &^ (VaxerManRun | VaxerManShoot)
	r.actions = r.actions | VaxerManIdle

	// VaxerManUpdate vaxerman state based on keyboard input
	// H - VaxerManLeft
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		r.actions = VaxerManLeft | VaxerManRun
		r.vx -= moveBy
	}

	// L - VaxerManRight
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		r.actions = VaxerManRight | VaxerManRun
		r.vx += moveBy
	}

	// K - VaxerManUp
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		r.actions = VaxerManUp | VaxerManRun
		r.vy -= moveBy
	}

	// J - VaxerManDown
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		r.actions = VaxerManDown | VaxerManRun
		r.vy += moveBy
	}

	// SPACE - Spacebar
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		if len(r.bullets) < maxBullets {
			direction := vaxermanDirToBulletDir(r)
			bullet := NewBullet(
				r.x,
				r.y,
				direction,
			)
			r.bullets = append(r.bullets, bullet)

		}
		r.actions = r.actions &^ (VaxerManIdle | VaxerManRun)
		r.actions = r.actions | VaxerManShoot
	}

	// VaxerManUpdate sprite's x & y positions based on velocity values and
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

	// VaxerManUpdate bullets if any
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

func (r VaxerMan) direction() VaxerManActions {
	switch {
	case r.actions.Has(VaxerManLeft):
		return VaxerManLeft
	case r.actions.Has(VaxerManRight):
		return VaxerManRight
	case r.actions.Has(VaxerManUp):
		return VaxerManUp
	case r.actions.Has(VaxerManDown):
		return VaxerManDown
	default:
		return VaxerManRight
	}
}

func (r VaxerMan) action() VaxerManActions {
	switch {
	case r.actions.Has(VaxerManIdle):
		return VaxerManIdle
	case r.actions.Has(VaxerManRun):
		return VaxerManRun
	case r.actions.Has(VaxerManShoot):
		return VaxerManShoot
	default:
		return VaxerManIdle
	}
}

func (r VaxerMan) getSprite() Sprite {
	direction := r.direction()
	action := r.action()

	return r.sprites[direction|action]
}

func (r *VaxerMan) draw(screen *ebiten.Image) {
	sprite := r.getSprite()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(r.x), float64(r.y))

	// Extract sprite frame
	i := (r.frameCount / sprite.numFrames) % sprite.numFrames
	sx, sy := sprite.frameOX+i*sprite.frameWidth, sprite.frameOY
	spriteSubImage := sprite.image.SubImage(image.Rect(sx, sy, sx+sprite.frameWidth, sy+sprite.frameHeight)).(*ebiten.Image)

	screen.DrawImage(spriteSubImage, op)
}

func (r *VaxerMan) drawBullets(screen *ebiten.Image) {
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

func vaxermanDirToBulletDir(r *VaxerMan) BulletActions {
	var direction BulletActions
	switch {
	case r.actions.Has(VaxerManLeft):
		direction = BulletLeft
	case r.actions.Has(VaxerManRight):
		direction = BulletRight
	case r.actions.Has(VaxerManUp):
		direction = BulletUp
	case r.actions.Has(VaxerManDown):
		direction = BulletDown
	}

	return direction
}

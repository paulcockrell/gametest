package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/paulcockrell/gametest/resources/images"
)

const (
	maxBullets = 3
)

var (
	vaxermanImage *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(images.VaxerMan_png))
	if err != nil {
		log.Fatalf("error decoding image: %v", err)
	}
	vaxermanImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
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

type VaxerMan struct {
	Health      int
	x, y        int
	vx, vy      int
	frameCount  int
	actions     VaxerManActions
	sprites     map[VaxerManActions]Sprite
	bullets     []*Bullet
	firingTimer int
}

func NewVaxerMan(x, y int) *VaxerMan {
	var a VaxerManActions = VaxerManIdle | VaxerManRight

	v := &VaxerMan{
		x:       x,
		y:       y,
		actions: a,
		Health:  100,
	}

	v.sprites = map[VaxerManActions]Sprite{
		VaxerManLeft | VaxerManIdle: {
			image:       vaxermanImage,
			numFrames:   6,
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
			numFrames:   6,
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
			numFrames:   6,
			frameOX:     32 * 0,
			frameOY:     32 * 9,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManUp | VaxerManRun: {
			image:       vaxermanImage,
			numFrames:   6,
			frameOX:     32 * 0,
			frameOY:     32 * 11,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManUp | VaxerManShoot: {
			image:       vaxermanImage,
			numFrames:   5,
			frameOX:     32 * 0,
			frameOY:     32 * 10,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManDown | VaxerManIdle: {
			image:       vaxermanImage,
			numFrames:   6,
			frameOX:     32 * 0,
			frameOY:     32 * 6,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManDown | VaxerManRun: {
			image:       vaxermanImage,
			numFrames:   6,
			frameOX:     32 * 0,
			frameOY:     32 * 8,
			frameHeight: 32,
			frameWidth:  32,
		},
		VaxerManDown | VaxerManShoot: {
			image:       vaxermanImage,
			numFrames:   5,
			frameOX:     32 * 0,
			frameOY:     32 * 7,
			frameHeight: 32,
			frameWidth:  32,
		},
	}

	return v
}

// Infect decrements the lives counter, if lives counter reaches zero it sets
// VaxerMan to dead
func (v *VaxerMan) Infect() {
	v.Health -= 10
	if v.Health == 0 {
		v.actions = VaxerManDead
	}
}

// IsDead returns true if vaxermans actions contains VaxerManDead
func (v *VaxerMan) IsDead() bool {
	return v.actions.Has(VaxerManDead)
}

func (v *VaxerMan) update() {
	const moveBy = 2

	// Update bullets
	var activeBullets []*Bullet
	for _, bullet := range v.bullets {
		bullet.Update()
		if bullet.IsLive() {
			activeBullets = append(activeBullets, bullet)
		}
	}
	v.bullets = activeBullets

	// If VaxerMan is dead, do nothing
	if v.IsDead() {
		return
	}

	// Reset velocity values
	v.vx = 0
	v.vy = 0

	// Reset movement to default idle
	v.actions = v.actions &^ (VaxerManRun | VaxerManShoot)
	v.actions = v.actions | VaxerManIdle

	// Respond to keyboard inputs
	// VaxerManLeft
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		v.actions = VaxerManLeft | VaxerManRun
		v.vx -= moveBy
	}

	// VaxerManRight
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		v.actions = VaxerManRight | VaxerManRun
		v.vx += moveBy
	}

	// VaxerManUp
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		v.actions = VaxerManUp | VaxerManRun
		v.vy -= moveBy
	}

	// VaxerManDown
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		v.actions = VaxerManDown | VaxerManRun
		v.vy += moveBy
	}

	// SPACE - Spacebar
	if ebiten.IsKeyPressed(ebiten.KeySpace) && v.canFire() {
		if len(v.bullets) < maxBullets {
			bullet := NewBullet(
				0,
				0,
				vaxermanDirToBulletDir(v),
			)

			bx, by := v.x, v.y
			if v.actions.Has(VaxerManLeft) {
				bx -= v.GetSprite().frameWidth / 2
				by += (v.GetSprite().frameWidth / 2) - (bullet.sprite.frameHeight / 2)
			}
			if v.actions.Has(VaxerManRight) {
				bx += v.GetSprite().frameWidth
				by += (v.GetSprite().frameWidth / 2) - (bullet.sprite.frameHeight / 2)
			}
			if v.actions.Has(VaxerManUp) {
				by -= v.GetSprite().frameHeight / 2
				bx += (v.GetSprite().frameHeight / 2) - (bullet.sprite.frameWidth / 2)
			}
			if v.actions.Has(VaxerManDown) {
				by += v.GetSprite().frameHeight
				bx += (v.GetSprite().frameHeight / 2) - (bullet.sprite.frameWidth / 2)
			}

			bullet.x = bx
			bullet.y = by

			v.firingTimer = 5
			v.bullets = append(v.bullets, bullet)
		}
		v.actions = v.actions &^ (VaxerManIdle | VaxerManRun)
		v.actions = v.actions | VaxerManShoot
	}

	// VaxerManUpdate sprite's x & y positions based on velocity values and
	// frame counter used by animation
	v.frameCount++
	v.x += v.vx
	v.y += v.vy

	if v.firingTimer > 0 {
		v.firingTimer -= 1
	}

	// Wall collision detection
	s := v.GetSprite()
	if v.x < 0 {
		v.x = 0
	}
	if v.x > screenWidth-s.frameWidth {
		v.x = screenWidth - s.frameWidth
	}
	if v.y < 0 {
		v.y = 0
	}
	if v.y > screenHeight-s.frameHeight {
		v.y = screenHeight - s.frameHeight
	}
}

func (v VaxerMan) direction() VaxerManActions {
	switch {
	case v.actions.Has(VaxerManLeft):
		return VaxerManLeft
	case v.actions.Has(VaxerManRight):
		return VaxerManRight
	case v.actions.Has(VaxerManUp):
		return VaxerManUp
	case v.actions.Has(VaxerManDown):
		return VaxerManDown
	default:
		return VaxerManRight
	}
}

func (v VaxerMan) action() VaxerManActions {
	switch {
	case v.actions.Has(VaxerManIdle):
		return VaxerManIdle
	case v.actions.Has(VaxerManRun):
		return VaxerManRun
	case v.actions.Has(VaxerManShoot):
		return VaxerManShoot
	default:
		return VaxerManIdle
	}
}

func (v VaxerMan) GetSprite() Sprite {
	direction := v.direction()
	action := v.action()

	return v.sprites[direction|action]
}

func (v *VaxerMan) draw(screen *ebiten.Image) {
	sprite := v.GetSprite()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(v.x), float64(v.y))

	// Extract sprite frame
	i := (v.frameCount / sprite.numFrames) % sprite.numFrames
	sx, sy := sprite.frameOX+i*sprite.frameWidth, sprite.frameOY
	spriteSubImage := sprite.image.SubImage(image.Rect(sx, sy, sx+sprite.frameWidth, sy+sprite.frameHeight)).(*ebiten.Image)

	screen.DrawImage(spriteSubImage, op)
}

func (v *VaxerMan) drawBullets(screen *ebiten.Image) {
	for _, bullet := range v.bullets {
		bullet.draw(screen)
	}
}

func (v *VaxerMan) hasShotEnemy(e *Enemy) bool {
	if e.status != EnemyAlive {
		return false
	}

	for _, bullet := range v.bullets {
		if bullet.HasHitEnemy(e) {
			e.SetNotInfectious()
			bullet.SetHit()
			boomPlayer.Rewind()
			boomPlayer.Play()
			return true
		}
	}

	return false
}

func (v VaxerMan) canFire() bool {
	return v.firingTimer == 0
}

func vaxermanDirToBulletDir(v *VaxerMan) BulletActions {
	var direction BulletActions
	switch {
	case v.actions.Has(VaxerManLeft):
		direction = BulletLeft
	case v.actions.Has(VaxerManRight):
		direction = BulletRight
	case v.actions.Has(VaxerManUp):
		direction = BulletUp
	case v.actions.Has(VaxerManDown):
		direction = BulletDown
	}

	return direction
}

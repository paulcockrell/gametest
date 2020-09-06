package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/paulcockrell/gametest/resources/images"
)

var (
	bulletImage *ebiten.Image
)

const (
	maxBullets = 3
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(images.Bullet_png))
	if err != nil {
		log.Fatalf("error decoding image: %v", err)
	}
	bulletImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
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
	sprite     Sprite
	actions    BulletActions
	x, y       int
	frameCount int
}

// NewBullet constructs a bullet sprite at the given position and direction
func NewBullet(x, y int, a BulletActions) *Bullet {
	b := &Bullet{
		x:       x, // starting x
		y:       y, // starting y
		actions: a, // holds direction
	}
	b.sprite = Sprite{
		image:       bulletImage,
		numFrames:   2,
		frameOX:     0,
		frameOY:     0,
		frameHeight: 32,
		frameWidth:  32,
	}

	return b
}

// SetHit sets action to Bullet
func (b *Bullet) SetHit() {
	b.actions = BulletHit
}

func (b *Bullet) draw(screen *ebiten.Image) {
	if b.actions.Has(BulletHit) {
		return
	}

	sprite := b.sprite

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))

	// Extract sprite frame
	i := (b.frameCount / sprite.numFrames) % sprite.numFrames
	sx, sy := sprite.frameOX+i*sprite.frameWidth, sprite.frameOY
	spriteSubImage := sprite.image.SubImage(image.Rect(sx, sy, sx+sprite.frameWidth, sy+sprite.frameHeight)).(*ebiten.Image)

	screen.DrawImage(spriteSubImage, op)
}

// Update updates the bullets location
func (b *Bullet) Update() {
	if b.actions.Has(BulletHit) {
		return
	}

	if b.actions.Has(BulletLeft) {
		b.x -= 5
	}
	if b.actions.Has(BulletRight) {
		b.x += 5
	}
	if b.actions.Has(BulletUp) {
		b.y -= 5
	}
	if b.actions.Has(BulletDown) {
		b.y += 5
	}
	b.frameCount++
}

// IsLive checks if bullet is on screen and state doesn't include BulletHit
func (b *Bullet) IsLive() bool {
	if b.actions.Has(BulletHit) {
		return false
	}

	// Is bullet on screen
	if b.x < 0 || b.x > screenWidth ||
		b.y < 0 || b.y > screenHeight {
		return false
	}

	return true
}

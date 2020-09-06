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
		numFrames:   2,
		frameOX:     0,
		frameOY:     0,
		frameHeight: 32,
		frameWidth:  32,
	}

	return b
}

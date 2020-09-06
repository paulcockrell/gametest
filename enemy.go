package main

import (
	"bytes"
	"image"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/paulcockrell/gametest/resources/images"
)

var (
	enemyImage *ebiten.Image
)

const (
	MaxEnemies = 5
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(images.Enemy_png))
	if err != nil {
		log.Fatalf("error decoding image: %v", err)
	}
	enemyImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

type EnemyActions uint8

const (
	EnemyAlive EnemyActions = iota
	EnemyHit
	EnemyDead
)

type Enemy struct {
	x, y           int
	vx, vy         int
	frameCount     int
	animFrameCount int // used to make sure we play a whole anim sequence at least once
	status         EnemyActions
	sprites        map[EnemyActions]Sprite
}

func NewEnemy(x, y, vx, vy int) *Enemy {
	e := &Enemy{
		x:      x,
		y:      y,
		vx:     vx,
		vy:     vy,
		status: EnemyAlive,
	}

	e.sprites = map[EnemyActions]Sprite{
		EnemyAlive: {
			image:       enemyImage,
			numFrames:   2,
			frameOX:     32 * 0,
			frameOY:     32 * 0,
			frameHeight: 32,
			frameWidth:  32,
		},
		EnemyHit: {
			image:       enemyImage,
			numFrames:   6,
			frameOX:     32 * 0,
			frameOY:     32 * 1,
			frameHeight: 32,
			frameWidth:  32,
		},
		EnemyDead: {
			image:       enemyImage,
			numFrames:   1,
			frameOX:     32 * 0,
			frameOY:     32 * 2,
			frameHeight: 32,
			frameWidth:  32,
		},
	}

	return e
}

func (e *Enemy) update() {
	e.frameCount++
	e.x += e.vx
	e.y += e.vy
	if e.x < 0 || e.y < 0 || e.x > screenWidth || e.y > screenHeight {
		e.status = EnemyDead
	}
}

func (e *Enemy) draw(screen *ebiten.Image) {
	sprite := e.GetSprite()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(e.x), float64(e.y))

	// Extract sprite frame
	i := (e.frameCount / sprite.numFrames) % sprite.numFrames
	sx, sy := sprite.frameOX+i*sprite.frameWidth, sprite.frameOY
	spriteSubImage := sprite.image.SubImage(image.Rect(sx, sy, sx+sprite.frameWidth, sy+sprite.frameHeight)).(*ebiten.Image)

	screen.DrawImage(spriteSubImage, op)
}

func (e *Enemy) GetSprite() (sprite Sprite) {
	sprite = e.sprites[e.status]
	return
}

func GenerateEnemyStartPos() (x, y, vx, vy int) {
	coinFlipOne := rand.Intn(1)
	coinFlipTwo := rand.Intn(1)
	if coinFlipOne == 1 {
		x = rand.Intn(screenWidth)
		if coinFlipTwo == 1 {
			y = 0
		} else {
			y = screenHeight
		}
	} else {
		y = rand.Intn(screenHeight)
		if coinFlipTwo == 1 {
			x = 0
		} else {
			x = screenWidth
		}
	}
	vx = rand.Intn(3-1) + 1
	vy = rand.Intn(3-1) + 1
	if x > screenWidth/2 {
		vx *= -1
	}
	if y > screenHeight/2 {
		vy *= -1
	}

	return
}
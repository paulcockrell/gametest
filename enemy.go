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
	// MaxEnemies sets the limit of enemies that can be 'alive' at any one time
	MaxEnemies = 3
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(images.Enemy_png))
	if err != nil {
		log.Fatalf("error decoding image: %v", err)
	}
	enemyImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

// EnemyActions is an integer type that holds the various Enemy actions
type EnemyActions uint8

const (
	EnemyAlive EnemyActions = iota
	EnemyHit
	EnemyDead
)

// Enemy defines an enemy
type Enemy struct {
	x, y          int
	vx, vy        int
	frameCount    int
	hitFrameCount int // used to make sure we play a whole hit anim sequence at least once
	status        EnemyActions
	isInfectious  bool
	sprites       map[EnemyActions]Sprite
}

// NewEnemy builds an enemy at the given position and velocity
func NewEnemy(x, y, vx, vy int) *Enemy {
	e := &Enemy{
		x:            x,
		y:            y,
		vx:           vx,
		vy:           vy,
		status:       EnemyAlive,
		isInfectious: true,
	}

	e.sprites = map[EnemyActions]Sprite{
		EnemyAlive: {
			image:       enemyImage,
			numFrames:   1,
			frameOX:     32 * 0,
			frameOY:     32 * 0,
			frameHeight: 32,
			frameWidth:  32,
		},
		EnemyHit: {
			image:       enemyImage,
			numFrames:   4,
			frameOX:     32 * 0,
			frameOY:     32 * 1,
			frameHeight: 32,
			frameWidth:  32,
		},
		EnemyDead: {
			image:       enemyImage,
			numFrames:   1,
			frameOX:     32 * 1,
			frameOY:     32 * 0,
			frameHeight: 32,
			frameWidth:  32,
		},
	}

	return e
}

func (e *Enemy) update() {
	if e.status == EnemyHit {
		e.hitFrameCount++

		if e.hitFrameCount > e.GetSprite().numFrames {
			e.status = EnemyDead
			return
		}
	}

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

// GetSprite returns the current sprite by status
func (e *Enemy) GetSprite() (sprite Sprite) {
	sprite = e.sprites[e.status]
	return
}

// IsDead return true if status is EnemyDead
func (e Enemy) IsDead() bool {
	return e.status == EnemyDead
}

// HasInfectedPlayer returns bool based on collision between enemy and player
func (e *Enemy) HasInfectedPlayer(v *VaxerMan) bool {
	// You can only infect VaxerMan once
	if !e.isInfectious {
		return false
	}
	if v.actions.Has(VaxerManDead) {
		return false
	}

	eSprite := e.GetSprite()
	vSprite := v.GetSprite()
	if e.x >= v.x+vSprite.frameWidth || v.x >= e.x+eSprite.frameWidth {
		return false
	}
	if e.y >= v.y+vSprite.frameHeight || v.y >= e.y+eSprite.frameHeight {
		return false
	}

	e.isInfectious = false

	return true
}

// GenerateEnemyStartPos randomly generates position and velocity values
func GenerateEnemyStartPos() (x, y, vx, vy int) {
	coinFlipOne := rand.Intn(2)
	coinFlipTwo := rand.Intn(2)
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

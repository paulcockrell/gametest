package main

import (
	"bytes"
	"encoding/json"
	"image"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/paulcockrell/gametest/resources/images"
)

// Tile constants
const (
	tileSize = 16
	tileXNum = 25
)

var (
	tilesImage *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatalf("error loading tiles images: %v", err)
	}
	tilesImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

type Level struct {
	layers [][]int
}

func NewLevel(path string) *Level {
	level := &Level{}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error loading level: %q", err)
	}

	err = json.Unmarshal([]byte(file), &level.layers)
	if err != nil {
		log.Fatalf("error parsing level: %q", err)
	}

	return level
}

func (l *Level) draw(screen *ebiten.Image) {
	const xNum = screenWidth / tileSize
	for _, l := range l.layers {
		for i, t := range l {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64((i%xNum)*tileSize), float64((i/xNum)*tileSize))

			sx := (t % tileXNum) * tileSize
			sy := (t / tileXNum) * tileSize
			tile := tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image)
			screen.DrawImage(tile, op)
		}
	}
}

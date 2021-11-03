package blit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"

	"github.com/octalide/blit/pkg/bgl"
)

type Spritesheet struct {
	Texture *bgl.Texture
	Sprites map[string][4]float64
	Filter  bgl.Filter
}

func NewSpritesheet() *Spritesheet {
	ss := &Spritesheet{
		Sprites: map[string][4]float64{},
		Filter:  bgl.Nearest,
	}

	return ss
}

func LoadSpritesheet(name string) (*Spritesheet, error) {
	raw, err := ioutil.ReadFile(fmt.Sprintf("%v.json", name))
	if err != nil {
		return nil, err
	}

	ss := NewSpritesheet()

	err = json.Unmarshal(raw, &ss.Sprites)
	if err != nil {
		return nil, err
	}

	raw, err = ioutil.ReadFile(fmt.Sprintf("%v.png", name))
	if err != nil {
		return nil, err
	}

	src, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	b := src.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)

	ss.Texture = bgl.NewTexture(img, ss.Filter)

	return ss, nil
}

func (ss *Spritesheet) Get(name string, shader *bgl.Program) (*Sprite, error) {
	raw := ss.Sprites[name]

	x := float32(raw[0])
	y := float32(raw[1])
	w := float32(raw[2])
	h := float32(raw[3])

	return NewSprite(shader, ss.Texture, x, y, w, h)
}

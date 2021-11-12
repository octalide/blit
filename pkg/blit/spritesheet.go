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
	Sprites map[string]Rect
	Filter  bgl.Filter
}

func NewSpritesheet() *Spritesheet {
	ss := &Spritesheet{
		Sprites: map[string]Rect{},
		Filter:  bgl.Nearest,
	}

	return ss
}

func GenSpritesheet(img *image.RGBA, desc map[string]interface{}) (*Spritesheet, error) {
	ss := NewSpritesheet()

	for k, v := range desc {
		raw := v.([]interface{})

		x, ok := raw[0].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid x value for \"%v\"", k)
		}
		y, ok := raw[1].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid y value for \"%v\"", k)
		}
		w, ok := raw[2].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid w value for \"%v\"", k)
		}
		h, ok := raw[3].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid h value for \"%v\"", k)
		}

		ss.Sprites[k] = Rect{
			float32(x),
			float32(y),
			float32(w),
			float32(h),
		}
	}

	b := img.Bounds()
	img = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), img, b.Min, draw.Src)

	ss.Texture = bgl.NewTexture(img, ss.Filter)

	return ss, nil
}

func LoadSpritesheet(imgPath, descPath string) (*Spritesheet, error) {
	raw, err := ioutil.ReadFile(imgPath)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	desc, err := ioutil.ReadFile(descPath)
	if err != nil {
		return nil, err
	}

	var descMap map[string]interface{}
	err = json.Unmarshal(desc, &descMap)
	if err != nil {
		return nil, err
	}

	// convert to RGBA
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return GenSpritesheet(rgba, descMap)
}

func (ss *Spritesheet) Get(name string, shader *bgl.Program) (*Sprite, error) {
	rect, ok := ss.Sprites[name]
	if !ok {
		return nil, fmt.Errorf("sprite name not found: \"%v\"", name)
	}

	return NewSprite(shader, ss.Texture, rect)
}

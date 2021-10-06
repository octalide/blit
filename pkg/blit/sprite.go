package blit

import "image/color"

type Sprite struct {
	frame Rect
	mat   Mat
	mask  color.RGBA
}

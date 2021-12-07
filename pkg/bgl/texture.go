package bgl

import (
	"image"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Filter int32

const (
	TextureBinding2D uint32 = gl.TEXTURE_BINDING_2D
	Linear                  = Filter(gl.LINEAR)
	Nearest                 = Filter(gl.NEAREST)
)

// Texture is an OpenGL texture.
type Texture struct {
	ID            uint32
	width, height int
	filter        int32
}

// NewTexture creates a new texture with the specified width and height with some initial
// pixel values. The pixels must be a sequence of RGBA values (one byte per component).
func NewTexture(img *image.RGBA, filter Filter) *Texture {
	width := img.Rect.Max.X - img.Rect.Min.X
	height := img.Rect.Max.Y - img.Rect.Min.Y

	tex := &Texture{
		width:  width,
		height: height,
	}

	gl.GenTextures(1, &tex.ID)

	tex.Bind()

	// initial data
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(width),
		int32(height),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(img.Pix),
	)

	/*
		borderColor := mgl32.Vec4{0, 0, 0, 0}
		gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
	*/

	tex.SetFilter(filter)

	tex.Unbind()

	runtime.SetFinalizer(tex, (*Texture).Delete)

	return tex
}

// Delete deletes the Texture.
func (t *Texture) Delete() {
	gl.DeleteTextures(1, &t.ID)
}

// Width returns the width of the Texture in pixels.
func (t *Texture) Width() int {
	return t.width
}

// Height returns the height of the Texture in pixels.
func (t *Texture) Height() int {
	return t.height
}

// SetPixels sets the content of a sub-region of the Texture
func (t *Texture) SetPixels(x, y, w, h int, pix []uint32) {
	if len(pix) != w*h*4 {
		panic("set pixels: wrong number of pixels")
	}

	gl.TexSubImage2D(
		gl.TEXTURE_2D,
		0,
		int32(x),
		int32(y),
		int32(w),
		int32(h),
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(pix),
	)
}

// Pixels returns the content of a sub-region of the Texture as an RGBA byte sequence.
func (t *Texture) Pixels(x, y, w, h int) []uint8 {
	pixels := make([]uint8, t.width*t.height*4)

	gl.GetTexImage(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(pixels),
	)

	subPixels := make([]uint8, w*h*4)
	for i := 0; i < h; i++ {
		row := pixels[(i+y)*t.width*4+x*4 : (i+y)*t.width*4+(x+w)*4]
		subRow := subPixels[i*w*4 : (i+1)*w*4]
		copy(subRow, row)
	}

	return subPixels
}

// SetFilter sets the filter of the Texture.
func (t *Texture) SetFilter(filter Filter) {
	t.filter = int32(filter)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, t.filter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, t.filter)
}

// Bind binds the Texture
func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}

// Unbind unbinds the Texture
func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// UV returns the uv coordinates of the Texture (utility function)
func (t *Texture) UV(x, y float32) [2]float32 {
	return [2]float32{
		x / float32(t.width),
		y / float32(t.height),
	}
}

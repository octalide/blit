package bgl

import (
	"image"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Filter int32

const (
	TextureBinding2D uint32 = gl.TEXTURE_BINDING_2D
	Linear                  = Filter(gl.LINEAR)
	Nearest                 = Filter(gl.NEAREST)
)

// Texture is an OpenGL texture.
type Texture struct {
	binder
	width, height int
	filter        int32
}

// NewTexture creates a new texture with the specified width and height with some initial
// pixel values. The pixels must be a sequence of RGBA values (one byte per component).
func NewTexture(width, height int, filter Filter, img *image.RGBA) *Texture {
	tex := &Texture{
		binder: binder{
			binding: TextureBinding2D,
			bindFunc: func(id uint32) {
				gl.BindTexture(gl.TEXTURE_2D, id)
			},
		},
		width:  width,
		height: height,
	}

	gl.GenTextures(1, &tex.id)

	tex.Begin()
	defer tex.End()

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

	borderColor := mgl32.Vec4{0, 0, 0, 0}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)

	tex.SetFilter(filter)

	runtime.SetFinalizer(tex, (*Texture).delete)

	return tex
}

func (t *Texture) delete() {
	gl.DeleteTextures(1, &t.id)
}

// ID returns the OpenGL ID of this Texture.
func (t *Texture) ID() uint32 {
	return t.id
}

// Width returns the width of the Texture in pixels.
func (t *Texture) Width() int {
	return t.width
}

// Height returns the height of the Texture in pixels.
func (t *Texture) Height() int {
	return t.height
}

// SetPixels sets the content of a sub-region of the Texture. Pixels must be an RGBA byte sequence.
func (t *Texture) SetPixels(x, y, w, h int, img image.RGBA) {
	if len(img.Pix) != w*h*4 {
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
		gl.Ptr(img.Pix),
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

// SetSmooth sets whether the Texture should be drawn "smoothly" or "pixely".
//
// It affects how the Texture is drawn when zoomed. Smooth interpolates between the neighbour
// pixels, while pixely always chooses the nearest pixel.
func (t *Texture) SetFilter(filter Filter) {
	t.filter = int32(filter)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, t.filter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, t.filter)
}

// Begin binds the Texture. This is necessary before using the Texture.
func (t *Texture) Begin() {
	t.bind()
}

// End unbinds the Texture and restores the previous one.
func (t *Texture) End() {
	t.restore()
}

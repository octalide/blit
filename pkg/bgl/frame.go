package bgl

import (
	"image"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	FramebufferBinding     uint32 = gl.FRAMEBUFFER_BINDING
	ReadFramebufferBinding uint32 = gl.READ_FRAMEBUFFER_BINDING
	DrawFramebufferBinding uint32 = gl.DRAW_FRAMEBUFFER_BINDING
)

type Frame struct {
	fb, rf, df binder
	tex        *Texture
}

// NewFrame creates a new fully transparent Frame with given dimensions in pixels.
func NewFrame(width, height int, filter Filter) *Frame {
	f := &Frame{
		fb: binder{
			binding: FramebufferBinding,
			bindFunc: func(obj uint32) {
				gl.BindFramebuffer(gl.FRAMEBUFFER, obj)
			},
		},
		rf: binder{
			binding: ReadFramebufferBinding,
			bindFunc: func(obj uint32) {
				gl.BindFramebuffer(gl.READ_FRAMEBUFFER, obj)
			},
		},
		df: binder{
			binding: DrawFramebufferBinding,
			bindFunc: func(obj uint32) {
				gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, obj)
			},
		},
		tex: NewTexture(width, height, filter, image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})),
	}

	gl.GenFramebuffers(1, &f.fb.id)

	f.fb.bind()
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, f.tex.id, 0)
	f.fb.restore()

	runtime.SetFinalizer(f, (*Frame).delete)

	return f
}

func (f *Frame) delete() {
	gl.DeleteFramebuffers(1, &f.fb.id)
}

// ID returns the OpenGL framebuffer ID of this Frame.
func (f *Frame) ID() uint32 {
	return f.fb.id
}

// Begin binds the Frame. All draw operations will target this Frame until End is called.
func (f *Frame) Begin() {
	f.fb.bind()
}

// End unbinds the Frame. All draw operations will go to whatever was bound before this Frame.
func (f *Frame) End() {
	f.fb.restore()
}

// Blit copies rectangle (sx0, sy0, sx1, sy1) in this Frame onto rectangle (dx0, dy0, dx1, dy1) in
// dst Frame.
//
// If the dst Frame is nil, the destination will be the framebuffer 0, which is the screen.
//
// If the sizes of the rectangles don't match, the source will be stretched to fit the destination
// rectangle. The stretch will be either smooth or pixely according to the source Frame's
// smoothness.
func (f *Frame) Blit(dst *Frame, sx0, sy0, sx1, sy1, dx0, dy0, dx1, dy1 int) {
	f.rf.id = f.fb.id
	if dst != nil {
		f.df.id = dst.fb.id
	} else {
		f.df.id = 0
	}
	f.rf.bind()
	f.df.bind()

	gl.BlitFramebuffer(
		int32(sx0), int32(sy0), int32(sx1), int32(sy1),
		int32(dx0), int32(dy0), int32(dx1), int32(dy1),
		gl.COLOR_BUFFER_BIT, uint32(f.tex.filter),
	)

	f.rf.restore()
	f.df.restore()
}

// Texture returns the Frame's underlying Texture that the Frame draws on.
func (f *Frame) Texture() *Texture {
	return f.tex
}

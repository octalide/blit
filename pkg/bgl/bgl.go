package bgl

import (
	"fmt"
	"image/color"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Init initializes OpenGL by loading function pointers from the active OpenGL context.
// This function must be manually run inside the main thread (using "github.com/faiface/mainthread"
// package).
//
// It must be called under the presence of an active OpenGL context, e.g., always after calling
// window.MakeContextCurrent(). Also, always call this function when switching contexts.
func Init() error {
	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to initialize opengl: %v", err)
	}

	// gl.Enable(gl.BLEND)
	// gl.BlendEquation(gl.FUNC_ADD)
	// gl.Enable(gl.SCISSOR_TEST)
	// gl.Enable(gl.DEPTH_TEST)

	// gl.Enable(gl.CULL_FACE)
	// gl.FrontFace(gl.CCW)
	// gl.CullFace(gl.BACK)

	return nil
}

func SetClearColor(c color.RGBA) {
	gl.ClearColor(
		float32(c.R)/255.0,
		float32(c.G)/255.0,
		float32(c.B)/255.0,
		float32(c.A)/255.0,
	)
}

func EnableMSAA() {
	gl.Enable(gl.MULTISAMPLE)
}

func DisableMSAA() {
	gl.Disable(gl.MULTISAMPLE)
}

func SetBounds(x, y, w, h int) {
	gl.Viewport(int32(x), int32(y), int32(w), int32(h))
	gl.Scissor(int32(x), int32(y), int32(w), int32(h))
}

// Viewport returns the viewport
func Viewport() [4]float32 {
	var vp [4]float32
	gl.GetFloatv(gl.VIEWPORT, &vp[0])
	return vp
}

// Aspect returns the aspect ratio of the current viewport.
func Aspect() float32 {
	vp := Viewport()
	return vp[2] / vp[3]
}

func Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// BlendFactor represents a source or destination blend factor.
type BlendFactor int

// Here's the list of all blend factors.
const (
	One              = BlendFactor(gl.ONE)
	Zero             = BlendFactor(gl.ZERO)
	SrcAlpha         = BlendFactor(gl.SRC_ALPHA)
	DstAlpha         = BlendFactor(gl.DST_ALPHA)
	OneMinusSrcAlpha = BlendFactor(gl.ONE_MINUS_SRC_ALPHA)
	OneMinusDstAlpha = BlendFactor(gl.ONE_MINUS_DST_ALPHA)
)

// BlendFunc sets the source and destination blend factor.
func BlendFunc(src, dst BlendFactor) {
	gl.BlendFunc(uint32(src), uint32(dst))
}

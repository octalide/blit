package blit

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/octalide/blit/pkg/bgl"
	"github.com/octalide/wisp/pkg/wisp"
)

const (
	OpenGLVersionMajor = 4
	OpenGLVersionMinor = 6
	OpenGLProfile      = glfw.OpenGLCoreProfile
)

func Init() error {
	wisp.Init()

	if err := bgl.Init(); err != nil {
		return fmt.Errorf("failed to initialize bgl: %v", err)
	}

	return nil
}

// Viewport returns the current viewport size.
func Viewport() Rect {
	vp := bgl.Viewport()

	return Rect{
		vp[0],
		vp[1],
		vp[2],
		vp[3],
	}
}

func Update() {
	glfw.GetCurrentContext().SwapBuffers()
	glfw.PollEvents()
}

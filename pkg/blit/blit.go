package blit

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/octalide/blit/pkg/bgl"
)

const (
	OpenGLVersionMajor = 4
	OpenGLVersionMinor = 6
	OpenGLProfile      = glfw.OpenGLCoreProfile
)

func Init() error {
	if err := bgl.Init(); err != nil {
		return fmt.Errorf("failed to initialize bgl: %v", err)
	}

	return nil
}

func Update() {
	glfw.GetCurrentContext().SwapBuffers()
	glfw.PollEvents()
}

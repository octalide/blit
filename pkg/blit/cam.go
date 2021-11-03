package blit

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/octalide/blit/pkg/bgl"
)

// Cam is a camera.
type Cam struct {
	*Orienter
	FOV      float32
	Viewport Rect
}

func NewCam() Cam {
	return Cam{
		FOV: 45,
		Orienter: &Orienter{
			Vec: Vec{0, 0, 1},
		},
	}
}

// Proj returns the projection matrix.
func (c Cam) Proj() Mat {
	return Perspective(mgl32.DegToRad(c.FOV), bgl.Aspect(), float32(0.1), float32(100))
}

// View returns the view matrix.
func (c Cam) View() Mat {
	eye := c.Vec

	center := c.Vec
	center[2] = 0

	up := Vec{0, 1, 0}

	return LookAt(eye, center, up)
}

// Use sets the matrices in the given shader using the uniforms "proj" and "view"
func (c Cam) Use(s *bgl.Program) {
	s.Bind()
	s.SetUniform("view", c.View().F())
	s.SetUniform("proj", c.Proj().F())
	s.Unbind()
}

// Pan moves the camera by the given vector.
func (c *Cam) Pan(v Vec) {
	c.Vec = c.Vec.Add(v)
}

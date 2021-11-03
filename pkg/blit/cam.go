package blit

import "github.com/octalide/blit/pkg/bgl"

// Cam is a camera.
type Cam struct {
	*Orienter
	Z float32
}

func NewCam() Cam {
	return Cam{
		Orienter: &Orienter{},
		Z:        1,
	}
}

// Proj returns the projection matrix.
func (c Cam) Proj() Mat {
	return Proj(c.Z)
}

// View returns the view matrix.
func (c Cam) View() Mat {
	return LookAt(Vec{}, c.Pos(), Vec{0, 1, 0})
}

// Use sets the matrices in the given shader using the uniforms "proj" and "view"
func (c Cam) Use(s *bgl.Program) {
	s.Bind()
	s.SetUniform("view", c.View().F())
	s.SetUniform("proj", c.Proj().F())
	s.Unbind()
}

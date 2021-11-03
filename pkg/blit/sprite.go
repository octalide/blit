package blit

import (
	"image/color"

	"github.com/octalide/blit/pkg/bgl"
)

var (
	quadDefault = []float32{
		// x, y, u, v
		-0.5, +0.5, 0, 0, // top left
		-0.5, -0.5, 0, 0, // bottom left
		+0.5, -0.5, 0, 0, // bottom right
		+0.5, +0.5, 0, 0, // top right
		-0.5, +0.5, 0, 0, // top left
		+0.5, -0.5, 0, 0, // bottom right
	}
)

type Sprite struct {
	O *Orienter

	shader *bgl.Program
	tex    *bgl.Texture

	vao *bgl.VAO
	vbo *bgl.VBO

	Mask    color.RGBA // Color mask
	Visible bool       // Visibility

	// texture rectangle
	Rect

	dirty bool
}

func NewSprite(shader *bgl.Program, texture *bgl.Texture, rect Rect) (*Sprite, error) {
	s := &Sprite{
		shader:  shader,
		tex:     texture,
		Rect:    rect,
		O:       &Orienter{},
		Visible: true,
	}

	// NOTE: VBO must be created and populated before VAO
	s.vbo = bgl.NewVBO()
	s.vbo.SetData(s.quad())
	s.vao = bgl.NewVAO(shader.VertexAttrs)

	return s, nil
}

func (s *Sprite) Dirty() {
	s.dirty = true
}

// Mat is a wrapper for s.O.Mat()
func (s *Sprite) Mat() Mat {
	return s.O.Mat()
}

func (s *Sprite) quad() []float32 {
	q := quadDefault

	uv := s.tex.UV
	// TL
	q[2] = uv(s.X(), s.Y())[0]
	q[3] = uv(s.X(), s.Y())[1]

	// BL
	q[6] = uv(s.X(), s.Y()+s.H())[0]
	q[7] = uv(s.X(), s.Y()+s.H())[1]

	// BR
	q[10] = uv(s.X()+s.W(), s.Y()+s.H())[0]
	q[11] = uv(s.X()+s.W(), s.Y()+s.H())[1]

	// TR
	q[14] = uv(s.X()+s.W(), s.Y())[0]
	q[15] = uv(s.X()+s.W(), s.Y())[1]

	// TL
	q[18] = uv(s.X(), s.Y())[0]
	q[19] = uv(s.X(), s.Y())[1]

	// BR
	q[22] = uv(s.X()+s.W(), s.Y()+s.H())[0]
	q[23] = uv(s.X()+s.W(), s.Y()+s.H())[1]

	return q
}

func (s *Sprite) mask() []float32 {
	return []float32{
		float32(s.Mask.R) / 255.0,
		float32(s.Mask.G) / 255.0,
		float32(s.Mask.B) / 255.0,
		float32(s.Mask.A) / 255.0,
	}
}

func (s *Sprite) Draw() {
	if s.Visible {
		s.shader.Bind()

		s.shader.SetUniform("color", s.mask())
		s.shader.SetUniform("modl", s.Mat().F())

		s.tex.Bind()
		s.vao.Bind()
		s.vbo.Bind()
		s.vbo.Draw()
		s.vbo.Unbind()
		s.vao.Unbind()
		s.tex.Unbind()
		s.shader.Unbind()
	}
}

type Animation struct {
	Frame  int
	Frames []Sprite
}

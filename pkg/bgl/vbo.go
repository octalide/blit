package bgl

import (
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// VBO is a vertex buffer object
type VBO struct {
	DrawMode
	Usage
	ID uint32

	size int
}

// NewVBO creates a new VBO
func NewVBO() *VBO {
	vbo := &VBO{
		DrawMode: Triangles,
		Usage:    DynamicDraw,
	}

	gl.GenBuffers(1, &vbo.ID)

	runtime.SetFinalizer(vbo, (*VBO).Delete)

	return vbo
}

// Bind binds the VBO
func (vbo *VBO) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.ID)
}

// Draw draws the VBO
func (vbo *VBO) Draw() {
	if vbo.size == 0 {
		// avoid drawing empty VBO
		return
	}

	vbo.Bind()
	gl.DrawArrays(uint32(vbo.DrawMode), 0, int32(vbo.Size()))
}

// Unbind unbinds the VBO
func (vbo *VBO) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

// Delete deletes the VBO
func (vbo *VBO) Delete() {
	gl.DeleteBuffers(1, &vbo.ID)
}

// Size returns the size of the VBO in floats
func (vbo *VBO) Size() int {
	return vbo.size
}

// Len returns the length of the VBO in bytes
func (vbo *VBO) Len() int {
	return vbo.size * 4
}

// SetData uploads data to the VBO
func (vbo *VBO) SetData(data []float32) {
	vbo.size = len(data)

	vbo.Bind()
	gl.BufferData(gl.ARRAY_BUFFER, vbo.Len(), gl.Ptr(data), uint32(vbo.Usage))
}

// GetData returns the data of the VBO
func (vbo *VBO) GetData() []float32 {
	vbo.Bind()

	data := make([]float32, vbo.size)
	gl.GetBufferSubData(gl.ARRAY_BUFFER, 0, vbo.Len(), gl.Ptr(data))

	return data
}

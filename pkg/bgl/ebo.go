package bgl

import (
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type EBO struct {
	Usage
	DrawMode
	ID uint32

	size int
}

// NewEBO creates a new empty EBO
func NewEBO() *EBO {
	ebo := &EBO{
		Usage:    DynamicDraw,
		DrawMode: Triangles,
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo.ID)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 0, nil, gl.STATIC_DRAW)

	runtime.SetFinalizer(ebo, (*EBO).Delete)

	return ebo
}

// Bind binds the EBO
func (ebo *EBO) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo.ID)
}

// Unbind unbinds the EBO
func (ebo *EBO) Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

// Size returns the size of the EBO in floats
func (ebo *EBO) Size() int {
	return ebo.size
}

// Len returns the length of the EBO in bytes
func (ebo *EBO) Len() int {
	return ebo.size * 4
}

// Draw draws the EBO
func (ebo *EBO) Draw() {
	ebo.Bind()
	gl.DrawElements(uint32(ebo.DrawMode), int32(ebo.size), gl.UNSIGNED_INT, nil)
}

// Delete
func (ebo *EBO) Delete() {
	gl.DeleteBuffers(1, &ebo.ID)
}

// SetData sets the data of the EBO
func (ebo *EBO) SetData(data []uint32) {
	ebo.size = len(data)

	ebo.Bind()
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, ebo.Len(), gl.Ptr(data), uint32(ebo.Usage))
}

package bgl

import (
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// VAO is a Vertex Array Object
type VAO struct {
	DrawMode
	AttrFormat
	ID uint32

	size int
}

// NewVAO creates a new VAO
func NewVAO(format AttrFormat) *VAO {
	vao := &VAO{
		AttrFormat: format,
		DrawMode:   Triangles,
	}

	gl.GenVertexArrays(1, &vao.ID)

	vao.Bind()

	offsets := vao.AttrFormat.Offsets()
	for i, attr := range vao.AttrFormat {
		gl.VertexAttribPointer(
			uint32(attr.Loc),
			int32(attr.Size()),
			gl.FLOAT,
			attr.Normalized,
			int32(vao.AttrFormat.Len()),
			gl.PtrOffset(offsets[i]),
		)
		gl.EnableVertexAttribArray(uint32(attr.Loc))
	}

	vao.Unbind()

	runtime.SetFinalizer(vao, (*VAO).Delete)

	return vao
}

// Bind binds the VAO
func (vao *VAO) Bind() {
	gl.BindVertexArray(vao.ID)
}

// Unbind unbinds the VAO
func (vao *VAO) Unbind() {
	gl.BindVertexArray(0)
}

// Draw draws the VAO
// func (vao *VAO) Draw(count int32) {
// 	if vao.size == 0 {
// 		// avoid drawing empty VAOs
// 		return
// 	}

// 	gl.DrawArrays(uint32(vao.DrawMode), 0, count)
// }

// Delete deletes the VAO
func (vao *VAO) Delete() {
	gl.DeleteVertexArrays(1, &vao.ID)
}

// Size returns the size of the VAO in floats
func (vao *VAO) Size() int {
	return vao.size
}

// Len returns the size of the VAO in bytes
func (vao *VAO) Len() int {
	return vao.size * 4
}

// SetData sets the data of the VAO
func (vao *VAO) SetData(data []float32) {
	vao.size = len(data)

	vao.Bind()
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, vao.Len(), gl.Ptr(data))
}

// GetData gets the data of the VAO
func (vao *VAO) GetData() []float32 {
	data := make([]float32, vao.size)

	vao.Bind()
	gl.GetBufferSubData(gl.ARRAY_BUFFER, 0, vao.Len(), gl.Ptr(data))

	return data
}

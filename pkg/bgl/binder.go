package bgl

import "github.com/go-gl/gl/v4.6-core/gl"

const (
	CurrentProgram     uint32 = gl.CURRENT_PROGRAM
	VertexArrayBinding uint32 = gl.VERTEX_ARRAY_BINDING
	ArrayBufferBinding uint32 = gl.ARRAY_BUFFER_BINDING
)

type binder struct {
	binding  uint32
	bindFunc func(uint32)

	id   uint32
	prev []uint32
}

func (b *binder) bind() *binder {
	var prev int32
	gl.GetIntegerv(b.binding, &prev)
	b.prev = append(b.prev, uint32(prev))

	if b.prev[len(b.prev)-1] != b.id {
		b.bindFunc(b.id)
	}

	return b
}

func (b *binder) restore() *binder {
	if b.prev[len(b.prev)-1] != b.id {
		b.bindFunc(b.prev[len(b.prev)-1])
	}

	b.prev = b.prev[:len(b.prev)-1]

	return b
}

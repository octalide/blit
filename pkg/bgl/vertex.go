package bgl

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type DrawMode uint32

const (
	Points                 = DrawMode(gl.POINTS)
	LineStrip              = DrawMode(gl.LINE_STRIP)
	LineLoop               = DrawMode(gl.LINE_LOOP)
	Lines                  = DrawMode(gl.LINES)
	LineStropAdjacency     = DrawMode(gl.LINE_STRIP_ADJACENCY)
	LinesAdjacency         = DrawMode(gl.LINES_ADJACENCY)
	TriangleStrip          = DrawMode(gl.TRIANGLE_STRIP)
	TriangleFan            = DrawMode(gl.TRIANGLE_FAN)
	Triangles              = DrawMode(gl.TRIANGLES)
	TriangleStripAdjacency = DrawMode(gl.TRIANGLE_STRIP_ADJACENCY)
	TrianglesAdjacency     = DrawMode(gl.TRIANGLES_ADJACENCY)
	Patches                = DrawMode(gl.PATCHES)
)

type VertexArray struct {
	DrawMode

	vao, vbo binder
	format   map[string]Attr
	prog     *Program

	cap    int
	stride int
	offset map[string]int
}

func newVertexArray(prog *Program, cap int, mode DrawMode) (*VertexArray, error) {
	// minimum capacity is 4
	if cap < 4 {
		cap = 4
	}

	va := &VertexArray{
		DrawMode: mode,
		vao: binder{
			binding: VertexArrayBinding,
			bindFunc: func(id uint32) {
				gl.BindVertexArray(id)
			},
		},
		vbo: binder{
			binding: ArrayBufferBinding,
			bindFunc: func(id uint32) {
				gl.BindBuffer(gl.ARRAY_BUFFER, id)
			},
		},
		cap:    cap,
		format: prog.VertexAttrs,
		stride: prog.VertexAttrs.Size(),
		offset: make(map[string]int, len(prog.VertexAttrs)),
		prog:   prog,
	}

	offset := 0
	for i, attr := range va.format {
		switch AttrType(attr.Type) {
		case Float, Vec2f, Vec3f, Vec4f:
		default:
			return nil, fmt.Errorf("invalid attribute type")
		}

		va.offset[i] = offset
		offset += int(attr.Size)
	}

	gl.GenVertexArrays(1, &va.vao.id)
	va.vao.bind()

	gl.GenBuffers(1, &va.vbo.id)
	defer va.vbo.bind().restore()

	emptyData := make([]byte, cap*va.stride)
	gl.BufferData(gl.ARRAY_BUFFER, len(emptyData), gl.Ptr(emptyData), gl.DYNAMIC_DRAW)

	for i, attr := range va.format {
		loc := gl.GetAttribLocation(prog.id, gl.Str(attr.Name+"\x00"))

		gl.VertexAttribPointer(
			uint32(loc),
			attr.Size,
			gl.FLOAT,
			false,
			int32(va.stride),
			gl.PtrOffset(va.offset[i]),
		)
		gl.EnableVertexAttribArray(uint32(loc))
	}

	va.vao.restore()

	runtime.SetFinalizer(va, (*VertexArray).delete)

	return va, nil
}

func (va *VertexArray) delete() {
	gl.DeleteVertexArrays(1, &va.vao.id)
	gl.DeleteBuffers(1, &va.vbo.id)
}

func (va *VertexArray) bind() {
	va.vao.bind()
	va.vbo.bind()
}

func (va *VertexArray) unbind() {
	va.vbo.restore()
	va.vao.restore()
}

func (va *VertexArray) draw(i, j int) {
	gl.DrawArrays(uint32(va.DrawMode), int32(i), int32(j-i))
}

func (va *VertexArray) setVertexData(i, j int, data []float32) {
	if j-i == 0 {
		// avoid setting 0 bytes of buffer data
		return
	}

	gl.BufferSubData(gl.ARRAY_BUFFER, i*va.stride, len(data)*4, gl.Ptr(data))
}

func (va *VertexArray) vertexData(i, j int) []float32 {
	if j-i == 0 {
		// avoid getting 0 bytes of buffer data
		return nil
	}

	data := make([]float32, (j-i)*va.stride/4)
	gl.GetBufferSubData(gl.ARRAY_BUFFER, i*va.stride, len(data)*4, gl.Ptr(data))

	return data
}

// VertexSlice points to a portion of (or possibly whole) vertex array. It is used as a pointer,
// contrary to Go's builtin slices. This is, so that append can be 'in-place'. That's for the good,
// because Begin/End-ing a VertexSlice would become super confusing, if append returned a new
// VertexSlice.
//
// It also implements all basic slice-like operations: appending, sub-slicing, etc.
//
// Note that you need to Begin a VertexSlice before getting or updating it's elements or drawing it.
// After you're done with it, you need to End it.
type VertexSlice struct {
	DrawMode
	va   *VertexArray
	i, j int
}

func NewVertexSlice(prog *Program, len, cap int, mode DrawMode) (*VertexSlice, error) {
	if len > cap {
		return nil, fmt.Errorf("length is greater than capacity")
	}

	vs := &VertexSlice{
		DrawMode: mode,
		i:        0,
		j:        len,
	}

	va, err := newVertexArray(prog, cap, mode)
	if err != nil {
		return nil, err
	}

	vs.va = va

	return vs, nil
}

// VertexFormat returns the format of vertex attributes inside the underlying vertex array of this
// VertexSlice.
func (vs *VertexSlice) VertexFormat() AttrFormat {
	return vs.va.format
}

// Stride returns the number of float32 elements occupied by one vertex.
func (vs *VertexSlice) Stride() int {
	return vs.va.stride / 4
}

// Len returns the length of the VertexSlice (number of vertices).
func (vs *VertexSlice) Len() int {
	return vs.j - vs.i
}

// Cap returns the capacity of an underlying vertex array.
func (vs *VertexSlice) Cap() int {
	return vs.va.cap - vs.i
}

// SetLen resizes the VertexSlice to length len.
func (vs *VertexSlice) SetLen(len int) {
	vs.Unbind() // vs must have been Begin-ed before calling this method
	*vs = vs.grow(len)
	vs.Bind()
}

// grow returns supplied vs with length changed to len. Allocates new underlying vertex array if
// necessary. The original content is preserved.
func (vs VertexSlice) grow(len int) VertexSlice {
	if len <= vs.Cap() {
		// capacity sufficient
		return VertexSlice{
			va: vs.va,
			i:  vs.i,
			j:  vs.i + len,
		}
	}

	// grow the capacity
	newCap := vs.Cap()
	if newCap < 1024 {
		newCap += newCap
	} else {
		newCap += newCap / 4
	}
	if newCap < len {
		newCap = len
	}

	va, _ := newVertexArray(vs.va.prog, newCap, vs.DrawMode)
	newVs := VertexSlice{
		va: va,
		i:  0,
		j:  len,
	}

	// preserve the original content
	newVs.Bind()
	newVs.Slice(0, vs.Len()).SetVertexData(vs.VertexData())
	newVs.Unbind()

	return newVs
}

// Slice returns a sub-slice of this VertexSlice covering the range [i, j) (relative to this
// VertexSlice).
//
// Note, that the returned VertexSlice shares an underlying vertex array with the original
// VertexSlice. Modifying the contents of one modifies corresponding contents of the other.
func (vs *VertexSlice) Slice(i, j int) *VertexSlice {
	if i < 0 || j < i || j > vs.va.cap {
		panic("failed to slice vertex slice: index out of range")
	}

	return &VertexSlice{
		va: vs.va,
		i:  vs.i + i,
		j:  vs.i + j,
	}
}

// SetVertexData sets the contents of the VertexSlice.
//
// The data is a slice of float32's, where each vertex attribute occupies a certain number of
// elements. Namely, Float occupies 1, Vec2 occupies 2, Vec3 occupies 3 and Vec4 occupies 4. The
// attribues in the data slice must be in the same order as in the vertex format of this Vertex
// Slice.
//
// If the length of vertices does not match the length of the VertexSlice, this methdo panics.
func (vs *VertexSlice) SetVertexData(data []float32) error {
	if len(data)/vs.Stride() != vs.Len() {
		return fmt.Errorf("invalid vertex data length")
	}

	vs.va.setVertexData(vs.i, vs.j, data)

	return nil
}

func (vs *VertexSlice) VertexData() []float32 {
	return vs.va.vertexData(vs.i, vs.j)
}

func (vs *VertexSlice) Draw() {
	vs.va.draw(vs.i, vs.j)
}

func (vs *VertexSlice) Bind() {
	vs.va.bind()
}

func (vs *VertexSlice) Unbind() {
	vs.va.unbind()
}

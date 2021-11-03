package bgl

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type AttrFormat []Attr

// Find returns the attribute with the given name
func (a AttrFormat) Find(name string) (attr Attr, ok bool) {
	for _, attr = range a {
		if attr.Name == name {
			ok = true
			return
		}
	}

	return
}

// Size gets the size in elements of the attribute format
func (a AttrFormat) Size() (total int) {
	for _, attr := range a {
		total += attr.Size()
	}

	return
}

// Len gets the size in bytes of the attrbute format
func (a AttrFormat) Len() (total int) {
	for _, attr := range a {
		total += attr.Len()
	}

	return
}

// Offsets returns the offsets of the attribute in the buffer
func (a AttrFormat) Offsets() (offsets []int) {
	offsets = make([]int, len(a))
	offset := 0
	for i, attr := range a {
		offsets[i] = offset
		offset += attr.Len()
	}

	return
}

type Attr struct {
	Type       AttrType
	Name       string
	Loc        int32
	Normalized bool
}

// Size gets the size in floats of the attribute
func (a Attr) Size() int {
	return int(a.Type.Size())
}

// Len gets the size in bytes of the attribute
func (a Attr) Len() int {
	return int(a.Type.Len())
}

func (a Attr) String() string {
	return fmt.Sprintf("{name: %v, type: %v, loc: %v, size: %v, len: %v, norm: %v}", a.Name, a.Type, a.Loc, a.Size(), a.Len(), a.Normalized)
}

type AttrType uint32

const (
	Float   = AttrType(gl.FLOAT)
	Vec2f   = AttrType(gl.FLOAT_VEC2)
	Vec3f   = AttrType(gl.FLOAT_VEC3)
	Vec4f   = AttrType(gl.FLOAT_VEC4)
	Mat2f   = AttrType(gl.FLOAT_MAT2)
	Mat3f   = AttrType(gl.FLOAT_MAT3)
	Mat4f   = AttrType(gl.FLOAT_MAT4)
	Mat2x3f = AttrType(gl.FLOAT_MAT2x3)
	Mat2x4f = AttrType(gl.FLOAT_MAT2x4)
	Mat3x2f = AttrType(gl.FLOAT_MAT3x2)
	Mat3x4f = AttrType(gl.FLOAT_MAT3x4)
	Mat4x2f = AttrType(gl.FLOAT_MAT4x2)
	Mat4x3f = AttrType(gl.FLOAT_MAT4x3)
	Int     = AttrType(gl.INT)
	Vec2i   = AttrType(gl.INT_VEC2)
	Vec3i   = AttrType(gl.INT_VEC3)
	Vec4i   = AttrType(gl.INT_VEC4)
	UInt    = AttrType(gl.UNSIGNED_INT)
	Vec2ui  = AttrType(gl.UNSIGNED_INT_VEC2)
	Vec3ui  = AttrType(gl.UNSIGNED_INT_VEC3)
	Vec4ui  = AttrType(gl.UNSIGNED_INT_VEC4)
	Double  = AttrType(gl.DOUBLE)
	Vec2d   = AttrType(gl.DOUBLE_VEC2)
	Vec3d   = AttrType(gl.DOUBLE_VEC3)
	Vec4d   = AttrType(gl.DOUBLE_VEC4)
	Mat2d   = AttrType(gl.DOUBLE_MAT2)
	Mat3d   = AttrType(gl.DOUBLE_MAT3)
	Mat4d   = AttrType(gl.DOUBLE_MAT4)
	Mat2x3d = AttrType(gl.DOUBLE_MAT2x3)
	Mat2x4d = AttrType(gl.DOUBLE_MAT2x4)
	Mat3x2d = AttrType(gl.DOUBLE_MAT3x2)
	Mat3x4d = AttrType(gl.DOUBLE_MAT3x4)
	Mat4x2d = AttrType(gl.DOUBLE_MAT4x2)
	Mat4x3d = AttrType(gl.DOUBLE_MAT4x3)
)

// Size gets the size in elements of an attribute type
func (a AttrType) Size() int {
	switch a {
	case Float:
		return 1
	case Vec2f, Vec2i, Vec2ui, Vec2d:
		return 2
	case Vec3f, Vec3i, Vec3ui, Vec3d:
		return 3
	case Vec4f, Vec4i, Vec4ui, Vec4d:
		return 4
	case Mat2f, Mat2d:
		return 4
	case Mat3f, Mat3d:
		return 9
	case Mat4f, Mat4d:
		return 16
	case Mat2x3f, Mat2x3d:
		return 6
	case Mat2x4f, Mat2x4d:
		return 8
	case Mat3x2f, Mat3x2d:
		return 6
	case Mat3x4f, Mat3x4d:
		return 12
	case Mat4x2f, Mat4x2d:
		return 8
	case Mat4x3f, Mat4x3d:
		return 12
	default:
		panic(fmt.Sprintf("Unknown attribute type: %d", a))
	}
}

// Len gets the size in bytes of an attribute type.
func (a AttrType) Len() int {
	switch a {
	case Float:
		return 4
	case Vec2f:
		return 8
	case Vec3f:
		return 12
	case Vec4f:
		return 16
	case Mat2f:
		return 16
	case Mat3f:
		return 36
	case Mat4f:
		return 64
	case Mat2x3f:
		return 24
	case Mat2x4f:
		return 32
	case Mat3x2f:
		return 24
	case Mat3x4f:
		return 48
	case Mat4x2f:
		return 32
	case Mat4x3f:
		return 48
	case Int:
		return 4
	case Vec2i:
		return 8
	case Vec3i:
		return 12
	case Vec4i:
		return 16
	case UInt:
		return 4
	case Vec2ui:
		return 8
	case Vec3ui:
		return 12
	case Vec4ui:
		return 16
	case Double:
		return 8
	case Vec2d:
		return 16
	case Vec3d:
		return 24
	case Vec4d:
		return 32
	case Mat2d:
		return 32
	case Mat3d:
		return 72
	case Mat4d:
		return 128
	case Mat2x3d:
		return 48
	case Mat2x4d:
		return 64
	case Mat3x2d:
		return 48
	case Mat3x4d:
		return 96
	case Mat4x2d:
		return 64
	case Mat4x3d:
		return 96
	default:
		return 0
	}
}

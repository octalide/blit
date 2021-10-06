package bgl

import "github.com/go-gl/gl/v4.6-core/gl"

type AttrFormat map[string]Attr

func (a AttrFormat) Size() (total int) {
	for _, attr := range a {
		total += int(attr.Size)
	}

	return
}

type Attr struct {
	Name string
	Type uint32
	Loc  int32
	Size int32
	Len  int32
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

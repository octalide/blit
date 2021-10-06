package blit

import (
	"fmt"
	"math"
)

type Vec struct {
	X, Y float64
}

func (v Vec) Normalize() Vec {
	l := 1.0 / v.Len()
	return Vec{v.X * l, v.Y * l}
}

func (v Vec) Inv() Vec {
	return Vec{-v.X, -v.Y}
}

func (v Vec) Len() float64 {
	return math.Hypot(v.X, v.Y)
}

func (v Vec) Eql(v2 Vec) bool {
	return v.X == v2.X && v.Y == v2.Y
}

func (v Vec) Add(v2 Vec) Vec {
	return Vec{v.X + v2.X, v.Y + v2.Y}
}

func (v Vec) Sub(v2 Vec) Vec {
	return v.Add(v2.Inv())
}

func (v Vec) Mul(s float64) Vec {
	return Vec{v.X * s, v.Y * s}
}

func (v Vec) Dot(v2 Vec) float64 {
	return (v.X * v2.X) + (v.Y * v2.Y)
}

func (v Vec) Rot(angle float64) Vec {
	sin, cos := math.Sincos(angle)
	return Vec{
		v.X*cos - v.Y*sin,
		v.X*sin + v.Y*cos,
	}
}

func (v Vec) F64() []float64 {
	return []float64{v.X, v.Y}
}

func (v Vec) F32() []float32 {
	return []float32{float32(v.X), float32(v.Y)}
}

type Rect struct {
	Min, Max Vec
}

func (r Rect) Width() float64 {
	return r.Max.X - r.Min.X
}

func (r Rect) Height() float64 {
	return r.Max.Y - r.Min.Y
}

func (r Rect) Size() Vec {
	return Vec{r.Width(), r.Height()}
}

func (r Rect) Area() float64 {
	return r.Width() * r.Height()
}

func (r Rect) Norm() Rect {
	return Rect{
		Min: Vec{
			math.Min(r.Min.X, r.Max.X),
			math.Min(r.Min.Y, r.Max.Y),
		},
		Max: Vec{
			math.Max(r.Min.X, r.Max.X),
			math.Max(r.Min.Y, r.Max.Y),
		},
	}
}

func (r Rect) Vert() [4]Vec {
	return [4]Vec{
		r.Min,
		{r.Min.X, r.Max.Y},
		r.Max,
		{r.Max.X, r.Min.Y},
	}
}

// Mat is a 2x3 affine matrix that can be used for all kinds of spatial transforms, such
// as movement, scaling and rotations.
//
// Mat has a handful of useful methods, each of which adds a transformation to the matrix. For
// example:
//
//   pixel.IM.Moved(pixel.V(100, 200)).Rotated(pixel.ZV, math.Pi/2)
//
// This code creates a Mat that first moves everything by 100 units horizontally and 200 units
// vertically and then rotates everything by 90 degrees around the origin.
//
// Layout is:
// [0] [2] [4]
// [1] [3] [5]
//  0   0   1  (implicit row)
type Mat [6]float64

func Ident() Mat {
	return Mat{1, 0, 0, 1, 0, 0}
}

func (m Mat) Add(m2 Mat) Mat {
	return Mat{m[0] + m2[0], m[1] + m2[1], m[2] + m2[2], m[3] + m2[3], m[4] + m2[4], m[5] + m2[5]}
}

func (m Mat) Sub(m2 Mat) Mat {
	return Mat{m[0] - m2[0], m[1] - m2[1], m[2] - m2[2], m[3] - m2[3], m[4] - m2[4], m[5] - m2[5]}
}

func (m Mat) Mul(s float64) Mat {
	return Mat{m[0] * s, m[1] * s, m[2] * s, m[3] * s, m[4] * s, m[5] * s}
}

func (m Mat) Pos(v Vec) Mat {
	m[4], m[5] = m[4]+v.X, m[5]+v.Y
	return m
}

func (m Mat) Scale(origin Vec, scale Vec) Mat {
	m[4] -= origin.X
	m[5] -= origin.Y
	m[0] *= scale.X
	m[2] *= scale.X
	m[4] *= scale.X
	m[1] *= scale.Y
	m[3] *= scale.Y
	m[5] *= scale.Y
	m[4] += origin.X
	m[5] += origin.Y

	return m
}

func (m Mat) ScaleS(around Vec, scale float64) Mat {
	return m.Scale(around, Vec{scale, scale})
}

func (m Mat) Rot(around Vec, rad float64) Mat {
	m[4] -= around.X
	m[5] -= around.Y

	sin, cos := math.Sincos(rad)
	m = m.Chain(Mat{cos, sin, -sin, cos, 0, 0})

	m[4] += around.X
	m[5] += around.Y

	return m
}

// Chain adds another Matrix to this one. All tranformations by the next Matrix will be applied
// after the transformations of this Matrix.
func (m Mat) Chain(m2 Mat) Mat {
	return Mat{
		m2[0]*m[0] + m2[2]*m[1],
		m2[1]*m[0] + m2[3]*m[1],
		m2[0]*m[2] + m2[2]*m[3],
		m2[1]*m[2] + m2[3]*m[3],
		m2[0]*m[4] + m2[2]*m[5] + m2[4],
		m2[1]*m[4] + m2[3]*m[5] + m2[5],
	}
}

// Project applies all transformations added to the Matrix to a vector u and returns the result.
//
// Time complexity is O(1).
func (m Mat) Project(u Vec) Vec {
	return Vec{
		m[0]*u.X + m[2]*u.Y + m[4],
		m[1]*u.X + m[3]*u.Y + m[5],
	}
}

// Unproject does the inverse operation to Project.
//
// Time complexity is O(1).
func (m Mat) Unproject(u Vec) Vec {
	det := m[0]*m[3] - m[2]*m[1]

	return Vec{
		(m[3]*(u.X-m[4]) - m[2]*(u.Y-m[5])) / det,
		(-m[1]*(u.X-m[4]) + m[0]*(u.Y-m[5])) / det,
	}
}

func (m Mat) String() string {
	return fmt.Sprintf("{%v %v %v | %v %v %v}", m[0], m[2], m[4], m[1], m[3], m[5])
}

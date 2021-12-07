package blit

import (
	"fmt"
	"math"
)

// Vec is a 4D vector
type Vec [4]float32

// F returns the vector as a [4]float32 array (convenience function)
func (v Vec) F() [4]float32 {
	return [4]float32(v)
}

// X returns the x component of the vector
func (v Vec) X() float32 {
	return v[0]
}

// Y returns the y component of the vector
func (v Vec) Y() float32 {
	return v[1]
}

// Z returns the z component of the vector
func (v Vec) Z() float32 {
	return v[2]
}

// W returns the w component of the vector
func (v Vec) W() float32 {
	return v[3]
}

// Nrm normalizes the vector
func (v Vec) Nrm() Vec {
	return v.Scl(1.0 / v.Len())
}

// Ref reflects the vector based on a normal
func (v Vec) Ref(n Vec) Vec {
	return v.Sub(n.Scl(2 * v.Dot(n)))
}

// Scl multiplies a vector by a scalar value
func (v Vec) Scl(s float32) Vec {
	return Vec{v[0] * s, v[1] * s, v[2] * s, v[3] * s}
}

// Len returns the length of the vector
func (v Vec) Len() float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2] + v[3]*v[3])))
}

// Add adds two vectors
func (v Vec) Add(v2 Vec) Vec {
	return Vec{v[0] + v2[0], v[1] + v2[1], v[2] + v2[2], v[3] + v2[3]}
}

// Sub subtracts two vectors
func (v Vec) Sub(v2 Vec) Vec {
	return Vec{v[0] - v2[0], v[1] - v2[1], v[2] - v2[2], v[3] - v2[3]}
}

// Mul multiplies two vectors
func (v Vec) Mul(v2 Vec) Vec {
	return Vec{v[0] * v2[0], v[1] * v2[1], v[2] * v2[2], v[3] * v2[3]}
}

// Crs returns the cross product of two vectors
func (v Vec) Crs(v2 Vec) Vec {
	return Vec{
		v[1]*v2[2] - v[2]*v2[1],
		v[2]*v2[0] - v[0]*v2[2],
		v[0]*v2[1] - v[1]*v2[0],
	}
}

// Dot returns the dot product of two vectors
func (v Vec) Dot(v2 Vec) float32 {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2] + v[3]*v2[3]
}

// Rot rotates the vector by the given angle
func (v Vec) Rot(angle float32) Vec {
	return v.Mat(Ident().Rot(angle))
}

// Mat multiples the vector by a matrix
func (v Vec) Mat(m Mat) Vec {
	return Vec{
		v[0]*m[0] + v[1]*m[4] + v[2]*m[8] + v[3]*m[12],
		v[0]*m[1] + v[1]*m[5] + v[2]*m[9] + v[3]*m[13],
		v[0]*m[2] + v[1]*m[6] + v[2]*m[10] + v[3]*m[14],
		v[0]*m[3] + v[1]*m[7] + v[2]*m[11] + v[3]*m[15],
	}
}

// Inv returns the inverse of the vector
func (v Vec) Inv() Vec {
	return v.Scl(-1)
}

// String returns a string representation of the vector, with each value
// truncated to 3 decimal places
func (v Vec) String() string {
	return fmt.Sprintf("[%+.3f, %+.3f, %+.3f, %+.3f]", v[0], v[1], v[2], v[3])
}

// Rect represents a rectangle
type Rect [4]float32

// X returns the x component of the rectangle
func (r Rect) X() float32 {
	return r[0]
}

// Y returns the y component of the rectangle
func (r Rect) Y() float32 {
	return r[1]
}

// W returns the w component of the rectangle
func (r Rect) W() float32 {
	return r[2]
}

// H returns the h component of the rectangle
func (r Rect) H() float32 {
	return r[3]
}

// Min returns the minimum vector of the rectangle
func (r Rect) Min() Vec {
	return Vec{r.X(), r.Y()}
}

// Max returns the maximum vector of the rectangle
func (r Rect) Max() Vec {
	return r.Min().Add(Vec{r.W(), r.H()})
}

// Size returns a vector with the size of the rectangle in each dimension
func (r Rect) Size() Vec {
	return Vec{r.W(), r.H()}
}

// Area returns the area of the rectangle
func (r Rect) Area() float32 {
	return r.W() * r.H()
}

// Mat is a 4x4 matrix
type Mat [16]float32

// Ident returns the identity matrix
func Ident() Mat {
	return Mat{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// F returns the matrix as a [16]float32 array (convenience function)
func (m Mat) F() [16]float32 {
	return [16]float32(m)
}

// Add adds two matrices
func (m Mat) Add(m2 Mat) Mat {
	return Mat{
		m[0] + m2[0], m[1] + m2[1], m[2] + m2[2], m[3] + m2[3],
		m[4] + m2[4], m[5] + m2[5], m[6] + m2[6], m[7] + m2[7],
		m[8] + m2[8], m[9] + m2[9], m[10] + m2[10], m[11] + m2[11],
		m[12] + m2[12], m[13] + m2[13], m[14] + m2[14], m[15] + m2[15],
	}
}

// Sub subtracts two matrices
func (m Mat) Sub(m2 Mat) Mat {
	return m.Add(m2.Scl(-1))
}

// Scl multiplies a matrix by a scalar
func (m Mat) Scl(s float32) Mat {
	return Mat{
		m[0] * s, m[1] * s, m[2] * s, m[3] * s,
		m[4] * s, m[5] * s, m[6] * s, m[7] * s,
		m[8] * s, m[9] * s, m[10] * s, m[11] * s,
		m[12] * s, m[13] * s, m[14] * s, m[15] * s,
	}
}

// Pos transforms a matrix by a position vector
func (m Mat) Pos(v Vec) Mat {
	return m.Mul(Mat{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		v.X(), v.Y(), v.Z(), 1,
	})
}

// Det gets the determinant of a 4x4 matrix
func (m Mat) Det() float32 {
	return m[0]*m[5]*m[10]*m[15] +
		m[0]*m[6]*m[9]*m[14] +
		m[0]*m[7]*m[8]*m[13] +
		m[1]*m[4]*m[10]*m[15] +
		m[1]*m[6]*m[8]*m[12] +
		m[1]*m[7]*m[9]*m[14] +
		m[2]*m[4]*m[9]*m[15] +
		m[2]*m[5]*m[8]*m[13] +
		m[2]*m[7]*m[10]*m[12] +
		m[3]*m[4]*m[8]*m[14] +
		m[3]*m[5]*m[9]*m[12] +
		m[3]*m[6]*m[10]*m[13] -
		m[0]*m[5]*m[9]*m[15] -
		m[0]*m[6]*m[10]*m[13] -
		m[0]*m[7]*m[11]*m[12] -
		m[1]*m[4]*m[10]*m[14] -
		m[1]*m[6]*m[11]*m[12] -
		m[1]*m[7]*m[8]*m[13] -
		m[2]*m[4]*m[11]*m[14] -
		m[2]*m[5]*m[8]*m[12] -
		m[2]*m[7]*m[9]*m[13] -
		m[3]*m[4]*m[9]*m[14] -
		m[3]*m[5]*m[10]*m[12] -
		m[3]*m[6]*m[8]*m[13]
}

// Mul performs a matrix product
func (m Mat) Mul(m2 Mat) Mat {
	return Mat{
		m[0]*m2[0] + m[4]*m2[1] + m[8]*m2[2] + m[12]*m2[3],
		m[1]*m2[0] + m[5]*m2[1] + m[9]*m2[2] + m[13]*m2[3],
		m[2]*m2[0] + m[6]*m2[1] + m[10]*m2[2] + m[14]*m2[3],
		m[3]*m2[0] + m[7]*m2[1] + m[11]*m2[2] + m[15]*m2[3],
		m[0]*m2[4] + m[4]*m2[5] + m[8]*m2[6] + m[12]*m2[7],
		m[1]*m2[4] + m[5]*m2[5] + m[9]*m2[6] + m[13]*m2[7],
		m[2]*m2[4] + m[6]*m2[5] + m[10]*m2[6] + m[14]*m2[7],
		m[3]*m2[4] + m[7]*m2[5] + m[11]*m2[6] + m[15]*m2[7],
		m[0]*m2[8] + m[4]*m2[9] + m[8]*m2[10] + m[12]*m2[11],
		m[1]*m2[8] + m[5]*m2[9] + m[9]*m2[10] + m[13]*m2[11],
		m[2]*m2[8] + m[6]*m2[9] + m[10]*m2[10] + m[14]*m2[11],
		m[3]*m2[8] + m[7]*m2[9] + m[11]*m2[10] + m[15]*m2[11],
		m[0]*m2[12] + m[4]*m2[13] + m[8]*m2[14] + m[12]*m2[15],
		m[1]*m2[12] + m[5]*m2[13] + m[9]*m2[14] + m[13]*m2[15],
		m[2]*m2[12] + m[6]*m2[13] + m[10]*m2[14] + m[14]*m2[15],
		m[3]*m2[12] + m[7]*m2[13] + m[11]*m2[14] + m[15]*m2[15],
	}
}

// Rot rotates a matrix along the z axis
func (m Mat) Rot(rad float32) Mat {
	s, c := math.Sincos(float64(rad))
	sin := float32(s)
	cos := float32(c)

	return m.Mul(Mat{
		cos, -sin, 0, 0,
		sin, cos, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})
}

// Inv inverts the matrix
func (m Mat) Inv() Mat {
	d := m.Det()
	if d == 0 {
		return m
	}

	return Mat{
		m[5]*m[10]*m[15] + m[6]*m[11]*m[13] + m[7]*m[9]*m[14] - m[5]*m[11]*m[14] - m[6]*m[9]*m[15] - m[7]*m[10]*m[13],
		m[1]*m[11]*m[14] + m[2]*m[9]*m[15] + m[3]*m[10]*m[13] - m[1]*m[10]*m[15] - m[2]*m[11]*m[13] - m[3]*m[9]*m[14],
		m[1]*m[6]*m[15] + m[2]*m[7]*m[13] + m[3]*m[5]*m[14] - m[1]*m[7]*m[14] - m[2]*m[5]*m[15] - m[3]*m[6]*m[13],
		m[1]*m[7]*m[10] + m[2]*m[5]*m[11] + m[3]*m[6]*m[9] - m[1]*m[6]*m[11] - m[2]*m[7]*m[9] - m[3]*m[5]*m[10],
		m[4]*m[11]*m[14] + m[6]*m[8]*m[15] + m[7]*m[10]*m[12] - m[4]*m[10]*m[15] - m[6]*m[11]*m[12] - m[7]*m[8]*m[14],
		m[0]*m[10]*m[15] + m[2]*m[11]*m[12] + m[3]*m[8]*m[14] - m[0]*m[11]*m[14] - m[2]*m[8]*m[15] - m[3]*m[10]*m[12],
		m[0]*m[7]*m[14] + m[2]*m[4]*m[15] + m[3]*m[6]*m[12] - m[0]*m[6]*m[15] - m[2]*m[7]*m[12] - m[3]*m[4]*m[14],
		m[0]*m[6]*m[11] + m[2]*m[7]*m[8] + m[3]*m[4]*m[10] - m[0]*m[7]*m[10] - m[2]*m[4]*m[11] - m[3]*m[6]*m[8],
		m[4]*m[9]*m[15] + m[5]*m[11]*m[12] + m[7]*m[8]*m[13] - m[4]*m[11]*m[13] - m[5]*m[8]*m[15] - m[7]*m[9]*m[12],
		m[0]*m[11]*m[13] + m[1]*m[8]*m[15] + m[3]*m[9]*m[12] - m[0]*m[9]*m[15] - m[1]*m[11]*m[12] - m[3]*m[8]*m[13],
		m[0]*m[5]*m[15] + m[1]*m[7]*m[12] + m[3]*m[4]*m[13] - m[0]*m[7]*m[13] - m[1]*m[4]*m[15] - m[3]*m[5]*m[12],
		m[0]*m[7]*m[9] + m[1]*m[4]*m[11] + m[3]*m[5]*m[8] - m[0]*m[5]*m[11] - m[1]*m[7]*m[8] - m[3]*m[4]*m[9],
		m[4]*m[10]*m[13] + m[5]*m[8]*m[14] + m[6]*m[9]*m[12] - m[4]*m[9]*m[14] - m[5]*m[10]*m[12] - m[6]*m[8]*m[13],
		m[0]*m[9]*m[14] + m[1]*m[10]*m[12] + m[2]*m[8]*m[13] - m[0]*m[10]*m[13] - m[1]*m[8]*m[14] - m[2]*m[9]*m[12],
		m[0]*m[6]*m[13] + m[1]*m[4]*m[14] + m[2]*m[5]*m[12] - m[0]*m[4]*m[14] - m[1]*m[6]*m[12] - m[2]*m[5]*m[13],
		m[0]*m[4]*m[10] + m[1]*m[5]*m[8] + m[2]*m[6]*m[9] - m[0]*m[5]*m[9] - m[1]*m[6]*m[8] - m[2]*m[4]*m[10],
	}
}

// String returns a string representation of the matrix, with each value
// truncated to 3 decimal places
func (m Mat) String() string {
	return fmt.Sprintf("[%+.3f, %+.3f, %+.3f, %+.3f, %+.3f, %+.3f, %+.3f, %+.3f, %+.3f, %+.3f, %+.3f, %+.3f, %+.3f, %+.3f, %+.3f, %+.3f]",
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15])
}

/*
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
type Mat [6]float32

func Ident() Mat {
	return Mat{1, 0, 0, 1, 0, 0}
}

func (m Mat) Add(m2 Mat) Mat {
	return Mat{m[0] + m2[0], m[1] + m2[1], m[2] + m2[2], m[3] + m2[3], m[4] + m2[4], m[5] + m2[5]}
}

func (m Mat) Sub(m2 Mat) Mat {
	return Mat{m[0] - m2[0], m[1] - m2[1], m[2] - m2[2], m[3] - m2[3], m[4] - m2[4], m[5] - m2[5]}
}

func (m Mat) Mul(s float32) Mat {
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

func (m Mat) ScaleS(around Vec, scale float32) Mat {
	return m.Scale(around, Vec{scale, scale, scale})
}

func (m Mat) Rot(around Vec, rad float32) Mat {
	m[4] -= around.X
	m[5] -= around.Y

	s, c := math.Sincos(float64(rad))
	sin := float32(s)
	cos := float32(c)

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
		u.Z,
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
		u.Z,
	}
}

func (m Mat) String() string {
	return fmt.Sprintf("{%v %v %v | %v %v %v}", m[0], m[2], m[4], m[1], m[3], m[5])
}
*/

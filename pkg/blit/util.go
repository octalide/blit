package blit

import (
	"math"
)

// Ortho creates an orthographic projection matrix
func Ortho(left, right, bottom, top, near, far float32) Mat {
	return Mat{
		2 / (right - left), 0, 0, 0,
		0, 2 / (top - bottom), 0, 0,
		0, 0, -2 / (far - near), 0,
		-(right + left) / (right - left), -(top + bottom) / (top - bottom), -(far + near) / (far - near), 1,
	}
}

// Perspective creates a perspective matrix
func Perspective(fov, aspect, near, far float32) Mat {
	f := 1 / float32(math.Tan(float64(fov/2)))

	return Mat{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (far + near) / (near - far), 2 * far * near / (near - far),
		0, 0, -1, 0,
	}
}

// LookAt creates a transform matrix from world space to eye space
func LookAt(eye, center, up Vec) Mat {
	z := (eye.Sub(center)).Nrm()
	x := up.Crs(z).Nrm()
	y := z.Crs(x)

	return Mat{
		x.X(), y.X(), z.X(), 0,
		x.Y(), y.Y(), z.Y(), 0,
		x.Z(), y.Z(), z.Z(), 0,
		-x.Dot(eye), -y.Dot(eye), -z.Dot(eye), 1,
	}
}

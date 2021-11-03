package bgl

import "github.com/go-gl/gl/v4.6-core/gl"

type Usage uint32

const (
	DynamicDraw = Usage(gl.DYNAMIC_DRAW)
	StaticDraw  = Usage(gl.STATIC_DRAW)
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

type ShaderType uint32

func (st ShaderType) String() string {
	switch st {
	case FragShader:
		return "FRAG"
	case VertShader:
		return "VERT"
	case CompShader:
		return "COMP"
	case GeomShader:
		return "GEOM"
	}

	return "UNKN"
}

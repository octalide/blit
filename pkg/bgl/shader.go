package bgl

import (
	"embed"
	"fmt"
	"io/fs"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	FragShader ShaderType = gl.FRAGMENT_SHADER
	VertShader ShaderType = gl.VERTEX_SHADER
	CompShader ShaderType = gl.COMPUTE_SHADER
	GeomShader ShaderType = gl.GEOMETRY_SHADER
)

var (
	//go:embed src/*
	embedded embed.FS

	DefaultFrag string // Default frag source included for convenience
	DefaultVert string // Default vertex source included for convenience
)

// init loads the default shader sources
func init() {
	emb, _ := fs.Sub(embedded, "src")

	frag, _ := fs.ReadFile(emb, "default.frag")
	DefaultFrag = string(frag) + "\x00"

	vert, _ := fs.ReadFile(emb, "default.vert")
	DefaultVert = string(vert) + "\x00"
}

// DefaultProgram returns a program with the default shaders
func DefaultProgram() (*Program, error) {
	return NewProgram([]*Shader{
		NewShader(DefaultFrag, FragShader),
		NewShader(DefaultVert, VertShader),
	})
}

// Shader is an OpenGL shader
type Shader struct {
	stype    ShaderType
	ID       uint32
	src      string
	compiled bool
}

// NewShader creates a new shader
func NewShader(src string, stype ShaderType) *Shader {
	s := &Shader{
		src:   src,
		stype: stype,
	}

	return s
}

// Type returns the type of the shader
func (s *Shader) Type() ShaderType {
	return s.stype
}

// Src returns the source of the program
func (s *Shader) Src() string {
	return s.src
}

// Compiled returns true if the program is compiled
func (s *Shader) Compiled() bool {
	return s.compiled
}

// getiv returns the value of the given parameter
func (s *Shader) getiv(iv uint32) int32 {
	var result int32
	gl.GetShaderiv(s.ID, iv, &result)

	return result
}

// getInfoLog returns the info log for the shader
func (s *Shader) getInfoLog() string {
	logLength := s.getiv(gl.INFO_LOG_LENGTH)

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(s.ID, logLength, nil, gl.Str(log))

	return log
}

// compile compiles the shader
func (s *Shader) compile() error {
	s.ID = gl.CreateShader(uint32(s.stype))

	csources, free := gl.Strs(s.src)
	gl.ShaderSource(s.ID, 1, csources, nil)
	gl.CompileShader(s.ID)
	free()

	status := s.getiv(gl.COMPILE_STATUS)
	if status == gl.FALSE {
		log := s.getInfoLog()

		return fmt.Errorf("failed to compile shader (type: %v):\n%v", s.stype, log)
	}

	s.compiled = true

	return nil
}

// delete deletes the shader
func (s *Shader) delete() {
	gl.DeleteShader(s.ID)
}

// Program is an OpenGL shader program
type Program struct {
	ID uint32

	VertexAttrs  AttrFormat // vertex attribute format
	UniformAttrs AttrFormat // uniform attribute format

	compiled bool
}

// NewProgram creates a new program
func NewProgram(shaders []*Shader) (*Program, error) {
	p := &Program{
		VertexAttrs:  AttrFormat{},
		UniformAttrs: AttrFormat{},
	}

	if len(shaders) == 0 {
		return nil, fmt.Errorf("no shaders defined")
	}

	p.ID = gl.CreateProgram()

	for _, s := range shaders {
		if err := s.compile(); err != nil {
			return nil, fmt.Errorf("failed to compile shader program: %w", err)
		}

		p.attach(s)
	}

	p.link()

	status := p.getiv(gl.LINK_STATUS)
	if status == gl.FALSE {
		log := p.GetInfoLog()

		return nil, fmt.Errorf("failed to link shader program: %v", log)
	}

	for _, s := range shaders {
		s.delete()
	}

	p.compiled = true

	p.findUniforms()
	p.findAttributes()

	runtime.SetFinalizer(p, (*Program).Delete)

	return p, nil
}

// Compiled returns true if the program is compiled
func (p *Program) Compiled() bool {
	return p.compiled
}

// attach attaches a shader to the program
func (p *Program) attach(shader *Shader) {
	gl.AttachShader(p.ID, uint32(shader.ID))
}

// link links the program
func (p *Program) link() {
	gl.LinkProgram(p.ID)
}

// getiv returns the value of the given parameter
func (p *Program) getiv(iv uint32) int32 {
	var result int32
	gl.GetProgramiv(p.ID, iv, &result)

	return result
}

// GetInfoLog returns the info log for the program
func (p *Program) GetInfoLog() string {
	logLength := p.getiv(gl.INFO_LOG_LENGTH)

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(p.ID, logLength, nil, gl.Str(log))

	return log
}

// findUniforms finds the uniforms in the program
func (p *Program) findUniforms() {
	p.Bind()
	defer p.Unbind()

	active := uint32(p.getiv(gl.ACTIVE_UNIFORMS))
	p.UniformAttrs = make(AttrFormat, active)
	for i := uint32(0); i < active; i++ {
		var l int32     // length
		var s int32     // size
		var t uint32    // type
		var b [256]byte // name
		gl.GetActiveUniform(p.ID, i, 256, &l, &s, &t, &b[0])

		a := Attr{
			Name: gl.GoStr(&b[0]),
			Type: AttrType(t),
			Loc:  gl.GetUniformLocation(p.ID, &b[0]),
		}

		p.UniformAttrs[i] = a
	}
}

// findAttributes finds the attributes in the program
func (p *Program) findAttributes() {
	p.Bind()
	defer p.Unbind()

	active := uint32(p.getiv(gl.ACTIVE_ATTRIBUTES))
	p.VertexAttrs = make(AttrFormat, active)
	for i := uint32(0); i < active; i++ {
		var l int32     // length
		var s int32     // size
		var t uint32    // type
		var b [256]byte // name
		gl.GetActiveAttrib(p.ID, i, 256, &l, &s, &t, &b[0])

		a := Attr{
			Name: gl.GoStr(&b[0]),
			Type: AttrType(t),
			Loc:  gl.GetAttribLocation(p.ID, &b[0]),
		}

		p.VertexAttrs[i] = a
	}
}

// Bind binds the program
func (p *Program) Bind() {
	gl.UseProgram(p.ID)
}

// Unbind unbinds the program
func (p *Program) Unbind() {
	gl.UseProgram(0)
}

// SetUniform sets the value of the given uniform
func (p *Program) SetUniform(name string, value interface{}) error {
	var attr Attr

	attr, ok := p.VertexAttrs.Find(name)
	if !ok {
		attr, ok = p.UniformAttrs.Find(name)
		if !ok {
			return fmt.Errorf("attribute not found: %v", name)
		}
	}

	switch AttrType(attr.Type) {
	case Float:
		v := value.(float32)
		gl.Uniform1fv(attr.Loc, 1, &v)
	// case Vec2f:
	// 	v := value.(mgl32.Vec2)
	// 	gl.Uniform2fv(attr.Loc, 1, &v[0])
	// case Vec3f:
	// 	v := value.(mgl32.Vec3)
	// 	gl.Uniform3fv(attr.Loc, 1, &v[0])
	case Vec4f:
		v := value.([4]float32)
		gl.Uniform4fv(attr.Loc, 1, &v[0])
	// case Mat2f:
	// 	v := value.(mgl32.Mat2)
	// 	gl.UniformMatrix2fv(attr.Loc, 1, false, &v[0])
	// case Mat3f:
	// 	v := value.(mgl32.Mat3)
	// 	gl.UniformMatrix3fv(attr.Loc, 1, false, &v[0])
	case Mat4f:
		v := value.([16]float32)
		gl.UniformMatrix4fv(attr.Loc, 1, false, &v[0])
	// case Mat2x3f:
	// 	v := value.(mgl32.Mat2x3)
	// 	gl.UniformMatrix2x3fv(attr.Loc, 1, false, &v[0])
	// case Mat2x4f:
	// 	v := value.(mgl32.Mat2x4)
	// 	gl.UniformMatrix2x4fv(attr.Loc, 1, false, &v[0])
	// case Mat3x2f:
	// 	v := value.(mgl32.Mat3x2)
	// 	gl.UniformMatrix3x2fv(attr.Loc, 1, false, &v[0])
	// case Mat3x4f:
	// 	v := value.(mgl32.Mat3x4)
	// 	gl.UniformMatrix3x4fv(attr.Loc, 1, false, &v[0])
	// case Mat4x2f:
	// 	v := value.(mgl32.Mat4x2)
	// 	gl.UniformMatrix4x2fv(attr.Loc, 1, false, &v[0])
	// case Mat4x3f:
	// 	v := value.(mgl32.Mat4x3)
	// 	gl.UniformMatrix4x3fv(attr.Loc, 1, false, &v[0])
	case Int:
		v := value.(int32)
		gl.Uniform1iv(attr.Loc, 1, &v)
	// case Vec2i:
	// 	v := value.(mgl32.Vec2)
	// 	gl.Uniform2iv(attr.Loc, 1, &v)
	// case Vec3i:
	// 	v := value.(int32)
	// 	gl.Uniform1iv(attr.Loc, 1, &v)
	// case Vec4i:
	// 	v := value.(int32)
	// 	gl.Uniform1iv(attr.Loc, 1, &v)
	case UInt:
		v := value.(uint32)
		gl.Uniform1uiv(attr.Loc, 1, &v)
		// case Vec2ui:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Vec3ui:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Vec4ui:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Double:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Vec2d:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Vec3d:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Vec4d:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Mat2d:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Mat3d:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Mat4d:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Mat2x3d:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Mat2x4d:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Mat3x2d:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Mat3x4d:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Mat4x2d:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
		// case Mat4x3d:
		// 	v := value.(int32)
		// 	gl.Uniform1iv(attr.Loc, 1, &v)
	default:
		return fmt.Errorf("invalid attribute type: %v", attr.Type)
	}

	return nil
}

// deleteProgram deletes the program.
func (p *Program) Delete() {
	gl.DeleteProgram(p.ID)
}

package bgl

import (
	"embed"
	"fmt"
	"io/fs"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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

func init() {
	emb, _ := fs.Sub(embedded, "src")

	frag, _ := fs.ReadFile(emb, "default.frag")
	DefaultFrag = string(frag) + "\x00"

	vert, _ := fs.ReadFile(emb, "default.vert")
	DefaultVert = string(vert) + "\x00"
}

func DefaultProgram() (*Program, error) {
	return NewProgram([]*Shader{
		NewShader(DefaultFrag, FragShader),
		NewShader(DefaultVert, VertShader),
	})
}

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

type Shader struct {
	stype    ShaderType
	id       uint32
	src      string
	compiled bool
}

func NewShader(src string, stype ShaderType) *Shader {
	s := &Shader{
		src:   src,
		stype: stype,
	}

	return s
}

func (s *Shader) Type() ShaderType {
	return s.stype
}

func (s *Shader) Id() uint32 {
	return s.id
}

func (s *Shader) Src() string {
	return s.src
}

func (s *Shader) Compiled() bool {
	return s.compiled
}

func (s *Shader) getiv(iv uint32) int32 {
	var result int32
	gl.GetShaderiv(s.Id(), iv, &result)

	return result
}

func (s *Shader) getInfoLog() string {
	logLength := s.getiv(gl.INFO_LOG_LENGTH)

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(s.Id(), logLength, nil, gl.Str(log))

	return log
}

func (s *Shader) compile() error {
	s.id = gl.CreateShader(uint32(s.stype))

	csources, free := gl.Strs(s.src)
	gl.ShaderSource(s.Id(), 1, csources, nil)
	gl.CompileShader(s.Id())
	free()

	status := s.getiv(gl.COMPILE_STATUS)
	if status == gl.FALSE {
		log := s.getInfoLog()

		return fmt.Errorf("failed to compile shader (type: %v):\n%v", s.stype, log)
	}

	s.compiled = true

	return nil
}

func (s *Shader) delete() {
	gl.DeleteShader(s.Id())
}

type Program struct {
	binder

	VertexAttrs  AttrFormat // vertex attribute format
	UniformAttrs AttrFormat // uniform attribute format

	compiled bool
}

func NewProgram(shaders []*Shader) (*Program, error) {
	p := &Program{
		binder: binder{
			binding: CurrentProgram,
			bindFunc: func(id uint32) {
				gl.UseProgram(id)
			},
		},
		VertexAttrs:  make(AttrFormat),
		UniformAttrs: make(AttrFormat),
	}

	if len(shaders) == 0 {
		return nil, fmt.Errorf("no shaders defined")
	}

	p.id = gl.CreateProgram()

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

	runtime.SetFinalizer(p, (*Program).delete)

	return p, nil
}

func (p *Program) Compiled() bool {
	return p.compiled
}

func (p *Program) attach(shader *Shader) {
	gl.AttachShader(p.id, uint32(shader.id))
}

func (p *Program) link() {
	gl.LinkProgram(p.id)
}

func (p *Program) getiv(iv uint32) int32 {
	var result int32
	gl.GetProgramiv(p.id, iv, &result)

	return result
}

func (p *Program) GetInfoLog() string {
	logLength := p.getiv(gl.INFO_LOG_LENGTH)

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(p.id, logLength, nil, gl.Str(log))

	return log
}

func (p *Program) findUniforms() {
	p.Bind()
	defer p.Unbind()

	active := uint32(p.getiv(gl.ACTIVE_UNIFORMS))
	for i := uint32(0); i < active; i++ {
		a := Attr{}
		var buf [256]byte

		gl.GetActiveUniform(p.id, i, 256, &a.Len, &a.Size, &a.Type, &buf[0])

		a.Loc = gl.GetUniformLocation(p.id, &buf[0])
		a.Name = gl.GoStr(&buf[0])

		p.UniformAttrs[a.Name] = a
	}
}

func (p *Program) findAttributes() {
	p.Bind()
	defer p.Unbind()

	active := uint32(p.getiv(gl.ACTIVE_ATTRIBUTES))
	for i := uint32(0); i < active; i++ {
		a := Attr{}
		var buf [256]byte
		gl.GetActiveAttrib(p.id, i, 256, &a.Len, &a.Size, &a.Type, &buf[0])

		a.Loc = gl.GetAttribLocation(p.id, &buf[0])
		a.Name = gl.GoStr(&buf[0])

		p.VertexAttrs[a.Name] = a
	}
}

func (p *Program) Bind() {
	p.binder.bind()
}

func (p *Program) Unbind() {
	p.binder.restore()
}

func (p *Program) SetUniformAttr(name string, value interface{}) error {
	var attr Attr

	attr, ok := p.VertexAttrs[name]
	if !ok {
		attr, ok = p.UniformAttrs[name]
		if !ok {
			return fmt.Errorf("attribute not found: %v", name)
		}
	}

	switch AttrType(attr.Type) {
	case Float:
		v := value.(float32)
		gl.Uniform1fv(attr.Loc, 1, &v)
	case Vec2f:
		v := value.(mgl32.Vec2)
		gl.Uniform2fv(attr.Loc, 1, &v[0])
	case Vec3f:
		v := value.(mgl32.Vec3)
		gl.Uniform3fv(attr.Loc, 1, &v[0])
	case Vec4f:
		v := value.(mgl32.Vec4)
		gl.Uniform4fv(attr.Loc, 1, &v[0])
	case Mat2f:
		v := value.(mgl32.Mat2)
		gl.UniformMatrix2fv(attr.Loc, 1, false, &v[0])
	case Mat3f:
		v := value.(mgl32.Mat3)
		gl.UniformMatrix3fv(attr.Loc, 1, false, &v[0])
	case Mat4f:
		v := value.(mgl32.Mat4)
		gl.UniformMatrix4fv(attr.Loc, 1, false, &v[0])
	case Mat2x3f:
		v := value.(mgl32.Mat2x3)
		gl.UniformMatrix2x3fv(attr.Loc, 1, false, &v[0])
	case Mat2x4f:
		v := value.(mgl32.Mat2x4)
		gl.UniformMatrix2x4fv(attr.Loc, 1, false, &v[0])
	case Mat3x2f:
		v := value.(mgl32.Mat3x2)
		gl.UniformMatrix3x2fv(attr.Loc, 1, false, &v[0])
	case Mat3x4f:
		v := value.(mgl32.Mat3x4)
		gl.UniformMatrix3x4fv(attr.Loc, 1, false, &v[0])
	case Mat4x2f:
		v := value.(mgl32.Mat4x2)
		gl.UniformMatrix4x2fv(attr.Loc, 1, false, &v[0])
	case Mat4x3f:
		v := value.(mgl32.Mat4x3)
		gl.UniformMatrix4x3fv(attr.Loc, 1, false, &v[0])
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

func (p *Program) delete() {
	gl.DeleteProgram(p.id)
}

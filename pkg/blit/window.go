package blit

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/octalide/blit/pkg/bgl"
	"github.com/octalide/wisp/pkg/wisp"
)

type WindowOptions struct {
	Title      string
	Resizable  bool
	Fullscreen bool
	Decorated  bool
	MSAA       bool
	VSync      bool
	Width      int
	Height     int
}

func DefaultWindowOptions() *WindowOptions {
	wo := &WindowOptions{}

	wo.Title = "Default"
	wo.Resizable = true
	wo.Fullscreen = false
	wo.Decorated = true
	wo.MSAA = true
	wo.VSync = true
	wo.Width = 512
	wo.Height = 512

	return wo
}

type Window struct {
	win     *glfw.Window
	options *WindowOptions

	focused bool
}

func NewWindow(options *WindowOptions) *Window {
	if options == nil {
		options = DefaultWindowOptions()
	}

	w := &Window{
		options: options,
	}

	return w
}

func (w *Window) Init() error {
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("glfw failed to initialize: %w", err)
	}

	mode := glfw.GetPrimaryMonitor().GetVideoMode()

	glfw.WindowHint(glfw.ContextVersionMajor, OpenGLVersionMajor)
	glfw.WindowHint(glfw.ContextVersionMinor, OpenGLVersionMinor)
	glfw.WindowHint(glfw.OpenGLProfile, OpenGLProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.False)

	glfw.WindowHint(glfw.RedBits, mode.RedBits)
	glfw.WindowHint(glfw.GreenBits, mode.GreenBits)
	glfw.WindowHint(glfw.BlueBits, mode.BlueBits)
	glfw.WindowHint(glfw.RefreshRate, mode.RefreshRate)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)

	glfw.WindowHint(glfw.Resizable, glfwBool(w.options.Resizable))
	glfw.WindowHint(glfw.Decorated, glfwBool(w.options.Decorated))

	if w.options.MSAA {
		glfw.WindowHint(glfw.Samples, 4)
	}

	var monitor *glfw.Monitor
	if w.options.Fullscreen {
		monitor = glfw.GetPrimaryMonitor()
	}

	win, err := glfw.CreateWindow(w.options.Width, w.options.Height, w.options.Title, monitor, nil)
	if err != nil {
		return fmt.Errorf("failed to create window: %w", err)
	}

	w.win = win
	w.BindCallbacks()

	w.win.MakeContextCurrent()

	w.SetVSync(w.options.VSync)

	w.Show()

	runtime.SetFinalizer(w, (*Window).Destroy)

	return nil
}

func (w *Window) BindCallbacks() {
	w.win.SetSizeCallback(w.onResize)
	w.win.SetFocusCallback(w.onFocus)
	w.win.SetKeyCallback(w.onKey)
	w.win.SetCharCallback(w.onChar)
	w.win.SetMouseButtonCallback(w.onMouseButton)
	w.win.SetScrollCallback(w.onScroll)
	w.win.SetCursorPosCallback(w.onCursorMove)
}

func (w *Window) Options() WindowOptions {
	return *w.options
}

func (w *Window) SetTitle(title string) {
	w.options.Title = title

	w.win.SetTitle(w.options.Title)
}

func (w *Window) GetResizable() bool {
	return w.options.Resizable
}

func (w *Window) SetResizable(resizable bool) {
	w.options.Resizable = resizable

	glfw.WindowHint(glfw.Resizable, glfwBool(w.options.Resizable))
}

func (w *Window) GetFullscreen() bool {
	return w.options.Fullscreen
}

func (w *Window) SetFullscreen(fullscreen bool) {
	w.options.Fullscreen = fullscreen

	mon := glfw.GetPrimaryMonitor()
	mode := mon.GetVideoMode()
	x := 0
	y := 0

	if w.options.Fullscreen {
		w.options.Width = mode.Width
		w.options.Height = mode.Height
	} else {
		mon = nil
		w.options.Width = DefaultWindowOptions().Width
		w.options.Height = DefaultWindowOptions().Height
		x = w.options.Width / 2
		y = w.options.Height / 2
	}

	w.win.SetMonitor(mon, x, y, w.options.Width, w.options.Height, mode.RefreshRate)
	bgl.SetBounds(0, 0, w.options.Width, w.options.Height)
}

func (w *Window) GetMSAA() bool {
	return w.options.MSAA
}

func (w *Window) GetVSync() bool {
	return w.options.VSync
}

func (w *Window) GetWidth() int {
	return w.options.Width
}

func (w *Window) GetHeight() int {
	return w.options.Height
}

func (w *Window) GetAspectRatio() float32 {
	return float32(w.GetWidth()) / float32(w.GetHeight())
}

func (w *Window) SetMSAA(msaa bool) {
	w.options.MSAA = msaa

	if w.options.MSAA {
		bgl.EnableMSAA()
	} else {
		bgl.DisableMSAA()
	}
}

func (w *Window) SetWidth(width int) {
	w.win.SetSize(width, w.options.Height)
}

func (w *Window) SetHeight(height int) {
	w.win.SetSize(w.options.Width, height)
}

func (w *Window) SetVSync(vsync bool) {
	w.options.VSync = vsync

	if w.options.VSync {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}
}

func (w *Window) LockAspectRatio(numer, denom int) {
	w.win.SetAspectRatio(numer, denom)
}

func (w *Window) ShowCursor() {
	w.win.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
}

func (w *Window) HideCursor() {
	w.win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

func (w *Window) Show() {
	w.win.Show()
}

func (w *Window) Hide() {
	w.win.Hide()
}

func (w *Window) Focus() {
	w.win.Focus()
}

func (w *Window) Focused() bool {
	return w.focused
}

func (w *Window) Close() {
	w.win.SetShouldClose(true)
}

func (w *Window) ShouldClose() bool {
	return w.win.ShouldClose()
}

func (w *Window) SwapBuffers() {
	w.win.SwapBuffers()
}

func (w *Window) Destroy() {
	w.win.Destroy()
}

func (w *Window) SetSize(width, height int) {
	w.options.Width = width
	w.options.Height = height

	w.win.SetSize(width, height)
}

var (
	eventResize          = wisp.NewEvent("core.window.resize", nil)
	eventFocus           = wisp.NewEvent("core.window.focus", nil)
	eventChar            = wisp.NewEvent("core.input.char", nil)
	eventKeyDown         = wisp.NewEvent("core.input.key.down", nil)
	eventKeyUp           = wisp.NewEvent("core.input.key.up", nil)
	eventKeyRepeat       = wisp.NewEvent("core.input.key.repeat", nil)
	eventMouseButtonDown = wisp.NewEvent("core.input.mouse.button.down", nil)
	eventMouseButtonUp   = wisp.NewEvent("core.input.mouse.button.up", nil)
	eventMouseMove       = wisp.NewEvent("core.input.mouse.move", nil)
	eventMouseScroll     = wisp.NewEvent("core.input.mouse.scroll", nil)
)

func (w *Window) onResize(_ *glfw.Window, width, height int) {
	w.options.Width = width
	w.options.Height = height

	bgl.SetBounds(0, 0, w.options.Width, w.options.Height)

	eventResize.Data = Vec{float32(width), float32(height)}
	wisp.Broadcast(eventResize)
}

func (w *Window) onFocus(_ *glfw.Window, focused bool) {
	w.focused = focused

	eventFocus.Data = focused
	wisp.Broadcast(eventFocus)
}

func (w *Window) onChar(_ *glfw.Window, char rune) {
	eventChar.Data = char
	wisp.Broadcast(eventChar)
}

func (w *Window) onKey(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		eventKeyDown.Data = Key(key)
		wisp.Broadcast(eventKeyDown)
	} else if action == glfw.Release {
		eventKeyUp.Data = Key(key)
		wisp.Broadcast(eventKeyUp)
	} else if action == glfw.Repeat {
		eventKeyRepeat.Data = Key(key)
		wisp.Broadcast(eventKeyRepeat)
	}
}

func (w *Window) onMouseButton(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		eventMouseButtonDown.Data = Key(button)
		wisp.Broadcast(eventMouseButtonDown)
	} else if action == glfw.Release {
		eventMouseButtonUp.Data = Key(button)
		wisp.Broadcast(eventMouseButtonUp)
	}
}

func (w *Window) onScroll(_ *glfw.Window, xoff float64, yoff float64) {
	eventMouseScroll.Data = Vec{float32(xoff), float32(yoff)}
	wisp.Broadcast(eventMouseScroll)
}

func (w *Window) onCursorMove(_ *glfw.Window, xpos float64, ypos float64) {
	eventMouseMove.Data = Vec{float32(xpos), float32(ypos)}
	wisp.Broadcast(eventMouseMove)
}

func glfwBool(b bool) int {
	if b {
		return glfw.True
	}

	return glfw.False
}

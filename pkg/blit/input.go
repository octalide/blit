package blit

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/octalide/wisp/pkg/wisp"
)

var (
	mouse struct {
		Cursor, Scroll, Delta Vec
	}

	buttons map[Key]bool
)

func init() {
	buttons = make(map[Key]bool)

	// wisp setup
	wisp.Init()
	wisp.AddHandler(&wisp.Handler{
		Callback: func(e *wisp.Event) bool {
			switch e.Tag {
			case "core.input.mouse.move":
				pos := e.Data.(Vec)
				mouse.Delta = pos.Sub(mouse.Cursor)
				mouse.Cursor = pos
			case "core.input.mouse.button.up":
				buttons[Key(e.Data.(Key))] = false
			case "core.input.mouse.button.down":
				buttons[Key(e.Data.(Key))] = true
			case "core.input.mouse.scroll":
				mouse.Scroll = e.Data.(Vec)
			case "core.input.key.up":
				buttons[Key(e.Data.(Key))] = false
			case "core.input.key.down":
				buttons[Key(e.Data.(Key))] = true
			case "core.input.key.repeat":
				buttons[Key(e.Data.(Key))] = true
			}

			return false
		},
		Tags:     []string{"core.input"},
		Blocking: true,
	})
}

func Keys(key Key) bool {
	v, ok := buttons[key]
	if !ok {
		buttons[key] = false
		return false
	}

	return v
}

func MouseDelta() Vec {
	return mouse.Delta
}

func MousePos() Vec {
	return mouse.Cursor
}

func MouseScroll() Vec {
	return mouse.Scroll
}

type Key int

const (
	MouseButton1      = Key(glfw.MouseButton1)
	MouseButton2      = Key(glfw.MouseButton2)
	MouseButton3      = Key(glfw.MouseButton3)
	MouseButton4      = Key(glfw.MouseButton4)
	MouseButton5      = Key(glfw.MouseButton5)
	MouseButton6      = Key(glfw.MouseButton6)
	MouseButton7      = Key(glfw.MouseButton7)
	MouseButton8      = Key(glfw.MouseButton8)
	MouseButtonLast   = Key(glfw.MouseButtonLast)
	MouseButtonLeft   = Key(glfw.MouseButtonLeft)
	MouseButtonRight  = Key(glfw.MouseButtonRight)
	MouseButtonMiddle = Key(glfw.MouseButtonMiddle)

	KeyUnknown      = Key(glfw.KeyUnknown)
	KeySpace        = Key(glfw.KeySpace)
	KeyApostrophe   = Key(glfw.KeyApostrophe)
	KeyComma        = Key(glfw.KeyComma)
	KeyMinus        = Key(glfw.KeyMinus)
	KeyPeriod       = Key(glfw.KeyPeriod)
	KeySlash        = Key(glfw.KeySlash)
	Key0            = Key(glfw.Key0)
	Key1            = Key(glfw.Key1)
	Key2            = Key(glfw.Key2)
	Key3            = Key(glfw.Key3)
	Key4            = Key(glfw.Key4)
	Key5            = Key(glfw.Key5)
	Key6            = Key(glfw.Key6)
	Key7            = Key(glfw.Key7)
	Key8            = Key(glfw.Key8)
	Key9            = Key(glfw.Key9)
	KeySemicolon    = Key(glfw.KeySemicolon)
	KeyEqual        = Key(glfw.KeyEqual)
	KeyA            = Key(glfw.KeyA)
	KeyB            = Key(glfw.KeyB)
	KeyC            = Key(glfw.KeyC)
	KeyD            = Key(glfw.KeyD)
	KeyE            = Key(glfw.KeyE)
	KeyF            = Key(glfw.KeyF)
	KeyG            = Key(glfw.KeyG)
	KeyH            = Key(glfw.KeyH)
	KeyI            = Key(glfw.KeyI)
	KeyJ            = Key(glfw.KeyJ)
	KeyK            = Key(glfw.KeyK)
	KeyL            = Key(glfw.KeyL)
	KeyM            = Key(glfw.KeyM)
	KeyN            = Key(glfw.KeyN)
	KeyO            = Key(glfw.KeyO)
	KeyP            = Key(glfw.KeyP)
	KeyQ            = Key(glfw.KeyQ)
	KeyR            = Key(glfw.KeyR)
	KeyS            = Key(glfw.KeyS)
	KeyT            = Key(glfw.KeyT)
	KeyU            = Key(glfw.KeyU)
	KeyV            = Key(glfw.KeyV)
	KeyW            = Key(glfw.KeyW)
	KeyX            = Key(glfw.KeyX)
	KeyY            = Key(glfw.KeyY)
	KeyZ            = Key(glfw.KeyZ)
	KeyLeftBracket  = Key(glfw.KeyLeftBracket)
	KeyBackslash    = Key(glfw.KeyBackslash)
	KeyRightBracket = Key(glfw.KeyRightBracket)
	KeyGraveAccent  = Key(glfw.KeyGraveAccent)
	KeyWorld1       = Key(glfw.KeyWorld1)
	KeyWorld2       = Key(glfw.KeyWorld2)
	KeyEscape       = Key(glfw.KeyEscape)
	KeyEnter        = Key(glfw.KeyEnter)
	KeyTab          = Key(glfw.KeyTab)
	KeyBackspace    = Key(glfw.KeyBackspace)
	KeyInsert       = Key(glfw.KeyInsert)
	KeyDelete       = Key(glfw.KeyDelete)
	KeyRight        = Key(glfw.KeyRight)
	KeyLeft         = Key(glfw.KeyLeft)
	KeyDown         = Key(glfw.KeyDown)
	KeyUp           = Key(glfw.KeyUp)
	KeyPageUp       = Key(glfw.KeyPageUp)
	KeyPageDown     = Key(glfw.KeyPageDown)
	KeyHome         = Key(glfw.KeyHome)
	KeyEnd          = Key(glfw.KeyEnd)
	KeyCapsLock     = Key(glfw.KeyCapsLock)
	KeyScrollLock   = Key(glfw.KeyScrollLock)
	KeyNumLock      = Key(glfw.KeyNumLock)
	KeyPrintScreen  = Key(glfw.KeyPrintScreen)
	KeyPause        = Key(glfw.KeyPause)
	KeyF1           = Key(glfw.KeyF1)
	KeyF2           = Key(glfw.KeyF2)
	KeyF3           = Key(glfw.KeyF3)
	KeyF4           = Key(glfw.KeyF4)
	KeyF5           = Key(glfw.KeyF5)
	KeyF6           = Key(glfw.KeyF6)
	KeyF7           = Key(glfw.KeyF7)
	KeyF8           = Key(glfw.KeyF8)
	KeyF9           = Key(glfw.KeyF9)
	KeyF10          = Key(glfw.KeyF10)
	KeyF11          = Key(glfw.KeyF11)
	KeyF12          = Key(glfw.KeyF12)
	KeyF13          = Key(glfw.KeyF13)
	KeyF14          = Key(glfw.KeyF14)
	KeyF15          = Key(glfw.KeyF15)
	KeyF16          = Key(glfw.KeyF16)
	KeyF17          = Key(glfw.KeyF17)
	KeyF18          = Key(glfw.KeyF18)
	KeyF19          = Key(glfw.KeyF19)
	KeyF20          = Key(glfw.KeyF20)
	KeyF21          = Key(glfw.KeyF21)
	KeyF22          = Key(glfw.KeyF22)
	KeyF23          = Key(glfw.KeyF23)
	KeyF24          = Key(glfw.KeyF24)
	KeyF25          = Key(glfw.KeyF25)
	KeyKP0          = Key(glfw.KeyKP0)
	KeyKP1          = Key(glfw.KeyKP1)
	KeyKP2          = Key(glfw.KeyKP2)
	KeyKP3          = Key(glfw.KeyKP3)
	KeyKP4          = Key(glfw.KeyKP4)
	KeyKP5          = Key(glfw.KeyKP5)
	KeyKP6          = Key(glfw.KeyKP6)
	KeyKP7          = Key(glfw.KeyKP7)
	KeyKP8          = Key(glfw.KeyKP8)
	KeyKP9          = Key(glfw.KeyKP9)
	KeyKPDecimal    = Key(glfw.KeyKPDecimal)
	KeyKPDivide     = Key(glfw.KeyKPDivide)
	KeyKPMultiply   = Key(glfw.KeyKPMultiply)
	KeyKPSubtract   = Key(glfw.KeyKPSubtract)
	KeyKPAdd        = Key(glfw.KeyKPAdd)
	KeyKPEnter      = Key(glfw.KeyKPEnter)
	KeyKPEqual      = Key(glfw.KeyKPEqual)
	KeyLeftShift    = Key(glfw.KeyLeftShift)
	KeyLeftControl  = Key(glfw.KeyLeftControl)
	KeyLeftAlt      = Key(glfw.KeyLeftAlt)
	KeyLeftSuper    = Key(glfw.KeyLeftSuper)
	KeyRightShift   = Key(glfw.KeyRightShift)
	KeyRightControl = Key(glfw.KeyRightControl)
	KeyRightAlt     = Key(glfw.KeyRightAlt)
	KeyRightSuper   = Key(glfw.KeyRightSuper)
	KeyMenu         = Key(glfw.KeyMenu)
	KeyLast         = Key(glfw.KeyLast)
)

// String returns a human-readable string describing the Button.
func (b Key) String() string {
	name, ok := buttonNames[b]
	if !ok {
		return "Invalid"
	}
	return name
}

var buttonNames = map[Key]string{
	MouseButton4:      "MouseButton4",
	MouseButton5:      "MouseButton5",
	MouseButton6:      "MouseButton6",
	MouseButton7:      "MouseButton7",
	MouseButton8:      "MouseButton8",
	MouseButtonLeft:   "MouseButtonLeft",
	MouseButtonRight:  "MouseButtonRight",
	MouseButtonMiddle: "MouseButtonMiddle",
	KeyUnknown:        "Unknown",
	KeySpace:          "Space",
	KeyApostrophe:     "Apostrophe",
	KeyComma:          "Comma",
	KeyMinus:          "Minus",
	KeyPeriod:         "Period",
	KeySlash:          "Slash",
	Key0:              "0",
	Key1:              "1",
	Key2:              "2",
	Key3:              "3",
	Key4:              "4",
	Key5:              "5",
	Key6:              "6",
	Key7:              "7",
	Key8:              "8",
	Key9:              "9",
	KeySemicolon:      "Semicolon",
	KeyEqual:          "Equal",
	KeyA:              "A",
	KeyB:              "B",
	KeyC:              "C",
	KeyD:              "D",
	KeyE:              "E",
	KeyF:              "F",
	KeyG:              "G",
	KeyH:              "H",
	KeyI:              "I",
	KeyJ:              "J",
	KeyK:              "K",
	KeyL:              "L",
	KeyM:              "M",
	KeyN:              "N",
	KeyO:              "O",
	KeyP:              "P",
	KeyQ:              "Q",
	KeyR:              "R",
	KeyS:              "S",
	KeyT:              "T",
	KeyU:              "U",
	KeyV:              "V",
	KeyW:              "W",
	KeyX:              "X",
	KeyY:              "Y",
	KeyZ:              "Z",
	KeyLeftBracket:    "LeftBracket",
	KeyBackslash:      "Backslash",
	KeyRightBracket:   "RightBracket",
	KeyGraveAccent:    "GraveAccent",
	KeyWorld1:         "World1",
	KeyWorld2:         "World2",
	KeyEscape:         "Escape",
	KeyEnter:          "Enter",
	KeyTab:            "Tab",
	KeyBackspace:      "Backspace",
	KeyInsert:         "Insert",
	KeyDelete:         "Delete",
	KeyRight:          "Right",
	KeyLeft:           "Left",
	KeyDown:           "Down",
	KeyUp:             "Up",
	KeyPageUp:         "PageUp",
	KeyPageDown:       "PageDown",
	KeyHome:           "Home",
	KeyEnd:            "End",
	KeyCapsLock:       "CapsLock",
	KeyScrollLock:     "ScrollLock",
	KeyNumLock:        "NumLock",
	KeyPrintScreen:    "PrintScreen",
	KeyPause:          "Pause",
	KeyF1:             "F1",
	KeyF2:             "F2",
	KeyF3:             "F3",
	KeyF4:             "F4",
	KeyF5:             "F5",
	KeyF6:             "F6",
	KeyF7:             "F7",
	KeyF8:             "F8",
	KeyF9:             "F9",
	KeyF10:            "F10",
	KeyF11:            "F11",
	KeyF12:            "F12",
	KeyF13:            "F13",
	KeyF14:            "F14",
	KeyF15:            "F15",
	KeyF16:            "F16",
	KeyF17:            "F17",
	KeyF18:            "F18",
	KeyF19:            "F19",
	KeyF20:            "F20",
	KeyF21:            "F21",
	KeyF22:            "F22",
	KeyF23:            "F23",
	KeyF24:            "F24",
	KeyF25:            "F25",
	KeyKP0:            "KP0",
	KeyKP1:            "KP1",
	KeyKP2:            "KP2",
	KeyKP3:            "KP3",
	KeyKP4:            "KP4",
	KeyKP5:            "KP5",
	KeyKP6:            "KP6",
	KeyKP7:            "KP7",
	KeyKP8:            "KP8",
	KeyKP9:            "KP9",
	KeyKPDecimal:      "KPDecimal",
	KeyKPDivide:       "KPDivide",
	KeyKPMultiply:     "KPMultiply",
	KeyKPSubtract:     "KPSubtract",
	KeyKPAdd:          "KPAdd",
	KeyKPEnter:        "KPEnter",
	KeyKPEqual:        "KPEqual",
	KeyLeftShift:      "LeftShift",
	KeyLeftControl:    "LeftControl",
	KeyLeftAlt:        "LeftAlt",
	KeyLeftSuper:      "LeftSuper",
	KeyRightShift:     "RightShift",
	KeyRightControl:   "RightControl",
	KeyRightAlt:       "RightAlt",
	KeyRightSuper:     "RightSuper",
	KeyMenu:           "Menu",
}

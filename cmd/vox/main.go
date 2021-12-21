package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
	"time"

	"github.com/octalide/blit/pkg/bgl"
	"github.com/octalide/blit/pkg/blit"
	"github.com/octalide/wisp/pkg/wisp"
)

func main() {
	runtime.LockOSThread()

	log.Println("creating window")
	opt := blit.DefaultWindowOptions()
	opt.VSync = false
	opt.MSAA = false
	win := blit.NewWindow(opt)

	log.Println("initializing window")
	err := win.Init()
	if err != nil {
		panic(err)
	}

	log.Println("initializing blit")
	err = blit.Init()
	if err != nil {
		panic(err)
	}

	bgl.SetClearColor(color.RGBA{0, 0, 255, 255})
	bgl.SetBounds(0, 0, win.GetWidth(), win.GetHeight())

	log.Println("creating shader...")
	shader, err := bgl.DefaultProgram()
	if err != nil {
		panic(err)
	}

	var minFOV, maxFOV, zoomFactor float32
	minFOV = 170
	maxFOV = 179.5
	zoomFactor = 0.4

	cam := blit.NewCam()
	cam.FOV = minFOV + (maxFOV-minFOV)/2
	cam.Viewport = blit.Viewport()

	// camera controller
	wisp.AddHandler(&wisp.Handler{
		Callback: func(e *wisp.Event) bool {
			switch e.Tag {
			case "core.input.mouse.scroll":
				delta := e.Data.(blit.Vec).Y()
				cam.FOV -= delta * zoomFactor
				if cam.FOV < minFOV {
					cam.FOV = minFOV
				}
				if cam.FOV > maxFOV {
					cam.FOV = maxFOV
				}
			case "core.input.mouse.move":
				// change camera position by mouse delta if middle mouse button is down
				delta := blit.MouseDelta()
				// invert Y axis
				delta[1] *= -1

				// change delta to world coordinates
				// delta = cam.Unproject(delta)
				pan := delta.Scl(0.01 * (cam.FOV - minFOV))

				if blit.Keys(blit.MouseButtonMiddle) {
					cam.Pan(pan.Inv())
				}

				fmt.Printf("%v : %v %v %v\r", cam.FOV, cam.Vec, delta, pan)
			}
			// log.Printf("event (%v): %v", e.Tag, e.Data)
			return false
		},
		Tags:     []string{"*"},
		Blocking: false,
	})

	log.Println("entering main loop...")

	lastPrint := time.Now()
	last := time.Now()
	for !win.ShouldClose() {
		time.Since(last) // delta
		last = time.Now()

		cam.Use(shader)

		bgl.Clear()

		blit.Update()

		if time.Since(lastPrint) > time.Second {
			// fmt.Printf("fps: %v\r", float64(1/delta.Seconds()))
			lastPrint = time.Now()
		}
	}
}

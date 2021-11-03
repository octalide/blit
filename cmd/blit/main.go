package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
	"time"

	"github.com/octalide/blit/pkg/bgl"
	"github.com/octalide/blit/pkg/blit"
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

	// load spritesheet
	log.Println("loading spritesheet...")
	ss, err := blit.LoadSpritesheet("tiles/ground")
	if err != nil {
		panic(err)
	}

	log.Println("creating shader...")
	shader, err := bgl.DefaultProgram()
	if err != nil {
		panic(err)
	}

	log.Println("creating sprites...")
	dirt, err := ss.Get("dirt", shader)
	if err != nil {
		panic(err)
	}

	grass, err := ss.Get("grass", shader)
	if err != nil {
		panic(err)
	}

	stone, err := ss.Get("stone", shader)
	if err != nil {
		panic(err)
	}

	cam := blit.NewCam()

	log.Println("entering main loop...")

	lastPrint := time.Now()
	last := time.Now()
	for !win.ShouldClose() {
		cam.Use(shader)

		bgl.Clear()
		dirt.Draw()
		grass.Draw()
		stone.Draw()

		blit.Update()

		// print fps calculated from last frame
		delta := time.Since(last)
		last = time.Now()

		// rotate modl by 0.1 radians per second according to delta
		dirt.O.R += float32(0.2 * delta.Seconds())
		grass.O.X += float32(0.2 * delta.Seconds())
		stone.O.Y -= float32(0.2 * delta.Seconds())
		// fmt.Printf("%.3f : %v\r", dirt.O.R, dirt.O.Mat())

		// print mouse position
		fmt.Printf("%v\r", blit.MousePos())

		if time.Since(lastPrint) > time.Second {
			// fmt.Printf("fps: %v\r", float64(1/delta.Seconds()))
			lastPrint = time.Now()
		}
	}
}

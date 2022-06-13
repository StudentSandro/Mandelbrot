// Copyright 2017 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build example
// +build example

package main

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"math"
)

const (
	screenWidth  = 640
	screenHeight = 640
	maxIt        = 128
	MaxGoRoutine = 10
)

var (
	palette [maxIt]byte
)

//einzelnen Pixel mit x, y und farbtiefe
type Pixel struct {
	X  int
	Y  int
	It int
}

//Channel für Pixel
var PixelChan = make(chan Pixel)

//Offset Mandelbrot (steuert Position und Zoom)
var OffsetX float64 = -0.75
var OffsetY float64 = 0.25
var Zoom float64 = 3

func init() {
	for i := range palette {
		palette[i] = byte(math.Sqrt(float64(i)/float64(len(palette))) * 0x80)
	}
}

func color(it int) (r, g, b byte) {
	if it == maxIt {
		return 0x00, 0xaf, 0xff //Standartfarbe ab maxIt Wert
	}
	c := palette[it]
	return 0, c / 2, c //färbt Mandelbrot blau ein
}

type Game struct {
	offscreen    *ebiten.Image
	offscreenPix []byte
}

func NewGame() *Game {
	g := &Game{
		offscreen:    ebiten.NewImage(screenWidth, screenHeight),
		offscreenPix: make([]byte, screenWidth*screenHeight*4),
	}
	//führt bei Start erstes mal Mandelbrot Berechnung durch
	g.updateOffscreen()

	return g
}

func (gm *Game) updateOffscreen() {
	ChanCounter := screenWidth * screenHeight //damit bekannt ist wie viele Pixel im Channel erwartet werden
	for i := 0; i < MaxGoRoutine; i++ {
		for j := 0; j < MaxGoRoutine; j++ {
			xStart := (screenWidth / MaxGoRoutine) * i
			xEnd := (screenWidth / MaxGoRoutine) * (i + 1)
			yStart := (screenHeight / MaxGoRoutine) * j
			yEnd := (screenHeight / MaxGoRoutine) * (j + 1)
			go CalcBlock(xStart, xEnd, yStart, yEnd) //startet funktion um Pixel Blöcke zu berrechen, rückgabe per channel
		}
	}

	for c := ChanCounter; c > 0; c-- { //holt Pixel aus channel und fügt sie dem offscreen zu
		Pixel := <-PixelChan
		r, g, b := color(Pixel.It)
		p := 4 * (int(Pixel.X) + int(Pixel.Y)*screenWidth)
		gm.offscreenPix[p] = r
		gm.offscreenPix[p+1] = g
		gm.offscreenPix[p+2] = b
		gm.offscreenPix[p+3] = 0xff
	}
	gm.offscreen.ReplacePixels(gm.offscreenPix) //Zeichnet anschliessend das Mandelbrot neu anhand der offscreenPixel

}

func (g *Game) Update() error {
	// Zoom hinein
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		Zoom *= 0.9
		g.updateOffscreen()
	}
	//Zoom hinaus
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		Zoom /= 0.9
		g.updateOffscreen()
	}
	//Kamera geht hoch
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		OffsetY += 0.05 * Zoom
		g.updateOffscreen()
	}
	//Kamera geht runter
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		OffsetY -= 0.05 * Zoom
		g.updateOffscreen()
	}
	//Kamera geht nach links
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		OffsetX -= 0.05 * Zoom
		g.updateOffscreen()
	}
	//Kamera geht nach rechts
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		OffsetX += 0.05 * Zoom
		g.updateOffscreen()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.offscreen, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

//Berrechnet Teil Blöcke des Mandelbrots
func CalcBlock(xStart int, xEnd int, yStart int, yEnd int) {
	for y := yStart; y < yEnd; y++ {
		for x := xStart; x < xEnd; x++ {
			it := CalcPixel(x, y)
			PixelChan <- Pixel{X: x, Y: y, It: it}
		}
	}
}

//Funktion zum berechnen der einzelnen Pixel mit der Mandelbrotformel
func CalcPixel(i, j int) int {
	x := float64(i)*Zoom/screenWidth - Zoom/2 + OffsetX
	y := (screenHeight-float64(j))*Zoom/screenHeight - Zoom/2 + OffsetY
	c := complex(x, y)
	z := complex(0, 0)
	it := 0
	for ; it < maxIt; it++ {
		z = z*z + c
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			break
		}
	}
	return it
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mandelbrot by Patrick Rizzo / Sandro Zogg")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

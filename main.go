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
	"time"
)

const (
	screenWidth  = 640
	screenHeight = 640
	maxIt        = 128
)

var (
	palette [maxIt]byte
)

var PixelList []Pixel

type Pixel struct {
	X  float64
	Y  float64
	It int
}

//Offset
var OffsetX float64 = -0.75
var OffsetY float64 = 0.25
var Zoom float64 = 2

func init() {
	for i := range palette {
		palette[i] = byte(math.Sqrt(float64(i)/float64(len(palette))) * 0x80)
	}
}

func color(it int) (r, g, b byte) {
	if it == maxIt {
		return 0xff, 0xff, 0xff
	}
	c := palette[it]
	return c, c, c
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
	// Now it is not feasible to call updateOffscreen every frame due to performance.
	g.updateOffscreen() //Todo start update

	return g
}

func (gm *Game) updateOffscreen() {
	//TODO Funktion auslagern zum Bloecke bilden

	for i := 0; i <= 9; i++ {
		for j := 0; j <= 9; j++ {
			xStart := (screenWidth / 10) * i
			xEnd := (screenWidth / 10) * (i + 1)
			yStart := (screenHeight / 10) * j
			yEnd := (screenHeight / 10) * (j + 1)
			go CalcBlock(xStart, xEnd, yStart, yEnd)
			print(xStart, ",", xEnd, ",", yStart, ",", yEnd, "\n")
		}
	}

	time.Sleep(time.Second * 3)
	print(len(PixelList))
	for k := 0; k < len(PixelList); k++ {

		r, g, b := color(PixelList[k].It)
		p := 4 * (int(PixelList[k].X) + int(PixelList[k].Y)*screenWidth)
		gm.offscreenPix[p] = r
		gm.offscreenPix[p+1] = g
		gm.offscreenPix[p+2] = b
		gm.offscreenPix[p+3] = 0xff
	}
	gm.offscreen.ReplacePixels(gm.offscreenPix)
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.offscreen, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func CalcBlock(xStart int, xEnd int, yStart int, yEnd int) {
	//var PixelBlock []Pixel
	for y := yStart; y < yEnd; y++ {
		for x := xStart; x < xEnd; x++ {
			it := CalcPixel(x, y)
			PixelList = append(PixelList, Pixel{X: float64(x), Y: float64(y), It: it})
		}
	}
}

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
	ebiten.SetWindowTitle("Mandelbrot (Ebiten Demo)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

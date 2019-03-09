package main

// Copyright 2018 The Ebiten Authors
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

// build example jsgo

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"

	"./images"
	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 240
	screenHeight = 240

	frameOX     = 0
	frameWidth  = 32
	frameHeight = 32

	tileSize = 16
	tileXNum = 25
)

var (
	layers = [][]int{
		{
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 244, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 219, 243, 243, 243, 219, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243,
			243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
		},
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 26, 27, 28, 29, 30, 31, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 51, 52, 53, 54, 55, 56, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 76, 77, 78, 79, 80, 81, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 101, 102, 103, 104, 105, 106, 0, 0, 0, 0,

			0, 0, 0, 0, 0, 126, 127, 128, 129, 130, 131, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 301, 302, 245, 242, 303, 303, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,

			0, 0, 45, 46, 47, 48, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 70, 71, 72, 73, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 95, 96, 97, 98, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 120, 121, 122, 123, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			0, 0, 145, 146, 147, 148, 0, 245, 242, 0, 0, 0, 0, 0, 0,
		},
	}

	direction = "right"
	frameOY   = 32
	frameNum  = 8
)

var (
	count        = 0
	runnerImage  *ebiten.Image
	x            = 0
	y            = 0
	tilesImage   *ebiten.Image
	runnermImage *ebiten.Image
)

func drawTile(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Draw each tile with each DrawImage call.
	// As the source images of all DrawImage calls are always same,
	// this rendering is done very effectively.
	// For more detail, see https://godoc.org/github.com/hajimehoshi/ebiten#Image.DrawImage
	const xNum = screenWidth / tileSize
	for _, l := range layers {
		for i, t := range l {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64((i%xNum)*tileSize), float64((i/xNum)*tileSize))

			sx := (t % tileXNum) * tileSize
			sy := (t / tileXNum) * tileSize
			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}
	return nil
}

func drawPeople(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	status := "stand"
	count++
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			if k.String() == "Up" {
				status = "run"
				y--
			} else if k.String() == "Down" {
				status = "run"
				y++
			} else if k.String() == "Left" {
				status = "run"
				direction = "left"
				x--
			} else if k.String() == "Right" {
				status = "run"
				direction = "right"
				x++
			} else {
				fmt.Println("Not Support Key Pressing")
			}
			if x > screenWidth-32 {
				x = screenWidth - 32
			}
			if x < 0 {
				x = 0
			}
			if y > screenHeight-32 {
				y = screenHeight - 32
			}
			if y < 0 {
				y = 0
			}
			//fmt.Println(fmt.Sprintf("x: %d, y %d", x, y))
		}
	}

	if status == "stand" {
		frameOY = 0
		frameNum = 5
	} else {
		frameOY = 32
		frameNum = 8
	}

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(x), float64(y))

	i := (count / 5) % frameNum
	sx, sy := frameOX+i*frameWidth, frameOY

	if direction == "right" {
		screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
	} else {
		screen.DrawImage(runnermImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
	}

	return nil
}

func update(screen *ebiten.Image) error {
	drawTile(screen)
	drawPeople(screen)
	return nil
}

func main() {
	// Decode image from a byte slice instead of a file so that
	// this example works in any working directory.
	// If you want to use a file, there are some options:
	// 1) Use os.Open and pass the file to the image decoder.
	//    This is a very regular way, but doesn't work on browsers.
	// 2) Use ebitenutil.OpenFile and pass the file to the image decoder.
	//    This works even on browsers.
	// 3) Use ebitenutil.NewImageFromFile to create an ebiten.Image directly from a file.
	//    This also works on browsers.
	rimg, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	rmimg, _, err := image.Decode(bytes.NewReader(images.Runnerm_png))
	if err != nil {
		log.Fatal(err)
	}
	timg, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage, _ = ebiten.NewImageFromImage(timg, ebiten.FilterDefault)
	runnerImage, _ = ebiten.NewImageFromImage(rimg, ebiten.FilterDefault)
	runnermImage, _ = ebiten.NewImageFromImage(rmimg, ebiten.FilterDefault)

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "haha test"); err != nil {
		log.Fatal(err)
	}
}

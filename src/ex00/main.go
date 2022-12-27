package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type MyLogo struct {
	width, height int
	img           *image.RGBA
}

// func generateImage(width int, height int, pixelColor color.RGBA) *image.RGBA {
// 	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
// 	for x := 0; x < 4; x++ {
// 		for y := 0; y < 4; y++ {
// 			img.Set(x, y, pixelColor)
// 		}
// 	}
// 	return img
// }

func (l *MyLogo) FillPixels(colorFunction func(int, int) color.RGBA) {
	for x := 0; x < l.width; x++ {
		for y := 0; y < l.height; y++ {
			l.img.SetRGBA(x, y, colorFunction(x, y))
		}
	}
}

// func randomColor() color.RGBA {
// 	rand := rand.New(rand.NewSource(time.Now().Unix()))
// 	return color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}
// }

func newLogoInstance(width int, height int) *MyLogo {
	return &MyLogo{
		width:  width,
		height: height,
		img:    image.NewRGBA(image.Rect(0, 0, width, height)),
	}
}

func (l *MyLogo) Save(path string) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	png.Encode(f, l.img)
}

func main() {
	width := 300
	height := 300
	myLogo := newLogoInstance(width, height)
	myLogo.FillPixels(func(x, y int) color.RGBA {
		return color.RGBA{
			uint8(x * 255 / width),
			uint8(y * 255 / height),
			100,
			255,
		}
		//return randomColor()
	})
	myLogo.Save("Logo.png")
}

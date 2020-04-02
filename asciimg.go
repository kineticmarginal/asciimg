package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
	// Side-effect import.
	// Сайд-эффект — добавление декодера PNG в пакет image.
	_ "image/png"
	"os"
	// Внешняя зависимость.
	"golang.org/x/image/draw"
	"github.com/buger/goterm"
	colorkit "github.com/gookit/color"
)

var (
	out   = flag.String("o", "", "file to write")
	noscale = flag.Bool("noscale", false, "Do not scale the image")
	width = flag.Int("w", 200, "Image width")
	height = flag.Int("h", 40, "Image height")
	termimal_size = flag.Bool("term", false, "Get termial size")
	colorize = flag.Bool("c", false, "Colorize image")
)

func resize(img image.Image, w int, h int) image.Image {
	dstImg := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.NearestNeighbor.Scale(dstImg, dstImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dstImg
}

func scale_image(img image.Image) image.Image {
	if *termimal_size {
		*width = goterm.Width()
		*height = goterm.Height()
	}
	if *out == "" && !*noscale {
		img = resize(img, *width, *height)
	}
	return img
}

func decodeImageFile(imgName string) (image.Image, error) {
	imgFile, err := os.Open(imgName)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(imgFile)

	return img, err
}

func processPixel(c color.Color) byte {
	symbols := []rune("MND8OZ$7I?+=~:,.  ")
	gc := color.GrayModel.Convert(c)
	r, _, _, _ := gc.RGBA()
	r = r >> 8
	pos := uint32(r * 16 / 255)
	return byte(symbols[pos])
}

func convertToAscii(img image.Image) [][]byte {
	textImg := make([][]byte, img.Bounds().Dy())
	for i := range textImg {
		textImg[i] = make([]byte, img.Bounds().Dx())
	}

	for i := range textImg {
		for j := range textImg[i] {
			textImg[i][j] = processPixel(img.At(j, i))
		}
	}
	return textImg
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: asciimg <image.jpg>")
		os.Exit(0)
	}
	imgName := flag.Arg(0)

	img, err := decodeImageFile(imgName)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	img = scale_image(img)

	textImg := convertToAscii(img)
	
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			print(err)
			os.Exit(1)
		}

		for i := range textImg {
			err = ioutil.WriteFile(*out, append([]byte(textImg[i]), []byte("\n")...), 0644)
			if err != nil {
				fmt.Println("Error:", err.Error())
				os.Exit(1)
			}
		}
		defer f.Close()
		fmt.Println("Image in your file!")
	} else {
		if *colorize {
			for i := range textImg {
				for j := range textImg[i] {
					r, g, b, _ := color.GrayModel.Convert(img.At(j, i)).RGBA()
					c := colorkit.RGB(uint8(r), uint8(g), uint8(b))
					c.Print(string(textImg[i][j]))
				}
				fmt.Println()
			}
		} else {
			for i := range textImg {
				for j := range textImg[i] {
					fmt.Print(string(textImg[i][j]))
				}
				fmt.Println()
			}
		}
	}
}
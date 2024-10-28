package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	var inputPath string
	var outputPath string

	flag.StringVar(&inputPath, "i", "", "[required] path to source image")
	flag.StringVar(&outputPath, "o", "", "[required] output path to new image")

	flag.Parse()

	if inputPath == "" || outputPath == "" {
		fmt.Println("No required arguments")
		os.Exit(1)
	}

	err := trim(inputPath, outputPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Complete, and saved to: %s\n", outputPath)
}

func cropImage(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	left, top, right, bottom := width, height, 0, 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			_, _, _, a := img.At(x, y).RGBA()

			a >>= 8

			if a != 0 {
				if x < left {
					left = x
				}
				if x > right {
					right = x
				}
				if y < top {
					top = y
				}
				if y > bottom {
					bottom = y
				}
			}
		}
	}

	if left > right || top > bottom {
		return nil, errors.New("not found alpha channel in border of this image")
	}

	newWidth := right - left + 1
	newHeight := bottom - top + 1
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := top; y <= bottom; y++ {
		for x := left; x <= right; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			r >>= 8
			g >>= 8
			b >>= 8
			a >>= 8

			newImg.Set(x-left, y-top, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	return newImg, nil
}

func trim(inputPath, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return errors.New(fmt.Sprint("Error open file:", err))
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return errors.New(fmt.Sprint("Error decode image:", err))
	}

	newImg, err := cropImage(img)
	if err != nil {
		return errors.New(fmt.Sprint("Error crop image:", err))
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return errors.New(fmt.Sprint("Error create file:", err))
	}
	defer outFile.Close()

	err = png.Encode(outFile, newImg)
	if err != nil {
		return errors.New(fmt.Sprint("Error save image:", err))
	}

	return nil
}

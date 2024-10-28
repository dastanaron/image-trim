package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"image"
	"image/color"
	"io"
	"os"
	"testing"
)

func calculateFileHash(file *os.File) (string, error) {
	hasher := sha256.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", errors.New("Error read file")
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash, nil
}

func TestCropImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	for y := 3; y < 7; y++ {
		for x := 3; x < 7; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	newImg, err := cropImage(img)

	if err != nil {
		t.Error(err)
	}

	// Проверяем размеры нового изображения
	if newImg.Bounds().Dx() != 4 || newImg.Bounds().Dy() != 4 {
		t.Errorf("Incorrect dimensions for new image: 4x4 expected %dx%d", newImg.Bounds().Dx(), newImg.Bounds().Dy())
	}

	// Проверяем, что все пиксели в новом изображении непрозрачные
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			_, _, _, a := newImg.At(x, y).RGBA()
			a >>= 8
			if a != 255 {
				t.Errorf("Pixel (%d, %d) has alpha channel %d, expected 255", x, y, a)
			}
		}
	}
}

func TestTrimImage(t *testing.T) {

	type testObject struct {
		ExpectedHash string
		InputPath    string
		OutputPath   string
	}

	tests := []testObject{
		{
			ExpectedHash: "bb4e7726922e4353059e5cfa54706b9519489dca5b031a3d156de5948dfa031e",
			InputPath:    "./tests/1.png",
			OutputPath:   "./tests/1_res.png",
		},
		{
			ExpectedHash: "0e09959eb3fcca7cabc27911407abb806fe934d33fd032c3aa6a96f7ea6a3dd8",
			InputPath:    "./tests/2.png",
			OutputPath:   "./tests/2_res.png",
		},
		{
			ExpectedHash: "5f4c9d383df1566e138f494992c7f9a4ee9bd7abd671fe18e32630ca529dc564",
			InputPath:    "./tests/3.png",
			OutputPath:   "./tests/3_res.png",
		},
	}

	for _, testObject := range tests {
		trim(testObject.InputPath, testObject.OutputPath)

		file, err := os.Open(testObject.OutputPath)

		if err != nil {
			t.Error("Error open file: ", testObject.OutputPath)
		}

		hash, err := calculateFileHash(file)

		if err != nil {
			t.Error("Cannot calculate hash")
		}

		if hash != testObject.ExpectedHash {
			t.Errorf("Incorrect hash for new image: %s expected %s", hash, testObject.ExpectedHash)
		}

		os.Remove(testObject.OutputPath)

	}
}

func BenchmarkCropImage(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	for y := 30; y < 70; y++ {
		for x := 30; x < 70; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cropImage(img)
	}
}

package main

import (
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"strconv"
)

var smallSq image.Image = image.NewRGBA(image.Rect(0, 0, 1, 1))

func TestResizeImage(t *testing.T) {
	assert := assert.New(t)
	filename := filepath.Join("testdata", "Lenna.png")

	testcases := []struct {
		inWidth, inHeight, outWidth, outHeight int
		inFormat, outFormat                    string
	}{
		{0, 0, 512, 512, "", "png"},
		{232, 232, 232, 232, "jpeg", "jpeg"},
		{123, 0, 123, 123, "png", "png"},
		{0, 222, 222, 222, "gif", "gif"},
		{300, 200, 300, 200, "bmp", "bmp"},
		{100, 100, 100, 100, "tiff", "tiff"},
	}

	for _, testcase := range testcases {
		func() {
			inWidth := testcase.inWidth
			inHeight := testcase.inHeight
			inFormat := testcase.inFormat
			outWidth := testcase.outWidth
			outHeight := testcase.outHeight
			outFormat := testcase.outFormat

			resizeImage(&cliFlags{
				width:  inWidth,
				height: inHeight,
				format: inFormat,
				interp: 0,
			}, filename)
			newFilename := filepath.Join("testdata", "Lenna_"+strconv.Itoa(outWidth)+"x"+strconv.Itoa(outHeight)+"."+outFormat)

			f, err := os.Open(newFilename)
			defer os.Remove(newFilename)
			checkError(err)

			img, format, err := image.Decode(f)
			checkError(err)

			assert.Equal(outFormat, format)
			imgBounds := img.Bounds()
			assert.Equal(outWidth, imgBounds.Dx())
			assert.Equal(outHeight, imgBounds.Dy())
		}()
	}
}

func TestShouldFill(t *testing.T) {
	// Valid
	assert.True(t, shouldFillCliFlags(123, 123, false))

	// Invalid
	assert.False(t, shouldFillCliFlags(0, 0, true))
	assert.False(t, shouldFillCliFlags(0, 0, false))
	assert.False(t, shouldFillCliFlags(123, 0, false))
	assert.False(t, shouldFillCliFlags(0, 123, false))
	assert.False(t, shouldFillCliFlags(0, 123, true))
	assert.False(t, shouldFillCliFlags(123, 0, true))
	assert.False(t, shouldFillCliFlags(123, 123, true))
}

func shouldFillCliFlags(width, height int, stretch bool) bool {
	return shouldPaintBackground(&cliFlags{
		width:   width,
		height:  height,
		stretch: stretch,
	})
}

func TestFitBounds(t *testing.T) {
	var (
		img1 = image.NewRGBA(image.Rect(0, 0, 120, 80)) // 3:2
		img2 = image.NewRGBA(image.Rect(0, 0, 80, 120)) // 2:3
	)

	// img1
	{
		r := fitBounds(img1, 100, 50, 0)
		rBounds := r.Bounds()

		assert.Equal(t, 75, rBounds.Dx())
		assert.Equal(t, 50, rBounds.Dy())
	}

	// img2
	{
		r := fitBounds(img2, 100, 9999, 0)
		rBounds := r.Bounds()

		assert.Equal(t, 100, rBounds.Dx())
		assert.Equal(t, 150, rBounds.Dy())
	}
}

func TestGetFillColor(t *testing.T) {
	// Valid
	{
		color, err := getBackgroundColor(1)
		assert.NoError(t, err)
		assert.Equal(t, image.Black, color)
	}

	{
		color, err := getBackgroundColor(2)
		assert.NoError(t, err)
		assert.Equal(t, image.White, color)
	}

	// Invalid
	{
		color, err := getBackgroundColor(0)
		assert.Nil(t, color)
		assert.Error(t, err)
		assert.Equal(t, "Unknown background color: 0", err.Error())
	}

	{
		color, err := getBackgroundColor(3)
		assert.Nil(t, color)
		assert.Error(t, err)
		assert.Equal(t, "Unknown background color: 3", err.Error())
	}
}

func TestAddResizedSuffix(t *testing.T) {
	assert.Equal(t, "test_1280x720.png", addResizedSuffix("test.jpg", "1280x720", "png"))
	assert.Equal(t, "/home/user/img_444x444.gif", addResizedSuffix("/home/user/img.png", "444x444", "gif"))
}

func TestWriteImageFile(t *testing.T) {
	// Valid
	filename := filepath.Join(os.TempDir(), "writeimagetest")
	defer os.Remove(filename)

	writeImageFile(filename, "png", smallSq)
	f, err := ioutil.ReadFile(filename)
	assert.NoError(t, err)
	assert.True(t, len(f) > 0)
}

func TestEncode(t *testing.T) {
	file := createTempFile()
	defer deleteFile(file)

	// Valid
	validExts := []string{
		"bmp",
		"gif",
		"jpg",
		"jpeg",
		"png",
		"tiff",
	}

	assert := assert.New(t)
	for _, ext := range validExts {
		err := encode(file, smallSq, ext)
		assert.NoError(err)
	}

	// Invalid
	err := encode(file, smallSq, "fail")
	assert.Error(err)
	assert.Equal("Unknown encode format: \"fail\"", err.Error())
}

func createTempFile() (file *os.File) {
	file, err := ioutil.TempFile("", "test")

	if err != nil {
		deleteFile(file)
		log.Fatal(err)
	}

	return
}

func deleteFile(file *os.File) {
	os.Remove(file.Name())
}

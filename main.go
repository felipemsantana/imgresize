package main

import (
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/nfnt/resize"
)

type cliFlags struct {
	width, height, interp, background int
	stretch                           bool
	format                            string
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Felipe Matos Santana",
			Email: "felipems@yahoo.com.br",
		},
	}
	app.Name = "Image Resizer"
	app.Usage = "Tool to resize images, it supports JPG, GIF and PNG"
	app.HelpName = "imgresize"
	app.HideHelp = true
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "height, h",
			Value: 0,
			Usage: "output image height, default is 0, which preserves original aspect ratio",
		},
		cli.IntFlag{
			Name:  "width, w",
			Value: 0,
			Usage: "output image width, default is 0, which preserves original aspect ratio",
		},
		cli.IntFlag{
			Name:  "interp, i",
			Value: 5,
			Usage: `interpolation function, from 0 (fastest) to 5 (slowest), default is 5:
	0: Nearest-neighbor interpolation
	1: Bilinear interpolation
	2: Bicubic interpolation
	3: Mitchell-Netravali interpolation
	4: Lanczos resampling with a=2
	5: Lanczos resampling with a=3`,
		},
		cli.IntFlag{
			Name:  "background, b",
			Value: 0,
			Usage: `background color, used only if stretch is false and aspect ratio from the output image is not the same as the input, default is 0:
	0: transparent
	1: black
	2: white`,
		},
		cli.BoolFlag{
			Name:  "stretch, s",
			Usage: "stretch image, default is false",
		},
		cli.StringFlag{
			Name: "format, f",
			Usage: `output image format, default is same as input:
	- png
	- jpg
	- gif`,
		},
	}
	app.Action = handleActions
	app.Run(os.Args)
}

func handleActions(c *cli.Context) {
	totalArgs := c.NArg()
	if totalArgs == 0 {
		cli.ShowAppHelp(c)
		return
	}

	flags := &cliFlags{
		width:      c.GlobalInt("width"),
		height:     c.GlobalInt("height"),
		interp:     c.GlobalInt("interp"),
		background: c.GlobalInt("background"),
		format:     c.GlobalString("format"),
		stretch:    c.GlobalBool("stretch"),
	}

	wg := new(sync.WaitGroup)
	for _, arg := range c.Args() {
		resizeImageAsync(flags, arg, wg)
	}
	wg.Wait()
}

func resizeImageAsync(flags *cliFlags, filename string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		resizeImage(flags, filename)
	}()
}

func resizeImage(flags *cliFlags, filename string) {
	file, err := os.Open(filename)
	checkError(err)

	img, inputFormat, err := image.Decode(file)
	defer file.Close()
	checkError(err)

	resized := doResize(flags, img)
	resizedBounds := resized.Bounds()

	outWidth := strconv.Itoa(resizedBounds.Dx())
	outHeight := strconv.Itoa(resizedBounds.Dy())
	resolution := outWidth + "x" + outHeight
	format := flags.format
	if format == "" {
		format = inputFormat
	}

	outputFilename := addResizedSuffix(filename, resolution, format)

	writeImageFile(outputFilename, format, resized)

	log.Println("Image resized successfully:", outputFilename)
}

func doResize(flags *cliFlags, img image.Image) (resized image.Image) {
	interpFunc := resize.InterpolationFunction(flags.interp)
	if shouldPaintBackground(flags) {
		return paintBackground(flags, img, interpFunc)
	}
	return resize.Resize(uint(flags.width), uint(flags.height), img, interpFunc)
}

func shouldPaintBackground(flags *cliFlags) bool {
	return flags.width > 0 && flags.height > 0 && !flags.stretch
}

func paintBackground(flags *cliFlags, img image.Image, interpFunc resize.InterpolationFunction) image.Image {
	width, height := flags.width, flags.height
	temp := fitBounds(img, width, height, interpFunc)
	tempBounds := temp.Bounds()

	r := image.Rect(0, 0, width, height)

	newImg := image.NewRGBA(r)

	// transparent is default, no need to repaint
	background := flags.background
	if background != 0 {
		backgroundColor, err := getBackgroundColor(background)
		checkError(err)

		draw.Draw(newImg, r, backgroundColor, image.ZP, draw.Src)
	}

	point := calcDrawPoint(tempBounds, r)

	draw.Draw(newImg, r, temp, point, draw.Src)
	return newImg
}

func calcDrawPoint(a, b image.Rectangle) (point image.Point) {
	return image.Point{
		X: (a.Dx() - b.Dx()) / 2,
		Y: (a.Dy() - b.Dy()) / 2,
	}
}

func fitBounds(img image.Image, width, height int, interpFunc resize.InterpolationFunction) image.Image {
	imgBounds := img.Bounds()
	ratioX := float32(width) / float32(imgBounds.Dx())
	ratioY := float32(height) / float32(imgBounds.Dy())

	if ratioX > ratioY {
		return resize.Resize(0, uint(height), img, interpFunc)
	}
	return resize.Resize(uint(width), 0, img, interpFunc)
}

func getBackgroundColor(background int) (*image.Uniform, error) {
	switch background {
	case 1:
		return image.Black, nil
	case 2:
		return image.White, nil
	default:
		return nil, errors.New("Unknown background color: " + strconv.Itoa(background))
	}
}

var filenameRegexp = regexp.MustCompile("^(.+)(\\.\\w+)$")

func addResizedSuffix(filename, resolution, format string) string {
	return filenameRegexp.ReplaceAllString(filename, "${1}_"+resolution+"."+format)
}

func writeImageFile(filename, format string, img image.Image) {
	out, err := os.Create(filename)
	defer func() {
		out.Close()
		if r := recover(); r != nil {
			os.Remove(filename)
			log.Fatal(r.(error))
		}
	}()
	panicError(err)

	err = encode(out, img, format)
	panicError(err)
}

func encode(out io.Writer, resized image.Image, format string) error {
	switch strings.ToLower(format) {
	case "png":
		return png.Encode(out, resized)
	case "jpg", "jpeg":
		return jpeg.Encode(out, resized, nil)
	case "gif":
		return gif.Encode(out, resized, nil)
	default:
		return errors.New("Unknown encode format: \"" + format + "\"")
	}
}

[![Latest Release](https://img.shields.io/github/release/felipemsantana/imgresize.svg?maxAge=2592000)](https://github.com/felipemsantana/imgresize/releases/latest)
[![Build Status](https://travis-ci.org/felipemsantana/imgresize.svg?branch=master)](https://travis-ci.org/felipemsantana/imgresize)

# imgresize
Simple command line tool to resize images

## Instalation
Download the latest binary file for your operating system and architecture [here](https://github.com/felipemsantana/imgresize/releases/latest).

Or download and install from the source:
```
$ go get github.com/felipemsantana/imgresize
```

## Usage
```
NAME:
   Image Resizer - Tool to resize images, it supports BMP, GIF, JPEG, PNG and TIFF

USAGE:
   imgresize [global options] [arguments...]

VERSION:
   1.3.0

AUTHOR(S):
   Felipe Matos Santana <felipems@yahoo.com.br>

GLOBAL OPTIONS:
   --height value, -h value	output image height, value 0 preserves original aspect ratio (default: 0)
   --width value, -w value	output image width, value 0 preserves original aspect ratio (default: 0)
   --interp value, -i value	interpolation function, from 0 (fastest) to 5 (slowest):
				0: Nearest-neighbor interpolation
				1: Bilinear interpolation
				2: Bicubic interpolation
				3: Mitchell-Netravali interpolation
				4: Lanczos resampling with a=2
				5: Lanczos resampling with a=3 (default: 5)
   --background value, -b value	background color, used only if stretch is false and aspect ratio from the output image is not the same as the input:
				0: transparent
				1: black
				2: white (default: 0)
   --stretch, -s		stretch image, default is false
   --format value, -f value	output image format, default is same as input:
				- bmp
				- gif
				- jpg
				- png
				- tif
   --version, -v		print the version

```

## Example
This will resize any supported image format to 1920x1080 without stretching and save it as JPG:
```
$ imgresize -w 1920 -h 1080 -f jpg myimage.jpg
```
The resized image file will be saved as "myimage_1920x1080.jpg" in the same directory.

## License
[MIT](LICENSE)

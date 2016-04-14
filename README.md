[![Latest Release](https://img.shields.io/github/release/fmatoss/imgresize.svg?maxAge=2592000)](https://github.com/fmatoss/imgresize/releases/latest)
[![Build Status](https://travis-ci.org/fmatoss/imgresize.svg?branch=master)](https://travis-ci.org/fmatoss/imgresize)
[![Coverage Status](https://coveralls.io/repos/github/fmatoss/imgresize/badge.svg?branch=master)](https://coveralls.io/github/fmatoss/imgresize?branch=master)

# imgresize
Simple command line tool to resize images

## Instalation
Download the latest binary file for your operating system and architecture [here](https://github.com/fmatoss/imgresize/releases/latest).

Or download and install from the source:
```
$ go get github.com/fmatoss/imgresize
```

## Usage
```
NAME:
   Image Resizer - Tool to resize images, it supports JPG, GIF and PNG

USAGE:
   imgresize [global options] [arguments...]
   
VERSION:
   1.1.0
   
AUTHOR(S):
   Felipe Matos Santana <felipems@yahoo.com.br> 
   
GLOBAL OPTIONS:
   --height, -h "0"	output image height, default is 0, which preserves original aspect ratio
   --width, -w "0"	output image width, default is 0, which preserves original aspect ratio
   --interp, -i "5"	interpolation function, from 0 (fastest) to 5 (slowest), default is 5:
			0: Nearest-neighbor interpolation
			1: Bilinear interpolation
			2: Bicubic interpolation
			3: Mitchell-Netravali interpolation
			4: Lanczos resampling with a=2
			5: Lanczos resampling with a=3
   --background, -b "0"	background color, used only if stretch is false and aspect ratio from the output image is not the same as the input, default is 0:
			0: transparent
			1: black
			2: white
   --stretch, -s	stretch image, default is false
   --format, -f 	output image format, default is same as input:
			- bmp
			- gif
			- jpg
			- png
   --version, -v	print the version

```

## Example
This will resize any supported image format to 1920x1080 without stretching and save it as JPG:
```
$ imgresize -w 1920 -h 1080 -f jpg myimage.jpg
```
The resized image file will be saved as "myimage_1920x1080.jpg" in the same directory.

## License
[MIT](LICENSE)

package mandelbrot

import (
	"image"
	"image/color"
	"math"
	//"math/cmplx"
)

type Pixel struct {
	Inside     bool
	Iterations int
}

type Bitmap struct {
	bitmap    [][]Pixel
	width     int
	GetColour func(Pixel) (r, g, b uint8)
}

func getIterations(c complex128, maxIterations int) Pixel {
	z := complex128(0)

	for i := 0; i < maxIterations; i++ {
		z = z*z + c
		if real(z)*real(z) + imag(z)*imag(z) > 4 {
			return Pixel{false, i}
		}
	}

	return Pixel{true, 0}
}

func MakeBitmap(nw, se complex128, width, maxIterations int) Bitmap {
	stepReal := (real(se) - real(nw)) / float64(width)
	stepImag := (imag(se) - imag(nw)) / float64(width)

	out := make([][]Pixel, width)

	for x := 0; x < width; x = x + 1 {
		out[x] = make([]Pixel, width)

		for y := 0; y < width; y = y + 1 {
			out[x][y] = getIterations(nw+complex(float64(x)*stepReal, float64(y)*stepImag), maxIterations)
		}
	}

	return Bitmap{bitmap: out, width: width, GetColour: BlackAndWhite}
}

func (bm Bitmap) ColorModel() color.Model {
	return color.RGBAModel
}

func (bm Bitmap) Bounds() image.Rectangle {
	return image.Rect(0, 0, bm.width, bm.width)
}

func (bm Bitmap) At(x, y int) color.Color {
	r, g, b := bm.GetColour(bm.bitmap[x][y])
	return color.RGBA{r, g, b, 255}
}

func BlackAndWhite(p Pixel) (r, g, b uint8) {
	if p.Inside {
		return 255, 255, 255
	}
	return 0, 0, 0
}

func HSLToRGB(h, s, l float64) (r, g, b uint8) {
	var r1, g1, b1 float64

	c := (1 - math.Abs(2*l-1)) * s
	hdash := h * 6
	x := c * (1 - math.Abs(math.Mod(hdash, 2)-1))

	switch {
	case 0 <= hdash && hdash < 1:
		r1, g1, b1 = c, x, 0
	case 1 <= hdash && hdash < 2:
		r1, g1, b1 = x, c, 0
	case 2 <= hdash && hdash < 3:
		r1, g1, b1 = 0, c, x
	case 3 <= hdash && hdash < 4:
		r1, g1, b1 = 0, x, c
	case 4 <= hdash && hdash < 5:
		r1, g1, b1 = x, 0, c
	case true:
		r1, g1, b1 = c, 0, x
	}

	m := l - 0.5*c

	r = uint8((r1 + m) * 255)
	g = uint8((g1 + m) * 255)
	b = uint8((b1 + m) * 255)

	return
}

func Flame(p Pixel) (r, g, b uint8) {
	if p.Inside {
		return 0, 0, 0
	}

	iterations := float64(p.Iterations % 64)

	var h float64
	
	if iterations < 32 {
		h = iterations / (32 * 8) + 0.05
	} else {
		h = (63 - iterations) / (32 * 8) + 0.05
	}

	r, g, b = HSLToRGB(h, 1, 0.5)

	return
}

func BlueGreen(p Pixel) (r, g, b uint8) {
	if p.Inside {
		return 0, 0, 0
	}

	iterations := float64(p.Iterations % 80)

	var temp, h, l float64
	
	if iterations < 40 {
		temp = iterations / (40 * 5)
	} else {
		temp = (79 - iterations) / (40 * 5)
	}

	h = temp + 0.4
	l = (0.5 - temp * 1.5)

	r, g, b = HSLToRGB(h, 0.5, l)

	return
}

func Multicolour(p Pixel) (r, g, b uint8) {
	if p.Inside {
		return 0, 0, 0
	}

	var y, cb, cr, iterations uint8

	iterations = uint8(p.Iterations % 64)
	y = 193

	switch {
	case iterations < 16:
		cb = iterations * 16
		cr = 0
	case iterations < 32:
		cb = 255
		cr = uint8((iterations % 16) * 16)
	case iterations < 48:
		cb = uint8((15 - (iterations % 16)) * 16)
		cr = 255
	case true:
		cb = 0
		cr = uint8((15 - (iterations % 16)) * 16)
	}

	return color.YCbCrToRGB(y, cb, cr)
}

func Stripey(p Pixel) (r, g, b uint8) {
	if p.Inside {
		return 0, 0, 0
	}

	c := uint8(p.Iterations % 2 * 255)

	return c, c, c
}

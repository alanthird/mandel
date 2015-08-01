package mandelbrot

import (
	"image"
	"image/color"
	"math/cmplx"
)

type Pixel struct {
	inside     bool
	iterations int
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
		if cmplx.Abs(z) > 2 {
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
	if p.inside {
		return 255, 255, 255
	}
	return 0, 0, 0
}

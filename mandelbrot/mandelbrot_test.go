package mandelbrot

import (
	"testing"
	"fmt"
)

func TestGetIterations(t *testing.T) {
	var p Pixel
	
	p = getIterations(2+2i, 32) 
	if p.Inside != false || p.Iterations != 0 {
		t.Error(fmt.Sprintf("2+2i: inside: %t, iterations: %d", p.Inside, p.Iterations))
	}

	p = getIterations(0+0i, 32)
	if  p.Inside != true || p.Iterations != 0 {
		t.Error(fmt.Sprintf("0+0i: inside: %t, iterations: %d", p.Inside, p.Iterations))
	}

	p = getIterations(0+2i, 32)
	if p.Inside != false || p.Iterations != 1 {
		t.Error(fmt.Sprintf("0+2i: inside: %t, iterations: %d", p.Inside, p.Iterations))
	}
}

func TestHSLToRGB(t *testing.T) {
	if r, g, b := HSLToRGB(0, 0, 1) ; r != 255 || g != 255 || b != 255 {
		t.Error(fmt.Sprintf("H: 0, S: 0, L: 1 should give r: 255, g: 255, b: 255. Actual: r: %d, g: %d, b: %d", r, g, b))
	}

	if r, g, b := HSLToRGB(0, 0, 0) ; r != 0 || g != 0 || b != 0 {
		t.Error(fmt.Sprintf("H: 0, S: 0, L: 0 should give r: 0, g: 0, b: 0. Actual: r: %d, g: %d, b: %d", r, g, b))
	}

	if r, g, b := HSLToRGB(1.0/3, 1, 0.250) ; r != 0 || g != 127 || b != 0 {
		t.Error(fmt.Sprintf("H: 1/3, S: 1, L: 0.25 should give r: 0, g: 127, b: 0. Actual: r: %d, g: %d, b: %d", r, g, b))
	}
}

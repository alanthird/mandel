package mandelbrot

import (
	"testing"
	"fmt"
)

func TestGetIterations(t *testing.T) {
	var p Pixel
	
	p = getIterations(2+2i, 32) 
	if p.inside != false || p.iterations != 0 {
		t.Error(fmt.Sprintf("2+2i: inside: %t, iterations: %d", p.inside, p.iterations))
	}

	p = getIterations(0+0i, 32)
	if  p.inside != true || p.iterations != 0 {
		t.Error(fmt.Sprintf("0+0i: inside: %t, iterations: %d", p.inside, p.iterations))
	}

	p = getIterations(0+2i, 32)
	if p.inside != false || p.iterations != 1 {
		t.Error(fmt.Sprintf("0+2i: inside: %t, iterations: %d", p.inside, p.iterations))
	}
}

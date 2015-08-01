package main

import (
	"fmt"
	"net/http"
	"strings"
	"strconv"
	"image/png"
	"github.com/alanthird/mandel/mandelbrot"
)

type appHandler func(http.ResponseWriter, *http.Request) error
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w,
			fmt.Sprintf("Sorry, I can't find what you're looking for!\n%s", err),
			http.StatusNotFound);
	}
}

func mapHandler(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/") 
	
	z, err := strconv.ParseUint(pathParts[2], 10, 64)
	if err != nil {
		return err
	}

	var imagesPerField uint64 = 1 << z
	
	y, err := strconv.ParseUint(pathParts[3], 10, 64)
	if err != nil || y >= imagesPerField {
		return err
	}

	x, err := strconv.ParseUint(pathParts[4], 10, 64)
	if err != nil || x >= imagesPerField {
		return err
	}

	unitsAcross := 4 / float64(imagesPerField)
	
	nw := complex(unitsAcross * float64(x) - 2, unitsAcross * float64(y) - 2)
	sw := complex(unitsAcross * float64(x+1) - 2, unitsAcross * float64(y+1) - 2)
	
	w.Header().Set("Content-Type", "image/png")

	b := mandelbrot.MakeBitmap(nw, sw, 256, 1000)

	png.Encode(w, b)

	return nil
}

func main() {
    http.Handle("/map/", appHandler(mapHandler))
	http.Handle("/", http.FileServer(http.Dir("static")))
    http.ListenAndServe(":8080", nil)
}

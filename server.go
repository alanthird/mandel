package main

import (
	"fmt"
	"github.com/alanthird/mandel/mandelbrot"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w,
			fmt.Sprintf("Sorry, I can't find what you're looking for!\n%s", err),
			http.StatusNotFound)
	}
}

type layer struct {
	name string
	f    func(mandelbrot.Pixel) (r, g, b uint8)
}

func cacheAllImages(layerList []layer, x, y, z uint64) error {
	var imagesPerField uint64 = 1 << z
	unitsAcross := 4 / float64(imagesPerField)

	nw := complex(unitsAcross*float64(x)-2, unitsAcross*float64(imagesPerField-y)-2)
	sw := complex(unitsAcross*float64(x+1)-2, unitsAcross*float64(imagesPerField-y-1)-2)

	b := mandelbrot.MakeBitmap(nw, sw, 256, 1024)

	for _, layer := range layerList {
		dirName := fmt.Sprintf("static/map/%s/%d/%d", layer.name, z, y)

		os.MkdirAll(dirName, 0755)

		b.GetColour = layer.f

		outFile, err := os.Create(fmt.Sprintf("%s/%d.png", dirName, x))
		if err != nil {
			return err
		}

		defer outFile.Close()

		if err := png.Encode(outFile, b); err != nil {
			return err
		}
	}

	return nil
}

func serveCachedFile(w http.ResponseWriter, r *http.Request) error {
	if _, err := os.Stat(fmt.Sprintf("static/%s.png", r.URL.Path)); err == nil {
		http.ServeFile(w, r, fmt.Sprintf("static/%s.png", r.URL.Path))
		return nil
	} else {
		return err
	}

}

func mapHandler(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")

	z, err := strconv.ParseUint(pathParts[3], 10, 64)
	if err != nil || z > 45 {
		return err
	}

	var imagesPerField uint64 = 1 << z

	y, err := strconv.ParseUint(pathParts[4], 10, 64)
	if err != nil || y >= imagesPerField {
		return err
	}

	x, err := strconv.ParseUint(pathParts[5], 10, 64)
	if err != nil || x >= imagesPerField {
		return err
	}

	if err := serveCachedFile(w, r); err == nil {
		return nil
	}

	layerList := []layer{
		layer{"beetlejuice", mandelbrot.Stripey},
		layer{"colour", mandelbrot.Multicolour},
		layer{"flame", mandelbrot.Flame},
		layer{"bluegreen", mandelbrot.BlueGreen}}
	
	if err := cacheAllImages(layerList, x, y, z); err != nil {
		return err
	}

	return serveCachedFile(w, r)
}

func main() {
	http.Handle("/map/", appHandler(mapHandler))
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":8080", nil)
}

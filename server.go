package main

import (
	"fmt"
	"net/http"
	// "strings"
	// "strconv"
	"image/png"
	"mandel/mandelbrot"
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
	// pathParts := strings.Split(r.URL.Path, "/") 
	
	// z, err := strconv.ParseFloat(pathParts[2], 64)
	// if err != nil {
	// 	return err
	// }

	// y, err := strconv.ParseFloat(pathParts[3], 64)
	// if err != nil {
	// 	return err
	// }

	// x, err := strconv.ParseFloat(pathParts[4], 64)
	// if err != nil {
	// 	return err
	// }

	w.Header().Set("Content-Type", "image/png")

	b := mandelbrot.MakeBitmap(-2+2i, 2-2i, 256, 1000)

	png.Encode(w, b)

	return nil
}

func main() {
    http.Handle("/map/", appHandler(mapHandler))
	http.Handle("/", http.FileServer(http.Dir("static")))
    http.ListenAndServe(":8080", nil)
}

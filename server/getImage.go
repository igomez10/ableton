package main

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

//handlerGetImage handles the path /image/<sha256>
func (d *imageDirectory) handlerGetImage(w http.ResponseWriter, r *http.Request) {
	//parse vars from path
	vars := mux.Vars(r)
	imageID := vars["sha256"]
	//Join all the strings in the image entry
	imageData, err := GetImageData((*d).builtImages, imageID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(nil)
	} else {
		w.Write([]byte(imageData))
	}
}

// GetImageData queries the imageID in the built images, if image not found in built images returns error
func GetImageData(builtImages map[string][]byte, imageID string) ([]byte, error) {

	var imageData []byte
	var err error

	if value, ok := builtImages[imageID]; ok {
		imageData = value
	} else {
		err = errors.New("Image not found")
	}
	return imageData, err
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// handlerRegisterImage handles the path POST /image
func (d *imageDirectory) handlerRegisterImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var response serverResponse
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "Error reading request"
	} else {
		var newImageRequest registerImageRequest
		json.Unmarshal(body, &newImageRequest)
		err = d.RegisterImage(newImageRequest.Sha256, newImageRequest.Size, newImageRequest.ChunkSize)

		if err != nil {
			w.WriteHeader(http.StatusConflict)
			response.Message = "Image already registered"
		} else {
			w.WriteHeader(http.StatusCreated)
			response.Message = "OK"
		}
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(responseBytes)
	}
}

// RegisterImage checks if the image already exists in built or pipeline
// if so, return a Conflict 409.
// if not, create the required slices to store the bytes of data
func (d *imageDirectory) RegisterImage(imageID string, size int, chunkSize int) error {
	//check if image SHA is in built images or is currently being uploaded
	_, imageExistsInProcessed := (*d).pipelineImages[imageID]
	_, imageWasBuilt := (*d).builtImages[imageID]
	imageExists := imageWasBuilt || imageExistsInProcessed

	if imageExists {
		//sha already exists in dictionary, respond with conflict status 409
		msg := fmt.Sprintf("Image %s already exists", imageID)
		log.Println(msg)
		return errors.New(msg)
	}

	// required chunks to upload image given the chunksize and image size
	requiredChunks := size / chunkSize

	if size%chunkSize > 0 {
		requiredChunks++
	}
	//for the new image create an array of array of bytes -> an array of bytes for every chunk
	var newImage splittedImage
	newImage.chunks = make([][]byte, requiredChunks)
	newImage.remainingChunks = requiredChunks

	//initialize chunks arrays with defined size chunkSize
	for i := 0; i < requiredChunks-1; i++ {
		newImage.chunks[i] = make([]byte, chunkSize)
	}

	//the last chunk will be smaller, define special size
	newImage.chunks[requiredChunks-1] = make([]byte, size%chunkSize)
	(*d).pipelineImages[imageID] = &newImage
	log.Printf("Registered %s - size: %d - chunk size: %d", imageID, size, chunkSize)
	return nil
}

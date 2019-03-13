package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// handlerUploadChunk handles the path /image/<sha256>/chunks
func (d *imageDirectory) handlerUploadChunk(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var imageChunkRequest uploadChunkRequest

	//parse variable in path
	vars := mux.Vars(r)

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request: %+v", err)
	} else {
		err = json.Unmarshal(body, &imageChunkRequest)
		if err != nil {
			log.Printf("Error parsing request: %+v", err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			errorUpload := d.UploadChunk(vars["sha256"], imageChunkRequest.ID, []byte(imageChunkRequest.Data))
			if errorUpload != nil {
				log.Println(errorUpload)
				w.WriteHeader(http.StatusBadRequest)
			} else {
			}
		}

	}
}

// UploadChunk pushes the new chunk of data into the image ID entry
// and reduces the counter of remainingChunks for the current image. When this counter
// reaches 0 BuildImage is called
func (d *imageDirectory) UploadChunk(imageID string, chunkID int, chunkData []byte) error {
	var err error
	//do not overwrite old images
	if _, isOldImage := (*d).builtImages[imageID]; isOldImage {
		msgPlaceholder := "Image %s was already uploaded entirely but received a chunk %d with data %s"
		err = fmt.Errorf(msgPlaceholder, imageID, chunkID, string(chunkData))
	} else if _, imageIsRegistered := (*d).pipelineImages[imageID]; imageIsRegistered {
		//log.Printf("Uploading chunk %d/%d to %s", chunkID, len((*d).pipelineImages[imageID].chunks), imageID)
		// push the chunk data to its image ID
		(*d).pipelineImages[imageID].chunks[chunkID] = chunkData

		// reduce the number of remaining chunks for the current image ID
		(*d).pipelineImages[imageID].remainingChunks--

		if (*d).pipelineImages[imageID].remainingChunks == 0 {
			//all chunks received, build final image
			(*d).BuildImage(imageID)
		}
	} else {
		err = fmt.Errorf("Image %s is not registered", imageID)
	}
	return err
}

// BuildImage joins all the chunks,
// stores the new single array of bytes in a
// different map and removes the old splitted data from the pipeline
func (d *imageDirectory) BuildImage(imageID string) {
	completeImage := bytes.Join((*d).pipelineImages[imageID].chunks, nil)
	//create new entry in map of built images
	d.builtImages[imageID] = completeImage
	//remove from old map, cleanup process
	delete(d.pipelineImages, imageID)
}

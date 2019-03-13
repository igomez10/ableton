package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	host := flag.String("host", "0.0.0.0", "hostname used to mount the server, defaults to localhost")
	port := flag.Int("port", 4444, "port number to mount the listen for connections, defaults to 4444")

	// create directory object where images will be stored
	var directory imageDirectory
	// initialize map for images being currently uploaded and complete/built images
	directory.pipelineImages = make(map[string]*splittedImage)
	directory.builtImages = make(map[string][]byte)

	router := mux.NewRouter()

	//register a new image in the directory
	router.HandleFunc("/image", directory.handlerRegisterImage).Methods("POST")

	//upload a chunk of data into an already registered image
	router.HandleFunc("/image/{sha256}/chunks", directory.handlerUploadChunk).Methods("POST")

	//get an image from the directory
	router.HandleFunc("/image/{sha256}", directory.handlerGetImage).Methods("GET")

	//health endpoint - for debugging and internal use
	router.HandleFunc("/hashes", directory.getHashes)

	//health endpoint - for debugging and internal use
	router.HandleFunc("/health", getHealth)

	//mount the server listening on host:port
	address := fmt.Sprintf("%s:%d", *host, *port)
	http.ListenAndServe(address, router)
}

func (d *imageDirectory) getHashes(w http.ResponseWriter, r *http.Request) {
	builder := strings.Builder{}
	for k, v := range d.builtImages {
		builder.WriteString(k)
		builder.WriteString("\n")
		builder.WriteString(string(v))
		builder.WriteString("\n")
	}
	w.Write([]byte(builder.String()))
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

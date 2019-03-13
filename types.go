package main

// upploadChunkRequest is the structure from the client request to push a new chunk
type uploadChunkRequest struct {
	ID   int    `json:"id"`
	Size int    `json:"size"`
	Data string `json:"data"`
}

// registerImageRequest is the structure from the client request to register a new image
type registerImageRequest struct {
	Sha256    string `json:"sha256"`
	Size      int    `json:"size"`
	ChunkSize int    `json:"chunk_size"`
}

// imageDirectory is the single directory where all the images are stored
type imageDirectory struct {
	pipelineImages map[string]*splittedImage // mapping to images in build process, waiting for other chunks
	builtImages    map[string][]byte         // mapping to images with all its the chunks and joined them into to single []byte
}

// splittedImage is an image in upload process
type splittedImage struct {
	chunks          [][]byte
	remainingChunks int //waiting to receive N chunks from client to have all chunks
}

type serverResponse struct {
	Message string `json:"message"`
}

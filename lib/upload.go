package lib

import (
	"bytes"
	"context"
	"sync"
)

//UploadService represents how our upload service should be structured
//in the future this will just be a way to communicate with our external
//upload service but currently it is being handled locally for simplicity
type UploadService interface {
	Save(ctx context.Context, file *File) (url string, err error)
}

//FileID is the primitive type for our file id. Used in upload memory store
//to identify a file and in our actual file struct to map it to said memory
//store
type FileID string

//File represents our generic file structure and what most (if not all)
//will require while the upload process takes place. This is only held
//in memory until the file has been successfully processed and is stored
//in the desired location.
type File struct {
	Data *bytes.Buffer
	Name string
	Type string
	ID   FileID
}

//UploadMemoryStore is our temporary upload store. It will hold the data
//for files and images while they are being chunked from the client and
//
type UploadMemoryStore struct {
	data  map[FileID]*File
	mutex sync.RWMutex
}

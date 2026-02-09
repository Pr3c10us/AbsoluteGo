package storage

import "os"

type UploadInput struct {
	ObjectKey string
	File      *os.File
}

type UploadResult struct {
	ObjectKey string
	URL       string
	Err       error
}

type DeleteResult struct {
	URL string
	Err error
}

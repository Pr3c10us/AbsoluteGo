package storage

import (
	"os"
)

type Interface interface {
	UploadFile(bucketName string, file *os.File) (string, error)
	UploadMany(bucketName string, files []*os.File) []UploadResult
	DeleteFile(fileURL string) error
	DeleteMany(urls []string) []DeleteResult
}

package storage

import (
	"os"
)

type Interface interface {
	UploadFile(bucketName, objectKey string, file *os.File) (string, error)
	UploadMany(bucketName string, files []UploadInput) []UploadResult
	DeleteFile(bucketName, objectKey string) error
}

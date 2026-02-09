package storage

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/minio/minio-go/v7"
)

type minioStorageRepository struct {
	client *minio.Client
}

func NewMinioStorageRepository(client *minio.Client) storage.Interface {
	return &minioStorageRepository{
		client: client,
	}
}

func (r *minioStorageRepository) UploadFile(bucketName, objectKey string, file *os.File) (string, error) {
	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}

	_, err = r.client.PutObject(context.TODO(), bucketName, objectKey, file, stat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	endpoint := r.client.EndpointURL()
	return fmt.Sprintf("%s/%s/%s", endpoint, bucketName, objectKey), nil
}

func (r *minioStorageRepository) UploadMany(bucketName string, files []storage.UploadInput) []storage.UploadResult {
	results := make([]storage.UploadResult, len(files))
	var wg sync.WaitGroup

	for i, f := range files {
		wg.Add(1)
		go func(idx int, input storage.UploadInput) {
			defer wg.Done()
			url, err := r.UploadFile(bucketName, input.ObjectKey, input.File)
			results[idx] = storage.UploadResult{
				ObjectKey: input.ObjectKey,
				URL:       url,
				Err:       err,
			}
		}(i, f)
	}

	wg.Wait()
	return results
}

func (r *minioStorageRepository) DeleteFile(bucketName, objectKey string) error {
	return r.client.RemoveObject(context.TODO(), bucketName, objectKey, minio.RemoveObjectOptions{})
}

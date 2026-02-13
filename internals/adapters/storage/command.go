package storage

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type minioStorageImplementation struct {
	client *minio.Client
}

func NewMinioStorageImplementation(client *minio.Client) storage.Interface {
	return &minioStorageImplementation{
		client: client,
	}
}

func (r *minioStorageImplementation) parseURL(rawURL string) (bucket, objectKey string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse URL %q: %w", rawURL, err)
	}

	trimmed := strings.TrimPrefix(u.Path, "/")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid storage URL %q: expected /<bucket>/<objectKey>", rawURL)
	}

	return parts[0], parts[1], nil
}

func (r *minioStorageImplementation) generateObjectKey(filename string) string {
	ext := filepath.Ext(filename)
	uniqueID := uuid.New().String()
	return uniqueID + ext
}

func (r *minioStorageImplementation) UploadFile(bucketName string, file *os.File) (string, error) {
	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}

	objectKey := r.generateObjectKey(stat.Name())

	_, err = r.client.PutObject(context.TODO(), bucketName, objectKey, file, stat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	endpoint := r.client.EndpointURL()
	return fmt.Sprintf("%s/%s/%s", endpoint, bucketName, objectKey), nil
}

func (r *minioStorageImplementation) UploadMany(bucketName string, files []*os.File) []storage.UploadResult {
	results := make([]storage.UploadResult, len(files))
	var wg sync.WaitGroup

	for i, file := range files {
		wg.Add(1)
		go func(idx int, f *os.File) {
			defer wg.Done()
			u, err := r.UploadFile(bucketName, f)

			var objectKey string
			if err == nil {
				_, objectKey, _ = r.parseURL(u)
			}

			results[idx] = storage.UploadResult{
				ObjectKey: objectKey,
				URL:       u,
				Err:       err,
			}
		}(i, file)
	}

	wg.Wait()
	return results
}

func (r *minioStorageImplementation) DeleteFile(fileURL string) error {
	bucket, key, err := r.parseURL(fileURL)
	if err != nil {
		return fmt.Errorf("DeleteFile: %w", err)
	}

	return r.client.RemoveObject(context.TODO(), bucket, key, minio.RemoveObjectOptions{})
}

func (r *minioStorageImplementation) DeleteMany(urls []string) []storage.DeleteResult {
	results := make([]storage.DeleteResult, len(urls))
	var wg sync.WaitGroup

	for i, rawURL := range urls {
		wg.Add(1)
		go func(idx int, u string) {
			defer wg.Done()
			err := r.DeleteFile(u)
			results[idx] = storage.DeleteResult{
				URL: u,
				Err: err,
			}
		}(i, rawURL)
	}

	wg.Wait()
	return results
}

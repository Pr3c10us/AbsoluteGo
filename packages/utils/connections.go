package utils

import (
	"database/sql"
	"fmt"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	_ "modernc.org/sqlite"
)

func NewS3Client(env *configs.EnvironmentVariables) *minio.Client {
	client, err := minio.New(env.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(env.S3.AccessKey, env.S3.SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("failed to create minio client: %v", err)
	}
	return client
}

func NewSQLClient(env *configs.EnvironmentVariables) *sql.DB {
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s", env.DatabasePath))
	if err != nil {
		log.Fatalf("failed to create minio client: %v", err)
	}
	return db
}

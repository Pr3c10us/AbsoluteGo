package adapters

import (
	"database/sql"
	ai2 "github.com/Pr3c10us/absolutego/internals/adapters/ai"
	book2 "github.com/Pr3c10us/absolutego/internals/adapters/book"
	storage2 "github.com/Pr3c10us/absolutego/internals/adapters/storage"
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/minio/minio-go/v7"
	"google.golang.org/genai"

	"github.com/Pr3c10us/absolutego/packages/configs"
)

type AdapterDependencies struct {
	EnvironmentVariables *configs.EnvironmentVariables
	GoogleGenAIClient    *genai.Client
	S3Client             *minio.Client
	DB                   *sql.DB
}

type Adapters struct {
	EnvironmentVariables *configs.EnvironmentVariables
	AiImplementation     ai.Interface
	StorageRepository    storage.Interface
	BookImplementation   book.Interface
}

func NewAdapters(dependencies AdapterDependencies) *Adapters {
	return &Adapters{
		EnvironmentVariables: dependencies.EnvironmentVariables,
		AiImplementation:     ai2.NewGoogleAI(dependencies.GoogleGenAIClient, dependencies.EnvironmentVariables.Gemini),
		StorageRepository:    storage2.NewMinioStorageRepository(dependencies.S3Client),
		BookImplementation:   book2.NewBookImplementation(dependencies.DB),
	}
}

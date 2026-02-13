package utils

import (
	"context"
	"database/sql"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"

	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/genai"
	_ "modernc.org/sqlite"
)

func NewS3Client(env *configs.EnvironmentVariables) *minio.Client {
	client, err := minio.New(env.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(env.S3.AccessKey, env.S3.SecretAccessKey, ""),
		Secure: false,
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

func NewGoogleGenAIClient(env *configs.EnvironmentVariables) *genai.Client {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: env.Gemini.APIKEY,
	})
	if err != nil {
		log.Fatalf("failed to create Google GenAI client: %v", err)
	}
	return client
}

func NewAMQChannel(conn *amqp.Connection) *amqp.Channel {
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel %v", err)
	}
	return channel
}

func NewAMQConnection(env *configs.EnvironmentVariables) *amqp.Connection {
	conn, err := amqp.Dial(env.AMQConnectionString)
	if err != nil {
		log.Printf("failed to connect to rabbitmq: %v", err)
		panic("failed to connect to rabbitmq")

	}
	return conn
}

package adapters

import (
	"database/sql"
	event2 "github.com/Pr3c10us/absolutego/internals/adapters/event"
	queue2 "github.com/Pr3c10us/absolutego/internals/adapters/queue"
	"github.com/Pr3c10us/absolutego/internals/domains/event"
	"github.com/Pr3c10us/absolutego/internals/domains/queue"
	amqp "github.com/rabbitmq/amqp091-go"

	ai2 "github.com/Pr3c10us/absolutego/internals/adapters/ai"
	book2 "github.com/Pr3c10us/absolutego/internals/adapters/book"
	script2 "github.com/Pr3c10us/absolutego/internals/adapters/script"
	storage2 "github.com/Pr3c10us/absolutego/internals/adapters/storage"
	"github.com/Pr3c10us/absolutego/internals/domains/ai"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
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
	AMQP                 *amqp.Channel
}

type Adapters struct {
	Dependencies          AdapterDependencies
	EnvironmentVariables  *configs.EnvironmentVariables
	AiImplementation      ai.Interface
	StorageImplementation storage.Interface
	BookImplementation    book.Interface
	ScriptImplementation  script.Interface
	QueueImplementation   queue.Interface
	EventImplementation   event.Interface
}

func NewAdapters(dependencies AdapterDependencies) *Adapters {
	return &Adapters{
		Dependencies:          dependencies,
		EnvironmentVariables:  dependencies.EnvironmentVariables,
		AiImplementation:      ai2.NewGoogleAI(dependencies.GoogleGenAIClient, dependencies.EnvironmentVariables.Gemini),
		StorageImplementation: storage2.NewMinioStorageImplementation(dependencies.S3Client),
		BookImplementation:    book2.NewBookImplementation(dependencies.DB),
		ScriptImplementation:  script2.NewScriptImplementation(dependencies.DB),
		QueueImplementation:   queue2.NewAMQImplementation(dependencies.AMQP, dependencies.EnvironmentVariables),
		EventImplementation:   event2.NewEventImplementation(dependencies.DB),
	}
}

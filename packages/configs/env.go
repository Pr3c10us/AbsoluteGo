package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type S3Credentials struct {
	Endpoint        string
	Region          string
	AccessKey       string
	SecretAccessKey string
}

type Buckets struct {
	ComicBucket string
	PageBucket  string
	PanelBucket string
	AudioBucket string
	VideoBucket string
	VabBucket   string
}

type GeminiConfig struct {
	APIKEY    string
	Model     string
	FastModel string
	LiveModel string
}

type EnvironmentVariables struct {
	Port                string
	DatabasePath        string
	AllowedOrigins      []string
	S3                  *S3Credentials
	Buckets             *Buckets
	Gemini              *GeminiConfig
	AMQConnectionString string
}

func loadEnv() {
	rootPath := GetRootPath()
	err := godotenv.Load(rootPath + `/.env`)

	if err != nil {
		log.Printf("Error loading .env file, falling back to system environment variables: %v", err)
	}
}

func LoadEnvironment() *EnvironmentVariables {
	loadEnv()
	return &EnvironmentVariables{
		Port:           getEnv("PORT", ":5000"),
		DatabasePath:   getEnv("DATABASE_PATH", "./database.db"),
		AllowedOrigins: strings.Split(getEnvOrError("ALLOWED_ORIGINS"), ","),
		S3: &S3Credentials{
			Endpoint:        getEnvOrError("S3_Endpoint"),
			Region:          getEnv("S3_REGION", ""),
			AccessKey:       getEnvOrError("S3_ACCESS_KEY"),
			SecretAccessKey: getEnvOrError("S3_SECRET_ACCESS_KEY"),
		},
		Buckets: &Buckets{
			ComicBucket: getEnvOrError("COMIC_BUCKET"),
			PageBucket:  getEnvOrError("PAGE_BUCKET"),
			PanelBucket: getEnvOrError("PANEL_BUCKET"),
			AudioBucket: getEnvOrError("AUDIOS_BUCKET"),
			VideoBucket: getEnvOrError("VIDEOS_BUCKET"),
			VabBucket:   getEnvOrError("VABS_BUCKET"),
		},
		Gemini: &GeminiConfig{
			APIKEY:    getEnvOrError("GEMINI_API_KEY"),
			Model:     getEnvOrError("GEMINI_MODEL"),
			FastModel: getEnvOrError("GEMINI_FAST_MODEL"),
			LiveModel: getEnvOrError("GEMINI_LIVE_MODEL"),
		},
		AMQConnectionString: getEnv("AMQ_CONNECTION_STRING", "amqp://guest:guest@localhost:5672/"),
	}
}

func getEnvOrError(key string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	panic("Environment variable " + key + " not set")
}

func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return fallback
}

func getEnvIntOrError(key string) int {
	value, exists := os.LookupEnv(key)
	if exists {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			log.Panicf("Environment variable \"%v\" not set properly", key)
		}
		return valueInt
	}
	errMsg := fmt.Sprintf("Environment variable \"%v\" not set", key)
	panic(errMsg)
}

func getEnvAsInt(key string, fallback int) int {
	value, exist := os.LookupEnv(key)
	if exist {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			log.Panicf("Environment variable \"%v\" not set properly", key)
		}
		return valueInt
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	value, exist := os.LookupEnv(key)
	if exist {
		valueBool, err := strconv.ParseBool(value)
		if err != nil {
			log.Panicf("Environment variable \"%v\" not set properly", key)
		}
		return valueBool
	}
	return fallback
}

package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	S3Url               string
	MongoHost           string
	MongoDatabase       string
	MongoCollectionName string
	WebSocketUrl        string
	AssetsPath          string
}

func getEnvValue(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		msg := fmt.Sprintf("not found %s", key)
		return "", errors.New(msg)
	}
	return value, nil
}

func NewConfig(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	s3Url, err := getEnvValue("S3_URL")
	if err != nil {
		return nil, err
	}

	mongoHost, err := getEnvValue("MONGO_HOST")
	if err != nil {
		return nil, err
	}

	mongoDatabase, err := getEnvValue("MONGO_DATABASE")
	if err != nil {
		return nil, err
	}

	mongoCollectionName, err := getEnvValue("MONGO_COLLECTION_NAME")
	if err != nil {
		return nil, err
	}

	webSocketUrl, err := getEnvValue("WEBSOCKET_URL")
	if err != nil {
		return nil, err
	}

	assetsPath, err := getEnvValue("ASSETS_PATH")
	if err != nil {
		return nil, err
	}

	return &Config{
		S3Url:               s3Url,
		MongoHost:           mongoHost,
		MongoDatabase:       mongoDatabase,
		MongoCollectionName: mongoCollectionName,
		WebSocketUrl:        webSocketUrl,
		AssetsPath:          assetsPath,
	}, nil

}

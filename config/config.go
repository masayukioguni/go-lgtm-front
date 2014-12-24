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
}

func getValue(key string) (string, error) {
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

	s3Url, err := getValue("S3_URL")
	if err != nil {
		return nil, err
	}

	mongoHost, err := getValue("MONGO_HOST")
	if err != nil {
		return nil, err
	}

	mongoDatabase, err := getValue("MONGO_DATABASE")
	if err != nil {
		return nil, err
	}

	mongoCollectionName, err := getValue("MONGO_COLLECTION_NAME")
	if err != nil {
		return nil, err
	}

	return &Config{
		S3Url:               s3Url,
		MongoHost:           mongoHost,
		MongoDatabase:       mongoDatabase,
		MongoCollectionName: mongoCollectionName,
	}, nil

}

package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	S3Url string
}

func NewConfig(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	s3Url := os.Getenv("S3_URL")
	if s3Url == "" {
		return nil, errors.New("not found S3_URL")
	}

	return &Config{S3Url: s3Url}, nil

}

package config

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfig_NewConfig(t *testing.T) {
	testPath := filepath.Join("./test-fixtures", ".env_test")

	c, _ := NewConfig(testPath)

	wantS3Url := "test"

	if !reflect.DeepEqual(c.S3Url, wantS3Url) {
		t.Errorf("TestConfig_NewConfig returned %+v, want %+v", c.S3Url, wantS3Url)
	}
}

func TestConfig_NewConfig_NoEnv(t *testing.T) {
	_, err := NewConfig("")

	if reflect.DeepEqual(err, nil) {
		t.Errorf("TestConfig_NewConfig_NoEnv returned %+v", err)
	}
}

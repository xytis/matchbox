package config

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// FileStoreConfig stores filesystem store configuration
type FileStoreConfig struct {
	Root string

	Logger *logrus.Logger
}

// NewFileStoreConfig creates an empty config
func NewFileStoreConfig(logger *logrus.Logger) *FileStoreConfig {
	return &FileStoreConfig{
		Logger: logger,
	}
}

// Validate checks if given config is viable
func (c *FileStoreConfig) Validate() error {
	if c.Root == "" {
		return errors.New("root path not provided")
	}
	if finfo, err := os.Stat(c.Root); err != nil || !finfo.IsDir() {
		return errors.Errorf(`root path "%s" is not a valid directory`, c.Root)
	}
	return nil
}

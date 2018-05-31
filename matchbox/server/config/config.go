package config

import (
	storage "github.com/coreos/matchbox/matchbox/storage/config"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	StoreBackendFile = "filesystem"
	StoreBackendEtcd = "etcd"
)

// Config configures a server implementation.
type Config struct {
	StoreBackend    string
	EtcdStoreConfig *storage.EtcdStoreConfig
	FileStoreConfig *storage.FileStoreConfig
}

// NewConfig creates an empty config
func NewConfig(logger *logrus.Logger) *Config {
	return &Config{
		EtcdStoreConfig: storage.NewEtcdStoreConfig(logger),
		FileStoreConfig: storage.NewFileStoreConfig(logger),
	}
}

// Validate performs configuration validation
func (c *Config) Validate() error {
	switch c.StoreBackend {
	case StoreBackendFile:
		if c.FileStoreConfig == nil {
			return errors.New(`unexpected empty configuration struct`)
		}
		return c.FileStoreConfig.Validate()
	case StoreBackendEtcd:
		if c.EtcdStoreConfig == nil {
			return errors.New(`unexpected empty configuration struct`)
		}
		return c.EtcdStoreConfig.Validate()
	default:
		return errors.Errorf(`invalid storage type "%v"`, c.StoreBackend)
	}
}

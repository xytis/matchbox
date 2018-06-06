package config

import (
	"github.com/pkg/errors"
)

const (
	//StoreBackendFile defines a string selector for filesystem backed store
	StoreBackendFile = "filesystem"
	//StoreBackendEtcd defines a string selector for etcd backed store
	StoreBackendEtcd = "etcd"
)

// Config configures a server implementation.
type Config struct {
	StoreBackend    string
	EtcdStoreConfig *EtcdStoreConfig
	FileStoreConfig *FileStoreConfig
}

// NewConfig creates an empty config
func NewConfig() *Config {
	return &Config{
		EtcdStoreConfig: NewEtcdStoreConfig(),
		FileStoreConfig: NewFileStoreConfig(),
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

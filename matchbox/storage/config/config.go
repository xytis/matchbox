package config

import (
	"github.com/pkg/errors"
)

// Available backend types
var (
	StoreBackendFile   = "filesystem"
	StoreBackendEtcd   = "etcd"
	StoreBackendMemory = "memory"
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
	case StoreBackendMemory:
		return nil
	default:
		return errors.Errorf(`invalid storage type "%v"`, c.StoreBackend)
	}
}

package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew_MissingEndpoints(t *testing.T) {
	cfg := &Config{
		Endpoint: "",
	}
	client, err := New(cfg)
	assert.Nil(t, client)
	assert.Error(t, err)
}

// gRPC expects host:port with no scheme (e.g. matchbox.example.com:8081)
func TestNew_InvalidEndpoints(t *testing.T) {
	invalid := []string{
		"matchbox.example.com",
		"http://matchbox.example.com:8081",
		"https://matchbox.example.com:8081",
	}

	for _, endpoint := range invalid {
		client, err := New(&Config{
			Endpoint: endpoint,
		})
		assert.Nil(t, client)
		assert.Error(t, err)
	}
}

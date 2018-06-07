package rpc

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/coreos/matchbox/matchbox/rpc/rpcpb"
	"github.com/coreos/matchbox/matchbox/server"
)

// Config prepares rpc server
type Config struct {
	Core server.Server
	TLS  *tls.Config
}

// NewServer wraps the matchbox Server to return a new gRPC Server.
func NewServer(config *Config) *grpc.Server {
	var opts []grpc.ServerOption
	if config.TLS != nil {
		// Add TLS Credentials as a ServerOption for server connections.
		opts = append(opts, grpc.Creds(credentials.NewTLS(config.TLS)))
	}

	grpcServer := grpc.NewServer(opts...)
	rpcpb.RegisterGroupsServer(grpcServer, newGroupServer(config.Core))
	rpcpb.RegisterProfilesServer(grpcServer, newProfileServer(config.Core))
	rpcpb.RegisterTemplatesServer(grpcServer, newTemplateServer(config.Core))
	rpcpb.RegisterSelectServer(grpcServer, newSelectServer(config.Core))
	rpcpb.RegisterVersionServer(grpcServer, newVersionServer(config.Core))
	return grpcServer
}

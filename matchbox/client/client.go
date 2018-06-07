package client

import (
	"crypto/tls"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/coreos/matchbox/matchbox/rpc/rpcpb"
	"github.com/pkg/errors"
)

// Config configures a Client.
type Config struct {
	// Endpoint to connect to
	Endpoint string
	// DialTimeout is the timeout for dialing a client connection
	DialTimeout time.Duration
	// Client TLS credentials
	TLS *tls.Config
}

// Client provides a matchbox client RPC session.
type Client struct {
	Groups    rpcpb.GroupsClient
	Profiles  rpcpb.ProfilesClient
	Templates rpcpb.TemplatesClient
	Select    rpcpb.SelectClient
	Version   rpcpb.VersionClient

	conn *grpc.ClientConn
}

// New creates a new Client from the given Config.
func New(config *Config) (*Client, error) {
	endpoint := config.Endpoint

	if len(endpoint) == 0 {
		return nil, errors.New("server endpoint not provided")
	}

	if _, _, err := net.SplitHostPort(endpoint); err != nil {
		return nil, errors.Wrap(err, "client: invalid host:port endpoint")
	}

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(config.DialTimeout),
	}
	if config.TLS != nil {
		creds := credentials.NewTLS(config.TLS)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "client: unable to dial")
	}

	client := &Client{
		conn:      conn,
		Groups:    rpcpb.NewGroupsClient(conn),
		Profiles:  rpcpb.NewProfilesClient(conn),
		Templates: rpcpb.NewTemplatesClient(conn),
		Select:    rpcpb.NewSelectClient(conn),
		Version:   rpcpb.NewVersionClient(conn),
	}
	return client, nil
}

// Close closes the client's connections.
func (c *Client) Close() error {
	return c.conn.Close()
}

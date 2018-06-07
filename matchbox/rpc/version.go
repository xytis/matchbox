package rpc

import (
	"golang.org/x/net/context"

	"github.com/coreos/matchbox/matchbox/rpc/rpcpb"
	"github.com/coreos/matchbox/matchbox/server"
	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/version"
)

// versionServer wraps a matchbox Server to be suitable for gRPC registration.
type versionServer struct {
	srv server.Server
}

func newVersionServer(s server.Server) rpcpb.VersionServer {
	return &versionServer{
		srv: s,
	}
}

func (s *versionServer) VersionReport(ctx context.Context, req *pb.VersionReportRequest) (*pb.VersionReportResponse, error) {
	return &pb.VersionReportResponse{Version: version.Version}, grpcError(nil)
}

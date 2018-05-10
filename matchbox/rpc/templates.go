package rpc

import (
	"golang.org/x/net/context"

	"github.com/coreos/matchbox/matchbox/rpc/rpcpb"
	"github.com/coreos/matchbox/matchbox/server"
	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// templateServer takes a matchbox Server and implements a gRPC TemplatesServer.
type templateServer struct {
	srv server.Server
}

func newTemplateServer(s server.Server) rpcpb.TemplatesServer {
	return &templateServer{
		srv: s,
	}
}

func (s *templateServer) TemplatePut(ctx context.Context, req *pb.TemplatePutRequest) (*pb.TemplatePutResponse, error) {
	_, err := s.srv.TemplatePut(ctx, req)
	return &pb.TemplatePutResponse{}, grpcError(err)
}

func (s *templateServer) TemplateGet(ctx context.Context, req *pb.TemplateGetRequest) (*pb.TemplateGetResponse, error) {
	template, err := s.srv.TemplateGet(ctx, req)
	return &pb.TemplateGetResponse{Template: template}, grpcError(err)
}

func (s *templateServer) TemplateDelete(ctx context.Context, req *pb.TemplateDeleteRequest) (*pb.TemplateDeleteResponse, error) {
	err := s.srv.TemplateDelete(ctx, req)
	return &pb.TemplateDeleteResponse{}, grpcError(err)
}

package server

import (
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"github.com/taeho-io/auth"
	tid "github.com/taeho-io/go-taeho/id"
	"github.com/taeho-io/id"
	"github.com/taeho-io/id/server/handler"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type IDServer struct {
	id.IdServer

	cfg Config
	id  tid.ID
}

func New(cfg Config) (*IDServer, error) {
	return &IDServer{
		cfg: cfg,
		id:  tid.New(),
	}, nil
}

func Mock() *IDServer {
	s, _ := New(MockConfig())
	return s
}

func (s *IDServer) Config() Config {
	return s.cfg
}

func (s *IDServer) ID() tid.ID {
	return s.id
}

func (s *IDServer) RegisterServer(srv *grpc.Server) {
	id.RegisterIdServer(srv, s)
}

func (s *IDServer) Create(ctx context.Context, req *id.NewRequest) (*id.NewResponse, error) {
	return handler.New(s.ID())(ctx, req)
}

func Serve(addr string, cfg Config) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			auth.TokenUnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			grpc_recovery.UnaryServerInterceptor(),
		),
	)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	idServer, err := New(cfg)
	if err != nil {
		return err
	}
	id.RegisterIdServer(grpcServer, idServer)

	reflection.Register(grpcServer)
	return grpcServer.Serve(lis)
}

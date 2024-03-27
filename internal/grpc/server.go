package internalgrpc

//go:generate protoc -I ../../../api --go_out=../../../ --go-grpc_out=../../../ // --grpc-gateway_out=../../../ --openapiv2_out=../../../ ../../../api/mbb/mbb.proto //nolint:lll

import (
	"context"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/mmart-pro/mart-brute-blocker/internal/grpc/pb"
	"github.com/mmart-pro/mart-brute-blocker/internal/model"
)

type Server struct {
	addr       string
	logger     Logger
	mbbService MbbService
	server     *grpc.Server
}

type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}

type MbbService interface {
	Allow(ctx context.Context, req model.Subnet) error
	Deny(ctx context.Context, req model.Subnet) error
	Remove(ctx context.Context, req model.Subnet) error
	Exists(ctx context.Context, req model.Subnet) (model.ListType, error)
	Contains(ctx context.Context, req model.IPAddr) (model.ListType, error)
	ClearBucket(_ context.Context, _ model.IPAddr, _ string) error
	Check(ctx context.Context, ip model.IPAddr, _, _ string) (bool, error)
}

func NewServer(addr string, logger Logger, mbbService MbbService) *Server {
	return &Server{
		addr:       addr,
		logger:     logger,
		mbbService: mbbService,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer(grpc.ChainUnaryInterceptor(s.loggerInterceptor()))
	pb.RegisterMBBServiceServer(s.server, NewService(s.mbbService))

	s.logger.Debugf("starting grpc on %s", s.addr)

	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}

func (s *Server) loggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()

		meta, ok := metadata.FromIncomingContext(ctx)
		agent := ""
		if ok {
			agent = strings.Join(meta.Get("user-agent"), " ")
		}

		resp, err := handler(ctx, req)
		// Логирование ошибки, если она произошла
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}

		// Получение RemoteAddr из контекста
		p, ok := peer.FromContext(ctx)
		var remoteAddr string
		if ok {
			remoteAddr = p.Addr.String()
		}

		s.logger.Debugf("addr: %s method: %s duration: %v user-agent: [%s] err: '%v'",
			remoteAddr, info.FullMethod, time.Since(startTime), agent, errMsg)

		return resp, err
	}
}

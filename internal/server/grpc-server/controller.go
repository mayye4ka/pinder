package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	public_api "github.com/mayye4ka/pinder-api/api/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	metadataContextKey      = "authorization"
	authorizationTrimPrefix = "Bearer "
	userIdContextKey        = "user_id"
)

type ServerCtrl struct {
	server     *Server
	port       int
	grpcServer *grpc.Server
}

func New(svc Service, auth Authenticator, port int) *ServerCtrl {
	return &ServerCtrl{
		server: &Server{
			service: svc,
			auth:    auth,
		},
		port: port,
	}
}

func (c *ServerCtrl) Start(ctx context.Context) error {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(c.authInterceptor),
		grpc.Creds(insecure.NewCredentials()),
	}
	c.grpcServer = grpc.NewServer(opts...)
	public_api.RegisterPinderServer(c.grpcServer, c.server)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.port))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-ctx.Done()
		c.grpcServer.GracefulStop()
	}()
	return c.grpcServer.Serve(lis)
}

func (c *ServerCtrl) Stop(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		c.grpcServer.Stop()
	}()
	c.grpcServer.GracefulStop()
	return nil
}

func (c *ServerCtrl) authInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	userId := c.getUserIdFromIncomingContext(ctx)
	ctx = context.WithValue(ctx, userIdContextKey, userId)
	return handler(ctx, req)
}

func (c *ServerCtrl) getUserIdFromIncomingContext(ctx context.Context) uint64 {
	token := c.getTokenFromIncomingContext(ctx)
	if token == "" {
		return 0
	}
	userId, err := c.server.auth.UnpackToken(ctx, token)
	if err != nil {
		return 0
	}
	return userId
}

func (c *ServerCtrl) getTokenFromIncomingContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if len(md[metadataContextKey]) != 1 {
		return ""
	}
	return strings.TrimPrefix(md[metadataContextKey][0], authorizationTrimPrefix)
}

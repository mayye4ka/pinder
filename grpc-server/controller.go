package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	public_api "github.com/mayye4ka/pinder-api/public_api/go"

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
	server *Server
}

func New(svc Service, auth Authenticator) *ServerCtrl {
	return &ServerCtrl{
		server: &Server{
			service: svc,
			auth:    auth,
		},
	}
}

func (c *ServerCtrl) Start(port int) error {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(c.authInterceptor),
		grpc.Creds(insecure.NewCredentials()),
	}
	srv := grpc.NewServer(opts...)
	public_api.RegisterPinderServer(srv, c.server)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	return srv.Serve(lis)
}

func (c *ServerCtrl) authInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	t1 := time.Now()
	userId := c.getUserIdFromIncomingContext(ctx)
	ctx = context.WithValue(ctx, userIdContextKey, userId)
	t2 := time.Now()
	resp, err := handler(ctx, req)
	t3 := time.Now()
	log.Printf("auth actions took %v, handler took %v", t2.Sub(t1), t3.Sub(t2))
	return resp, err
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

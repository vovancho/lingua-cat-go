package auth

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func (s *AuthService) AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Missing metadata")
	}

	authHeaders, ok := md["authorization"]
	if !ok || len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Missing Authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeaders[0], "Bearer ")
	if tokenStr == authHeaders[0] {
		return nil, status.Error(codes.Unauthenticated, "Invalid Authorization header format")
	}

	userID, err := s.VerifyToken(tokenStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	ctx = context.WithValue(ctx, "userID", userID)
	return handler(ctx, req)
}

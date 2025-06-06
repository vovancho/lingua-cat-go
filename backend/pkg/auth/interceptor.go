package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

	ctx = s.withUserID(ctx, userID)
	ctx = s.withJWTToken(ctx, tokenStr)

	return handler(ctx, req)
}

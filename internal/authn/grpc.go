package authn

import (
	"context"

	"github.com/Di-nis/shortener-url/internal/constants"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const headerAuthorization = "authorization"

// Interceptor - аутентификация пользователя.
func Interceptor(JWTSecret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var (
			token, userID, sessionID string
			err                      error
		)

		mdReq, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := mdReq.Get(headerAuthorization)
			if len(values) > 0 {
				token = values[0]
			}
		}

		if token == "" {
			userID = GenerateUserID()
			sessionID = GenerateSessionID()
			token, err = BuildJWTString(JWTSecret, userID, sessionID)
			if err != nil {
				return nil, status.Error(codes.Internal, "internal error")
			}
		} else {
			claims, isTokenValid := GetClaims(token, JWTSecret)
			if !isTokenValid {
				return nil, status.Error(codes.Unauthenticated, "token not valid")
			}
			userID = claims.UserID
			sessionID = claims.SID
			if sessionID == "" {
				return nil, status.Error(codes.Unauthenticated, "sessionID not valid")
			}
		}

		mdOut := metadata.Pairs(
			headerAuthorization, token,
		)

		// отправка headers клиенту
		err = grpc.SendHeader(ctx, mdOut)
		if err != nil {
			return nil, status.Error(codes.Internal, "internal error")
		}

		ctx = context.WithValue(ctx, constants.UserIDKey, userID)
		return handler(ctx, req)
	}
}

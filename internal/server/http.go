package server

import (
	"context"
	"errors"

	userPb "kratos-starter/api/v1/user"
	"kratos-starter/internal/conf"
	"kratos-starter/internal/constant"
	"kratos-starter/internal/service"

	kErrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c *conf.Server,
	logger log.Logger,
	rds redis.UniversalClient,
	user *service.UserService,
) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			validate.Validator(),
			selector.Server(
				jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
					return []byte(constant.JwtSignKey), nil
				},
					jwt.WithSigningMethod(jwtv5.SigningMethodHS256),
					jwt.WithClaims(func() jwtv5.Claims {
						return &jwtv5.RegisteredClaims{}
					}),
				),
				rdsTokenValidator(rds),
			).Match(NewWhiteListMatcher()).Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	userPb.RegisterUserServiceHTTPServer(srv, user)
	return srv
}

// NewWhiteListMatcher 路由白名单
func NewWhiteListMatcher() selector.MatchFunc {
	whitelist := make(map[string]struct{})
	whitelist[userPb.OperationUserServiceLogin] = struct{}{}

	return func(ctx context.Context, operation string) bool {
		if _, ok := whitelist[operation]; ok {
			return false
		}
		return true
	}
}

func rdsTokenValidator(rds redis.UniversalClient) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			claims, ok := jwt.FromContext(ctx)
			if !ok {
				return nil, jwt.ErrMissingJwtToken
			}

			registerClams, ok := claims.(*jwtv5.RegisteredClaims)
			if !ok {
				return nil, jwt.ErrMissingJwtToken
			}

			clamsIssuer, err := registerClams.GetIssuer()
			if err != nil {
				return nil, jwt.ErrTokenInvalid
			}

			issuer, err := rds.Get(ctx, constant.GenUserTokenKey(registerClams.ID)).Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					return err, jwt.ErrTokenInvalid
				}
				return nil, kErrors.InternalServer("INTERNAL ERROR", err.Error())
			}

			if issuer != clamsIssuer || clamsIssuer == "" {
				return nil, jwt.ErrTokenInvalid
			}

			return handler(ctx, req)
		}
	}
}

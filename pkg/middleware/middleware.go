package middleware

import (
	"crypto/rsa"

	ginJwt "github.com/appleboy/gin-jwt/v2"
	"github.com/rs/zerolog"

	"employee-management-system/model"
	"employee-management-system/pkg/environment"
	"employee-management-system/pkg/helper"
	"employee-management-system/storage"
)

const (
	// RequestBodyInContext context key holder
	RequestBodyInContext = "request_body_in_context"
	// RequestUserIDInContext context for API interceptor system user_id
	RequestUserIDInContext = "request_user_id_in_context"
	packageName            = "middleware"
)

type (

	// Middleware object
	Middleware struct {
		logger          zerolog.Logger
		env             environment.Env
		employeeStorage storage.EmployeeDatabase
		userStorage     storage.UserDatabase
		jwt             *ginJwt.GinJWTMiddleware
		pKey            *rsa.PrivateKey
	}
)

// NewMiddleware new instance of our custom ginJwt middleware
func NewMiddleware(z zerolog.Logger, env environment.Env, s *storage.Storage) *Middleware {
	mWare, _ := jwtMiddleware(env.Get("SIGNING_SECRET_KEY"))
	l := z.With().Str(helper.LogStrKeyModule, packageName).Logger()

	return &Middleware{
		logger:          l,
		env:             env,
		userStorage:     *storage.NewUser(s),
		employeeStorage: *storage.NewEmployee(s),
		jwt:             mWare,
	}
}

func jwtMiddleware(secretKey string) (*ginJwt.GinJWTMiddleware, error) {
	return ginJwt.New(&ginJwt.GinJWTMiddleware{
		Realm:      realm,
		Key:        []byte(secretKey),
		MaxRefresh: maxRefresh,
		PayloadFunc: func(data interface{}) ginJwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return ginJwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return ginJwt.MapClaims{}
		},
		IdentityKey: identityKey,
		Timeout:     accessTimeout,
	})
}

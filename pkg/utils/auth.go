package utils

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

func GetUid(ctx context.Context) (uint, error) {
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		return 0, fmt.Errorf("no claims")
	}

	registerClams, ok := claims.(*jwtv5.RegisteredClaims)
	if !ok {
		return 0, fmt.Errorf("no register claims")
	}

	uid, err := strconv.ParseUint(registerClams.ID, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(uid), nil
}

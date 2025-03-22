package mid

import (
	"context"
	"errors"

	"github.com/ardanlabs/service/app/api/auth"
	"github.com/ardanlabs/service/app/api/errs"
)

var ErrInvalidID = errors.New("ID is not in its proper form")

func Authorize(ctx context.Context, ath *auth.Auth, rule string, handler Handler) error {
	userID, err := GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	claims := GetClaims(ctx)
	if err := ath.Authorize(ctx, claims, userID, rule); err != nil {
		return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)
	}

	return handler(ctx)
}

package mid

import (
	"context"
	"fmt"

	"github.com/ardanlabs/service/app/api/auth"
	"github.com/ardanlabs/service/app/api/errs"
	"github.com/google/uuid"
)

func Bearer(ctx context.Context, ath *auth.Auth, authorization string, handler Handler) error {
	claims, err := ath.Authenticate(ctx, authorization)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	if claims.Subject == "" {
		return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims")
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return errs.New(errs.Unauthenticated, fmt.Errorf("parsing subject: %w", err))
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return handler(ctx)
}

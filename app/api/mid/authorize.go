package mid

import (
	"context"
	"errors"

	"github.com/ardanlabs/service/app/api/authclient"
	"github.com/ardanlabs/service/app/api/errs"
	"github.com/ardanlabs/service/foundation/logger"
)

var ErrInvalidID = errors.New("ID is not in its proper form")

func Authorize(ctx context.Context, log *logger.Logger, client *authclient.Client, rule string, handler Handler) error {
	userID, err := GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	auth := authclient.Authorize{
		Claims: GetClaims(ctx),
		UserID: userID,
		Rule:   rule,
	}

	if err := client.Authorize(ctx, auth); err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	return handler(ctx)
}

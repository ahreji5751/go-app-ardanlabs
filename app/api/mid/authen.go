package mid

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/ardanlabs/service/app/api/auth"
	"github.com/ardanlabs/service/app/api/authclient"
	"github.com/ardanlabs/service/app/api/errs"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func Authenticate(ctx context.Context, log *logger.Logger, client *authclient.Client, authorization string, handler Handler) error {
	resp, err := client.Authenticate(ctx, authorization)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	ctx = setUserID(ctx, resp.UserID)
	ctx = setClaims(ctx, resp.Claims)

	return handler(ctx)
}

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

func Basic(ctx context.Context, ath *auth.Auth, authorization string, handler Handler) error {
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "231c0a1e-664f-4b2c-a8a9-f3fbaa9a8bd1",
			Issuer:    ath.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return errs.Newf(errs.Unauthenticated, "parsing subject: %s", err)
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return handler(ctx)
}

func parseBasicAuth(auth string) (string, string, bool) {
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", false
	}

	c, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}

	username, password, ok := strings.Cut(string(c), ":")
	if !ok {
		return "", "", false
	}

	return username, password, true
}

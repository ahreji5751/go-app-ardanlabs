package authapi

import (
	"context"
	"net/http"

	"github.com/ardanlabs/service/app/api/auth"
	"github.com/ardanlabs/service/app/api/errs"
	"github.com/ardanlabs/service/app/api/mid"
	"github.com/ardanlabs/service/foundation/web"
	"github.com/google/uuid"
)

type api struct {
	auth *auth.Auth
}

func newAPI(ath *auth.Auth) *api {
	return &api{
		auth: ath,
	}
}

func (api *api) token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	kid := web.Param(r, "kid")
	if kid == "" {
		return errs.Newf(errs.FailedPrecondition, "missing kid")
	}

	claims := mid.GetClaims(ctx)

	tkn, err := api.auth.GenerateToken(kid, claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	token := struct {
		Token string `json:"token"`
	}{
		Token: tkn,
	}

	return web.Respond(ctx, w, token, http.StatusOK)
}

func (api *api) authenticate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	resp := struct {
		UserID uuid.UUID
		Claims auth.Claims
	}{
		UserID: userID,
		Claims: mid.GetClaims(ctx),
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

func (api *api) authorize(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var ath struct {
		Claims auth.Claims
		UserID uuid.UUID
		Rule   string
	}
	if err := web.Decode(r, &ath); err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	if err := api.auth.Authorize(ctx, ath.Claims, ath.UserID, ath.Rule); err != nil {
		return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", ath.Claims.Roles, ath.Rule, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

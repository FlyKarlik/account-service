package server

import (
	"account-service/internal/models"
	"context"
	"github.com/hashicorp/go-hclog"
	protos "protos/account"
)

// CheckAccess check access token
func (a *AccountService) CheckAccess(ctx context.Context, rr *protos.CheckAccessRequest) (*protos.CheckAccessResponse, error) {
	log := hclog.Default()

	ctx, span := a.trace.Start(ctx, "CheckAccess")
	defer span.End()

	accessToken := rr.GetToken()

	tok, err := a.tokenSrv.ParseJWT(
		ctx,
		accessToken,
	)
	if err != nil {
		log.Error("[server.CheckAccess] a.tokenSrv.ParseJWT", "accessToken", accessToken, "error", err)
		return nil, models.UnauthenticatedAuthTokenError
	}

	if err := a.tokenSrv.Validate(tok, models.AccessToken); err != nil {
		log.Error("[server.CheckAccess] a.tokenSrv.Validate", "accessToken", accessToken, "error", err)
		return nil, models.UnauthenticatedAuthTokenError
	}

	user, err := a.db.GetUserByID(tok.Identity)
	if err != nil {
		log.Error("[server.CheckAccess] a.db.GetUserByID", "error", err)
		return nil, models.InternalError
	}
	if user == nil {
		return nil, models.UserNotFoundError
	}

	return &protos.CheckAccessResponse{
		UserId:       tok.Identity,
		Email:        user.Email,
		DepartmentId: user.DepartmentID,
		RoleId:       user.RoleID,
	}, nil
}

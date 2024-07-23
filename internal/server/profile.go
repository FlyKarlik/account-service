package server

import (
	"account-service/internal/models"
	"comet/utils"
	"context"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/types/known/emptypb"
	protos "protos/account"
)

// GetAccountInfo return account info
func (a *AccountService) GetAccountInfo(ctx context.Context, _ *emptypb.Empty) (*protos.GetAccountInfoResponse, error) {
	log := hclog.Default()

	tr := a.trace
	ctx, span := tr.Start(ctx, "GetAccuntInfo")
	defer span.End()

	accessToken, err := utils.GetAccessHeader(&ctx)
	if err != nil {
		log.Error("[server.GetAccountInfo] utils.GetAccessHeader", "error", err)
		return nil, models.InvalidAccessTokenError
	}

	tok, err := a.tokenSrv.ParseJWT(
		ctx,
		accessToken,
	)
	if err != nil {
		log.Error("[server.GetAccountInfo] a.tokenSrv.ParseJWT", "error", err)
		return nil, models.UnauthenticatedAccessTokenError
	}

	if err := a.tokenSrv.Validate(tok, models.AccessToken); err != nil {
		log.Error("[server.GetAccountInfo] a.tokenSrv.Validate", "error", err)
		return nil, models.UnauthenticatedAccessTokenError
	}

	user, err := a.db.GetUserByID(tok.Identity)
	if err != nil {
		log.Error("[server.GetAccountInfo] a.db.GetUserByID", "uid", tok.Identity, "error", err)
		return nil, models.InternalError
	}
	if user == nil {
		return nil, models.UserNotFoundError
	}

	return &protos.GetAccountInfoResponse{
		UserId:       user.ID,
		FirstName:    user.FirstName,
		SecondName:   user.LastName,
		Email:        user.Email,
		DepartmentId: user.DepartmentID,
		RoleId:       user.RoleID,
	}, nil
}

//func (a *AccountService) GetAccountsInfo(ctx context.Context, _ *emptypb.Empty) (*protos.GetAccountsInfoResponse, error) {
//
//}

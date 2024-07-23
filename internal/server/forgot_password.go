package server

import (
	"account-service/internal/models"
	"account-service/internal/validators"
	"comet/utils"
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/types/known/emptypb"
	protos "protos/account"
)

// ForgotPassword request forgot password for specific user
func (a *AccountService) ForgotPassword(ctx context.Context, _ *emptypb.Empty) (*protos.ForgotPasswordResponse, error) {
	log := hclog.Default()

	tr := a.trace
	ctx, span := tr.Start(ctx, "ForgotPassword")
	defer span.End()

	authToken, err := utils.GetAccessHeader(&ctx)
	if err != nil {
		log.Error("[server.ForgotPassword] utils.GetAccessHeader", "error", err)
		return nil, models.InvalidAccessTokenError
	}

	tok, err := a.tokenSrv.ParseJWT(
		ctx,
		authToken,
	)
	if err != nil {
		log.Error("[server.ForgotPassword] a.tokenSrv.ParseJWT", "registerToken", authToken, "error", err)
		return nil, models.UnauthenticatedAccessTokenError
	}

	errLoginToken := a.tokenSrv.Validate(tok, models.FirstLoginToken)
	errAuthToken := a.tokenSrv.Validate(tok, models.AuthToken)
	if errLoginToken != nil && errAuthToken != nil {
		log.Error("[server.ForgotPassword] a.tokenSrv.Validate", "registerToken", authToken, "errLoginToken", errLoginToken, "errAuthToken", errAuthToken)
		return nil, models.UnauthenticatedAccessTokenError
	}

	user, err := a.db.GetUserByID(tok.Identity)
	if err != nil {
		log.Error("[server.ForgotPassword] a.db.GetUserByID", "authToken", authToken, "userID", tok.Identity, "error", err)
		return nil, models.InternalError
	}
	if user == nil {
		return nil, models.UserNotFoundError
	}

	return &protos.ForgotPasswordResponse{
		ForgotToken: "",
		SecretEmail: utils.HideEmail(user.Email),
	}, nil
}

// ResetPassword reset password for specific user
func (a *AccountService) ResetPassword(ctx context.Context, rr *protos.ResetPasswordRequest) (*protos.ResetPasswordResponse, error) {
	log := hclog.Default()

	tr := a.trace
	ctx, span := tr.Start(ctx, "ResetPassword")
	defer span.End()

	resetToken := rr.GetResetPasswordToken()

	tok, err := a.tokenSrv.ParseJWT(
		ctx,
		resetToken,
	)
	if err != nil {
		log.Error("[server.ResetPassword] a.tokenSrv.ParseJWT", "error", err)
		return nil, models.InvalidResetPasswordError
	}

	if err := a.tokenSrv.Validate(tok, models.ResetPasswordToken); err != nil {
		log.Error("[server.ResetPassword] a.tokenSrv.Validate", "error", err)
		return nil, models.InvalidResetPasswordError
	}

	password := rr.GetNewPassword()
	err = validators.ValidatePassword(password)
	if err != nil {
		return nil, models.PasswordNotValidError(err)
	}

	user, err := a.db.ChangePasswordByID(tok.Identity, password)
	if err != nil {
		log.Error("[server.ResetPassword] a.db.ChangePasswordByID", "userID", tok.Identity, "error", err)
		return nil, models.InternalError
	}

	token, err := a.tokenSrv.NewJWT(
		ctx,
		models.FirstLoginToken,
		user.ID,
		user.Email,
		jwt.MapClaims{},
	)
	if err != nil {
		log.Error("[server.Authorize] a.tokenSrv.NewJWT", "userID", user.ID, "error", err)
		return nil, models.InternalError
	}

	//TODO Revoke all tokens of users
	a.tokenSrv.Revoke(tok)

	return &protos.ResetPasswordResponse{
		LoginToken: token.ToJWTString(),
	}, nil
}

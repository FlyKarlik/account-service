package server

import (
	"account-service/internal/models"
	"account-service/internal/validators"
	"comet/utils"
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	protos "protos/account"
	"time"
)

// LoginUser recurring login user
func (a *AccountService) LoginUser(ctx context.Context, rr *protos.LoginUserRequest) (*protos.LoginUserResponse, error) {
	log := hclog.Default()

	tr := a.trace
	ctx, span := tr.Start(ctx, "Login")
	defer span.End()

	email := rr.GetEmail()

	err := validators.ValidateEmail(email)
	if err != nil {
		log.Error("[server.LoginUser] validators.ValidateEmail", "error", err)
		return nil, models.EmailNotValidError
	}

	user, err := a.db.GetUserByEmail(email)
	if err != nil {
		log.Error("[server.LoginUser] a.db.GetUserByEmail", "error", err)
		return nil, models.UserNotFoundError
	}

	password := rr.GetPassword()
	if len(password) > 0 {
		log.Error("[server.LoginUser] password less then 1", "error")
		return nil, models.PasswordNotValidError(fmt.Errorf("password less then 1"))
	}

	if password != "" {
		isValid, err := utils.VerifyArgon(user.Password, password)
		if err != nil {
			log.Error("[server.LoginUser] utils.VerifyArgon", "userID", user.ID, "error", err)
			return nil, models.InternalError
		}

		if !isValid {
			return nil, models.NotMatchError
		}
	} else {
		return nil, models.BadRequestError
	}

	token, refresh, err := a.tokenSrv.CreateAccessJWT(
		ctx,
		user.ID,
		user.Email,
	)

	if err != nil {
		log.Error("[server.LoginUser] a.tokenSrv.CreateAccessJWT", "userID", user.ID, "error", err)
		return nil, models.InternalError
	}

	accessTokenExpirationTime, ok := models.Expirations["AccessToken"]

	if !ok {
		log.Error("[server.LoginUser] models.Expirations", "error", "accessToken expiration not found")
	}

	refreshTokenExpirationTime, ok := models.Expirations["RefreshToken"]

	if !ok {
		log.Error("[server.LoginUser] models.Expirations", "error", "refreshToken expiration not found")
	}

	return &protos.LoginUserResponse{
		AccessToken:           token.ToJWTString(),
		AccessTokenExpiredAt:  time.Now().Add(accessTokenExpirationTime).Unix(),
		RefreshToken:          refresh.ToJWTString(),
		RefreshTokenExpiredAt: time.Now().Add(refreshTokenExpirationTime).Unix(),
		User: &protos.GetAccountInfoResponse{
			UserId:       user.ID,
			FirstName:    user.FirstName,
			SecondName:   user.LastName,
			Email:        email,
			DepartmentId: user.DepartmentID,
			RoleId:       user.RoleID,
		},
	}, nil
}

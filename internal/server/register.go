package server

import (
	"account-service/internal/models"
	"account-service/internal/validators"
	"comet/utils"
	"context"
	"github.com/hashicorp/go-hclog"
	protos "protos/account"
	"strings"
)

func (a *AccountService) RegisterUser(ctx context.Context, rr *protos.RegisterUserRequest) (*protos.RegisterUserResponse, error) {
	log := hclog.Default()

	tr := a.trace
	_, span := tr.Start(ctx, "RegisterUser")
	defer span.End()

	registerToken, err := utils.GetAccessHeader(&ctx)
	if err != nil {
		log.Error("[server.RegisterUsername] utils.GetAccessHeader", "error", err)
		return nil, models.InvalidAuthTokenError
	}

	tok, err := a.tokenSrv.ParseJWT(
		ctx,
		registerToken,
	)
	if err != nil {
		log.Error("[server.RegisterUsername] a.tokenSrv.ParseJWT", "error", err)
		return nil, models.UnauthenticatedAuthTokenError
	}

	if err := a.tokenSrv.Validate(tok, models.RegisterToken); err != nil {
		log.Error("[server.RegisterUsername] a.tokenSrv.Validate", "error", err)
		return nil, models.UnauthenticatedAuthTokenError
	}

	email := strings.TrimSpace(strings.ToLower(rr.GetEmail()))
	err = validators.ValidateEmail(email)
	if err != nil {
		log.Error("[server.RegisterUser] validators.ValidateEmail", "userID", tok.Identity, "email", email, "error", err)
		return nil, models.UsernameNotValidError(err)
	}

	emailExist, err := a.db.EmailExist(email)
	if err != nil {
		log.Error("[server.RegisterUser] a.db.EmailExists", "userID", tok.Identity, "error", err)
		return nil, models.InternalError
	}
	if emailExist {
		return nil, models.AlreadyExistError
	}

	firstName := rr.GetFirstName()
	err = validators.ValidateFIO(firstName)
	if err != nil {
		log.Error("[server.RegisterUser] validators.ValidateFIO", "error", err)
		return nil, models.FioNotValidError
	}

	lastName := rr.GetSecondName()
	err = validators.ValidateFIO(lastName)
	if err != nil {
		log.Error("[server.RegisterUser] validators.ValidateFIO", "error", err)
		return nil, models.FioNotValidError
	}

	password := rr.GetPassword()
	err = validators.ValidatePassword(password)
	if err != nil {
		log.Error("[server.RegisterUser] validators.ValidatePassword", "error", err)
		return nil, models.PasswordNotValidError(err)
	}

	b, err := utils.HashArgon(password)
	if err != nil {
		log.Error("[server.RegisterUser] utils.HashArgon", err)
		return nil, models.InternalError
	}
	hashPassword := string(b)

	departmentID := rr.GetDepartmentId()
	department, err := a.db.GetUserDepartmentByID(departmentID)
	if err != nil {
		log.Error("[server.RegisterUser] a.db.GetUserDepartmentByID", "error", err)
		return nil, models.DepartmentNotFoundError
	}

	if department == nil {
		log.Error("[server.RegisterUser] a.db.GetUserDepartmentByID", "error", err)
		return nil, models.DepartmentNotFoundError
	}

	roleID := rr.GetRoleId()
	role, err := a.db.GetUserRoleByID(roleID)
	if err != nil {
		log.Error("[server.RegisterUser] a.db.GetUserRoleByID", "error", err)
		return nil, models.RoleNotFoundError
	}

	if role == nil {
		log.Error("[server.RegisterUser] a.db.GetUserRoleByID", "error", err)
		return nil, models.RoleNotFoundError
	}

	user, err := a.db.CreateUserIfNotExist(email, firstName, lastName, hashPassword, departmentID, roleID)
	if err != nil {
		log.Error("[server.RegisterUser] a.db.CreateUserIfNotExist", "error", err)
		return nil, models.InternalError
	}

	return &protos.RegisterUserResponse{
		Uuid:         user.ID,
		Email:        user.Email,
		FirstName:    firstName,
		SecondName:   lastName,
		DepartmentId: departmentID,
		RoleId:       roleID,
	}, nil
}

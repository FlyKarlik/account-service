package models

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BadRequestError bad request
var BadRequestError = status.Errorf(
	codes.InvalidArgument,
	"Bad request error",
)

// InternalError internal error
var InternalError = status.Errorf(
	codes.Internal,
	"Internal server error",
)

// UserNotFoundError user is not found
var UserNotFoundError = status.Errorf(
	codes.NotFound,
	"User not found",
)

// DepartmentNotExist department not exist
var DepartmentNotFoundError = status.Errorf(
	codes.NotFound,
	"Department not exist",
)

// RoleNotExist department not exist
var RoleNotFoundError = status.Errorf(
	codes.NotFound,
	"Role not exist",
)

// InvalidEmailError invalid email address
var InvalidEmailError = status.Errorf(
	codes.InvalidArgument,
	"Invalid email",
)

// InvalidOtpJwtError invalid otp token
var InvalidOtpJwtError = status.Errorf(
	codes.InvalidArgument,
	"Invalid otp jwt",
)

// InvalidOtpError invalid otp code
var InvalidOtpError = status.Errorf(
	codes.InvalidArgument,
	"Invalid otp",
)

// InvalidAuthTokenError invalid auth token
var InvalidAuthTokenError = status.Errorf(
	codes.Unauthenticated,
	"Auth token is invalid",
)

// InvalidAccessTokenError invalid access token
var InvalidAccessTokenError = status.Errorf(
	codes.Unauthenticated,
	"Access token is invalid",
)

// UnauthenticatedAuthTokenError invalid auth token
var UnauthenticatedAuthTokenError = status.Errorf(
	codes.Unauthenticated,
	"Invalid auth token",
)

// UnauthenticatedAccessTokenError invalid access token
var UnauthenticatedAccessTokenError = status.Errorf(
	codes.Unauthenticated,
	"Invalid auth token",
)

// InvalidForgotPasswordOTPToken invalid forgot password token
var InvalidForgotPasswordOTPToken = status.Errorf(
	codes.InvalidArgument,
	"Invalid otp token",
)

// InvalidPhoneNumberError invalid phone number
var InvalidPhoneNumberError = status.Errorf(
	codes.InvalidArgument,
	"Invalid phone number",
)

// PhoneNumberIsReservedError phone number is reserved
var PhoneNumberIsReservedError = status.Errorf(
	codes.Unavailable,
	"Phone number is reserved",
)

// EmailIsReservedError email is reserved
var EmailIsReservedError = status.Errorf(
	codes.Unavailable,
	"Email is reserved",
)

// InvalidOtpChangeNumberJwtError invalid change number token
var InvalidOtpChangeNumberJwtError = status.Errorf(
	codes.InvalidArgument,
	"Invalid otp change number jwt",
)

// CantCheckUserError can't check user
var CantCheckUserError = status.Errorf(
	codes.Unknown,
	"Can't check user",
)

// RegistrationNotCompletedError registration is not finished
var RegistrationNotCompletedError = status.Errorf(
	codes.Unavailable,
	"Your registration is not completed",
)

// InvalidResetPasswordError invalid reset password token
var InvalidResetPasswordError = status.Errorf(
	codes.InvalidArgument,
	"Invalid reset password token",
)

// NotMatchError data is not match
var NotMatchError = status.Errorf(
	codes.InvalidArgument,
	"Data is not match",
)

// AuthTokenAndDeviceTokenNotMatchError auth and device token is not match
var AuthTokenAndDeviceTokenNotMatchError = status.Errorf(
	codes.InvalidArgument,
	"Auth token and device token is not match",
)

// EmailNotValidError invalid email
var EmailNotValidError = status.Errorf(
	codes.Canceled,
	"Email is not valid",
)

// FIO firstName or lastName not valid
var FioNotValidError = status.Errorf(codes.Canceled, "FIO is not valid")

// PasswordNotValidError password is invalid
func PasswordNotValidError(err error) error {
	return status.Errorf(
		codes.Canceled,
		fmt.Sprintf("Password is not valid: %s", err.Error()),
	)
}

// UsernameNotValidError username iis invalid
func UsernameNotValidError(err error) error {
	return status.Errorf(
		codes.Canceled,
		fmt.Sprintf("Username is not valid: %s", err.Error()),
	)
}

// TokenInvalidError token is invalid
var TokenInvalidError = status.Errorf(
	codes.Unauthenticated,
	"Token is invalid",
)

// TokenTypeInvalidError token type is invalid
var TokenTypeInvalidError = status.Errorf(
	codes.Unauthenticated,
	"Token type is invalid",
)

// PhoneNumberIsSameError phone number is same
var PhoneNumberIsSameError = status.Errorf(
	codes.InvalidArgument,
	"Phone number is same",
)

// InvalidRefreshTokenError invalid refresh token error
var InvalidRefreshTokenError = status.Errorf(
	codes.Unauthenticated,
	"Invalid refresh token",
)

// AlreadyExistError Already exist
var AlreadyExistError = status.Errorf(
	codes.AlreadyExists,
	"Already exist",
)

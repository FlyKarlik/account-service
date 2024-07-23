package models

import (
	"comet/utils"
	"github.com/golang-jwt/jwt"
	"github.com/hashicorp/go-hclog"
	"regexp"
	"time"
)

const (
	// PhoneOTPToken phone otp token
	PhoneOTPToken = "PHONE_OTP"
	// AuthorizeOTPToken Authorize otp token
	AuthorizeOTPToken = "AUTHORIZE_OTP"
	// FirstLoginToken login token
	FirstLoginToken = "FIRST_LOGIN"
	// RegisterToken register token
	RegisterToken = "REGISTER"
	// AuthToken auth token
	AuthToken = "AUTH"
	// RefreshAuthToken refresh auth token
	RefreshAuthToken = "REFRESH_AUTH"
	// AccessToken access token
	AccessToken = "ACCESS"
	// RefreshAccessToken refresh access token
	RefreshAccessToken = "REFRESH_ACCESS"
	// DeviceToken device token
	DeviceToken = "DEVICE"
	// ForgotOTPToken forgot otp token
	ForgotOTPToken = "FORGOT_OTP"
	// ResetPasswordToken reset password token
	ResetPasswordToken = "RESET_PASSWORD"
	// ChangeNumberOTPToken change number otp token
	ChangeNumberOTPToken = "CHANGE_NUMBER_OTP"
)

// Expirations expire time of tokens
var Expirations = map[string]time.Duration{
	PhoneOTPToken:        time.Minute * 5,
	AuthorizeOTPToken:    time.Minute * 5,
	FirstLoginToken:      time.Minute * 30,
	RegisterToken:        time.Minute * 15,
	AuthToken:            time.Hour * 24,
	RefreshAuthToken:     time.Hour * 24 * 7,
	AccessToken:          time.Minute * 15,
	RefreshAccessToken:   time.Hour * 2,
	DeviceToken:          time.Hour * 24 * 1000,
	ForgotOTPToken:       time.Hour * 24,
	ResetPasswordToken:   time.Hour * 1,
	ChangeNumberOTPToken: time.Hour * 1,
}

// RefreshRegex refresh token regex
var RefreshRegex = regexp.MustCompile(`^REFRESH_(.*)$`)

// JWT jwt object
type JWT struct {
	ID          string
	Identity    string
	Variety     string
	Email       string
	JwtSecret   string
	Exp         int64
	IsRevoked   bool
	LastUse     *time.Time
	TokenObject *Token
	Extra       jwt.MapClaims
}

// ToJWTString to jwt token string
func (j *JWT) ToJWTString() string {
	log := hclog.Default()

	atClaims := jwt.MapClaims{
		"variety":  j.Variety,
		"identity": j.Identity,
		"exp":      j.Exp,
		"id":       j.ID,
	}

	if j.Email != "" {
		atClaims["username"] = j.Email
	}

	if j.Extra != nil {
		for k := range j.Extra {
			atClaims[k] = j.Extra[k]
		}
	}

	token, err := utils.EncryptJWT(&atClaims, j.JwtSecret)
	if err != nil {
		log.Error("[tokens.ToJWTString] utils.EncryptJWT", "error", err)
		return ""
	}

	return token
}

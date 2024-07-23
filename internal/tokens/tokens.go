package tokens

import (
	"account-service/config"
	"account-service/internal/models"
	"comet/utils"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"go.opentelemetry.io/otel/trace"
	"time"
)

// TokenService token service
type TokenService struct {
	db     Repository
	tracer trace.Tracer
	cfg    *config.Config
}

// NewToken Create a new Token
func NewToken(db Repository, tr trace.Tracer, cfg *config.Config) *TokenService {
	return &TokenService{
		db:     db,
		tracer: tr,
		cfg:    cfg,
	}
}

func (t *TokenService) NewJWT(ctx context.Context, variety string, identity string, email string, extra jwt.MapClaims) (*models.JWT, error) {
	tr := t.tracer
	ctx, span := tr.Start(ctx, "new-jwt")
	defer span.End()

	tokenObj, err := t.db.CreateToken(ctx, variety, identity)
	if err != nil {
		return nil, fmt.Errorf("t.db.CreateToken error: %w", err)
	}

	expiration, ok := models.Expirations[variety]
	if !ok {
		return nil, fmt.Errorf("invalid variety: %s", variety)
	}

	return &models.JWT{
		ID:          tokenObj.ID,
		Variety:     variety,
		Identity:    identity,
		Email:       email,
		JwtSecret:   t.cfg.JwtSecret,
		Exp:         time.Now().Add(expiration).Unix(),
		IsRevoked:   false,
		LastUse:     nil,
		TokenObject: tokenObj,
		Extra:       extra,
	}, nil
}

// ParseJWT parse jwt token
func (t *TokenService) ParseJWT(ctx context.Context, token string) (*models.JWT, error) {
	tr := t.tracer
	ctx, span := tr.Start(ctx, "parse-jwt")
	defer span.End()
	verifyJWT, err := utils.VerifyJWT(token, t.cfg.JwtSecret)
	if err != nil {
		return nil, fmt.Errorf("utils.VerifyJWT error: %w", err)
	}

	claim := verifyJWT.Claims.(jwt.MapClaims)

	variety, ok := claim["variety"]
	if !ok {
		return nil, fmt.Errorf("token have not variety: %s", claim)
	}

	id, ok := claim["id"]
	if !ok {
		return nil, fmt.Errorf("token have not id: %s", claim)
	}

	tokenObj, err := t.db.GetTokenByID(ctx, id.(string), variety.(string))
	if err != nil {
		return nil, fmt.Errorf("utils.GetTokenByID error: %w", err)
	}
	if tokenObj == nil {
		return nil, fmt.Errorf("tokenObj is nil")
	}

	identity, ok := claim["identity"]
	if !ok {
		return nil, fmt.Errorf("token have not identity: %s", claim)
	}

	exp, ok := claim["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("token have not exp: %s", claim)
	}

	rawEmail, ok := claim["email"]
	var email string
	if ok {
		email = rawEmail.(string)
	}

	return &models.JWT{
		ID:          tokenObj.ID,
		Variety:     variety.(string),
		Identity:    identity.(string),
		Email:       email,
		JwtSecret:   t.cfg.JwtSecret,
		Exp:         int64(exp),
		IsRevoked:   tokenObj.IsRevoked,
		LastUse:     &tokenObj.LastUse,
		TokenObject: tokenObj,
	}, nil
}

// ValidateRefreshJWT validate refresh token
func (t *TokenService) ValidateRefreshJWT(ctx context.Context, token string) (*models.JWT, string, error) {
	tr := t.tracer
	ctx, span := tr.Start(ctx, "validate-refresh-jwt")
	defer span.End()

	tok, err := t.ParseJWT(ctx, token)
	if err != nil {
		return nil, "", fmt.Errorf("ParseJWT error: %w", err)
	}

	if tok.Variety != models.RefreshAuthToken && tok.Variety != models.RefreshAccessToken {
		return nil, "", fmt.Errorf("token is not valid")
	}

	result := models.RefreshRegex.FindStringSubmatch(tok.Variety)
	if len(result) < 2 {
		return nil, "", fmt.Errorf("length result findStringSubmatch no valid")
	}

	tokenType := result[1]

	if err := t.Validate(tok, "REFRESH_"+tokenType); err != nil {
		return nil, "", fmt.Errorf("token is not valid: %w", err)
	}

	return tok, tokenType, nil
}

// Validate validate JWT by variety
func (t *TokenService) Validate(j *models.JWT, variety string) error {
	if j.IsRevoked {
		return fmt.Errorf("jwt not revoked")
	}

	if j.Variety != variety {
		return fmt.Errorf("variety failed")
	}

	t.db.SetUse(j.TokenObject)

	return nil
}

// CreateAccessJWT create access and access refresh
func (t *TokenService) CreateAccessJWT(ctx context.Context, identity string, email string) (*models.JWT, *models.JWT, error) {
	authToken, err := t.NewJWT(ctx, models.AccessToken, identity, email, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("[tokens.CreateAccessJWT] t.NewJWT access: %w", err)
	}

	refreshAuthToken, err := t.NewJWT(ctx, models.RefreshAccessToken, identity, email, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("[tokens.CreateAccessJWT] t.NewJWT refresh access: %w", err)
	}

	return authToken, refreshAuthToken, nil
}

// CreateAuthJWT create auth and auth refresh
func (t *TokenService) CreateAuthJWT(ctx context.Context, identity string, email string) (*models.JWT, *models.JWT, error) {
	authToken, err := t.NewJWT(ctx, models.AuthToken, identity, email, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("[tokens.CreateAuthJWT] t.NewJWT auth: %w", err)
	}

	refreshAuthToken, err := t.NewJWT(ctx, models.RefreshAuthToken, identity, email, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("[tokens.CreateAuthJWT] t.NewJWT refresh auth: %w", err)
	}

	return authToken, refreshAuthToken, nil
}

// Revoke revoke JWT
func (t *TokenService) Revoke(j *models.JWT) {
	t.db.Revoke(j.TokenObject)
}

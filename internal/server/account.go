package server

import (
	"account-service/config"
	"account-service/internal/server/interfaces"
	"go.opentelemetry.io/otel/trace"
	protos "protos/account"
)

// AccountService grpc service declaration
type AccountService struct {
	protos.UnimplementedAccountServiceServer
	db       interfaces.Repository
	tokenSrv interfaces.TokenService
	trace    trace.Tracer
	cfg      *config.Config
}

// NewAccount Creates a new Account server
func NewAccount(db interfaces.Repository, t interfaces.TokenService, tracer trace.Tracer, cfg *config.Config) *AccountService {
	return &AccountService{
		db:       db,
		tokenSrv: t,
		trace:    tracer,
		cfg:      cfg,
	}
}

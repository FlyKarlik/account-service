package main

import (
	"account-service/config"
	"account-service/internal/models"
	"account-service/internal/server"
	"account-service/internal/server/repository"
	"account-service/internal/tokens"
	tokensRepository "account-service/internal/tokens/repository"
	"comet/db"
	"comet/utils"
	"flag"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/go-hclog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	protos "protos/account"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()
	log := hclog.Default()

	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}

	err = utils.InitSentry(cfg.SentryDSN)
	if err != nil {
		return fmt.Errorf("failed init sentry: %w", err)
	}
	defer sentry.Flush(2 * time.Second)

	tracer, err := utils.InitTracerProvider(cfg.JaegerHost, "account-service")
	if err != nil {
		return fmt.Errorf("failed init Tracer: %w", err)
	}

	database, err := db.NewDatabase(cfg.DbDsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = database.AutoMigrate(&models.User{}, &models.Token{}, &models.Department{}, &models.Role{})
	if err != nil {
		return fmt.Errorf("failed AutoMigrate database: %w", err)
	}

	creds, err := credentials.NewServerTLSFromFile("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return fmt.Errorf("failed to setup TLS: %w", err)
	}

	// Create a new gRPC srv
	gs := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	repoAccount := repository.NewRepository(database)
	repoToken := tokensRepository.NewRepository(database, tracer)
	tokenSrv := tokens.NewToken(repoToken, tracer, cfg)

	srv := server.NewAccount(repoAccount, tokenSrv, tracer, cfg)

	protos.RegisterAccountServiceServer(gs, srv)

	reflection.Register(gs)

	l, err := net.Listen("tcp", cfg.ServerHost)
	if err != nil {
		log.Error("Unable to create listener", "error", err)
		os.Exit(1)
	}

	log.Info("service is running.")

	// listen for requests
	err = gs.Serve(l)
	if err != nil {
		return fmt.Errorf("failed serving: %w", err)
	}

	return nil
}

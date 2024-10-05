package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/rmishgoog/adopt-a-dog/apis/services/api/debug"
	"github.com/rmishgoog/adopt-a-dog/core/api/auth"
	"github.com/rmishgoog/adopt-a-dog/foundations/keystore"
	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

var build = "develop"

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "****sending alert****")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "Adoptions", traceIDFn, events)

	ctx := context.Background()
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "start-up errors", "msg", err)
		os.Exit(1)
	}

}

func run(ctx context.Context, log *logger.Logger) error {

	log.Info(ctx, "server-bootstrap", "GOMAXPROCS", runtime.GOMAXPROCS(0), "build", build)

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:3010"`
			CORSAllowedOrigins []string      `conf:"default:*,mask"`
		}
		Auth struct {
			Host string `conf:"default:https://local.auth.adoptadog.com/realms/adoptadog/.well-known/openid-configuration"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "adoptions-auth",
		},

		Web: struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:3010"`
			CORSAllowedOrigins []string      `conf:"default:*,mask"`
		}{
			ShutdownTimeout: 60 * time.Second,
		},
	}

	const prefix = "adoptions-auth"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	log.Info(ctx, "starting auth service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("failed marshalling config: %w", err)
	}
	log.Info(ctx, "auth service startup", "config", out)
	expvar.NewString("build").Set(cfg.Build)

	ks := keystore.New()

	if err := ks.PublicKey(cfg.Auth.Host, true); err != nil {
		return fmt.Errorf("failed to fetch public key, likely the OIDC service is not up or having issues: %w", err)
	}
	// TODO: Implement this with a Kubernetes configmap or secret mounted as a volume on the container.
	authCfg := auth.Config{
		Log:         log,
		JWTValidate: ks,
		Issuer:      "https://local.auth.adoptadog.com/realms/adoptadog",
	}

	auth, err := auth.New(authCfg)
	if err != nil {
		return fmt.Errorf("failed to create auth data struct: %w", err)
	}
	fmt.Printf("Issuer: %s\n", auth.Issuer())

	// Start the debug service, this is omnipresent in all services. A debug host running locally to help with debugging.
	go func() {
		log.Info(ctx, "auth-debug-service", "status", "started", "host", cfg.Web.DebugHost)
		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "auth-debug-service", "status", "shutdown", "host", cfg.Web.DebugHost, "err", err)
		}
	}()

	// Start the authn/authz API service.
	log.Info(ctx, "startup", "status", "initializing authentication & authorization support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	apirouter := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      nil, //mux.WebAPI(log, shutdown),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "starting authentication & authorization router", "host", apirouter.Addr)
		serverErrors <- apirouter.ListenAndServe()
	}()

	// This is a blocking select statement, which will wait for either a server error on start up or a shutdown signal.
	// This is where the application shall hanlde the graceful shutdown of the server or an interrupt signal.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("auth server error: %w", err)
	case sig := <-shutdown:
		log.Info(ctx, "auth server-shutdown", "status", "auth server shutdown started", "signal", sig)
		defer log.Info(ctx, "auth server-shutdown", "status", "auth server shutdown complete", "signal", sig)
		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()
		if err := apirouter.Shutdown(ctx); err != nil {
			apirouter.Close()
			return fmt.Errorf("could not stop api server gracefully: %w", err)
		}
	}

	return nil

}

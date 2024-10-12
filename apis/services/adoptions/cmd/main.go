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
	"github.com/rmishgoog/adopt-a-dog/apis/services/adoptions/builder"
	"github.com/rmishgoog/adopt-a-dog/apis/services/adoptions/mux"
	"github.com/rmishgoog/adopt-a-dog/apis/services/api/debug"
	"github.com/rmishgoog/adopt-a-dog/app/api/authclient"
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

	// Load the configuration required for the service to function correctly.
	log.Info(ctx, "server-bootstrap", "GOMAXPROCS", runtime.GOMAXPROCS(0), "build", build)

	authHost := os.Getenv("AUTH_HOST")
	if authHost == "" {
		return fmt.Errorf("fatal error, no authentication endpoint set: %w", errors.New("environment variable AUTH_HOST not set"))
	}

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
			Host string `conf:"default:http://auth-service.sales-system.svc.cluster.local:6000"`
		}
		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:postgres,mask"`
			HostPort     string `conf:"default:database-service.sales-system.svc.cluster.local"`
			Name         string `conf:"default:postgres"`
			MaxIdleConns int    `conf:"default:2"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "adoptions",
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
		Auth: struct {
			Host string "conf:\"default:http://auth-service.sales-system.svc.cluster.local:6000\""
		}{
			Host: authHost,
		},
	}

	const prefix = "adoptions"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("failed marshalling config: %w", err)
	}
	log.Info(ctx, "service startup", "config", out)
	expvar.NewString("build").Set(cfg.Build)

	// Start the debug service.
	go func() {
		log.Info(ctx, "debug-service", "status", "started", "host", cfg.Web.DebugHost)
		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "debug-service", "status", "shutdown", "host", cfg.Web.DebugHost, "err", err)
		}
	}()

	// Start the core API service.
	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	logFunc := func(ctx context.Context, msg string, v ...any) {
		log.Info(ctx, msg, v...)
	}
	authclient := authclient.New(cfg.Auth.Host, logFunc)

	// Create a comprehensive mux configuration & pass it along, no discrete passing of values for the mux.
	// This shall contain everything from logger, shutdown channels to build information that web api needs!
	muxConfig := mux.Config{
		Build:      build,
		Log:        log,
		AuthClient: authclient,
		Shutdown:   shutdown,
	}

	apirouter := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      mux.WebAPI(muxConfig, builder.Routes()),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)
	// This shall server as the parent goroutine for the API requests, each request will be handled by a child goroutine created by the routine running the router.
	go func() {
		log.Info(ctx, "startup", "status", "starting api router", "host", apirouter.Addr)
		serverErrors <- apirouter.ListenAndServe()
	}()

	// Router shutdown & handling process interruption signals.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info(ctx, "server-shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "server-shutdown", "status", "shutdown complete", "signal", sig)
		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()
		if err := apirouter.Shutdown(ctx); err != nil {
			apirouter.Close()
			return fmt.Errorf("could not stop api server gracefully: %w", err)
		}
	}

	return nil

}

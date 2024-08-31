package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
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
		//return web.GetTraceID(ctx)
		return ""
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

	// Load the configuration.
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

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	sig := <-shutdown

	log.Info(ctx, "server-shutdown", "status", "shutdown started", "signal", sig)
	defer log.Info(ctx, "server-shutdown", "status", "shutdown complete", "signal", sig)

	return nil

}

package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
)

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

	log.Info(ctx, "server-bootstrap", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	sig := <-shutdown

	log.Info(ctx, "server-shutdown", "status", "shutdown started", "signal", sig)
	defer log.Info(ctx, "server-shutdown", "status", "shutdown complete", "signal", sig)

	return nil

}

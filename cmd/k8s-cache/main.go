package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/grafana/beyla-k8s-cache/pkg/service"
)

const defaultPort = 50055

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})))

	ic := service.InformersCache{
		Port: defaultPort,
		// TODO: make configurable
		ResyncPeriod: 30 * time.Minute,
	}
	portStr := os.Getenv("BEYLA_K8S_CACHE_PORT")
	if portStr != "" {
		var err error
		if ic.Port, err = strconv.Atoi(portStr); err != nil {
			slog.Error("invalid BEYLA_K8S_CACHE_PORT, using default port", "error", err)
			ic.Port = defaultPort
		}
	}

	if err := ic.Run(context.Background()); err != nil {
		slog.Error("starting informers' cache service", "error", err)
		os.Exit(-1)
	}
	slog.Info("service stopped. Exiting now")
}

package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/am3o/TP-Link-Router_exporter/pkg/metrics"
	"github.com/am3o/TP-Link-Router_exporter/pkg/router"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Version = "0.0.0"

func main() {
	routerURL, exists := os.LookupEnv("TP-LINK-ROUTER-URL")
	if !exists {
		routerURL = "192.168.0.1"
	}

	username, _ := os.LookupEnv("TP-LINK-ROUTER-USERNAME")
	password, _ := os.LookupEnv("TP-LINK-ROUTER-PASSWORD")

	collector := metrics.New()
	prometheus.MustRegister(collector)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx := context.Background()
	defer ctx.Done()

	go func(ctx context.Context, url, user, password string) {
		router, err := router.New(url, user, password)
		if err != nil {
			logger.With(slog.Any("error", err)).Error("could not initialize router")
			panic(err)
		}

		ticker := time.NewTicker(15 * time.Second)
		for ; ; <-ticker.C {
			sessionCtx, err := router.Login(ctx)
			if err != nil {
				logger.With(slog.Any("error", err)).Error("could not login to router")
				continue
			}

			RxLAN, TxLAN, err := router.LANTraffic(sessionCtx)
			if err != nil {
				logger.With(slog.Any("error", err)).Error("could not get to LAN information")
				continue
			}

			for client, value := range RxLAN {
				collector.RxLAN(client, value)
			}

			for client, value := range TxLAN {
				collector.TxLAN(client, value)
			}

			RxWAN, TxWAN, err := router.WANTraffic(sessionCtx)
			if err != nil {
				logger.With(slog.Any("error", err)).Error("could not get to WAN information")
				continue
			}

			for client, value := range RxWAN {
				collector.RxWAN(client, value)
			}

			for client, value := range TxWAN {
				collector.TxWAN(client, value)
			}

			if err := router.Logout(sessionCtx); err != nil {
				logger.With(slog.Any("error", err)).Error("could not logout to router")
				continue
			}
		}
	}(ctx, routerURL, username, password)

	http.Handle("/metrics", promhttp.Handler())
	logger.With(
		slog.String("version", Version),
	).InfoContext(ctx, "start TP-Link router exporter")
	server := &http.Server{
		Addr:              net.JoinHostPort("", "8080"),
		ReadHeaderTimeout: 3 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		logger.With(slog.Any("error", err)).ErrorContext(ctx, "could not start service")
		os.Exit(1)
	}
}

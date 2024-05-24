package main

import (
	"context"
	"fmt"
	stdlog "log"
	"log/slog"
	stdhttp "net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/jessevdk/go-flags"
	"golang.org/x/sync/errgroup"

	"github.com/belliel/crypto-price-aggregator/internal/server/http"
	cryptoHandler "github.com/belliel/crypto-price-aggregator/internal/server/http/crypto"
	cryptoService "github.com/belliel/crypto-price-aggregator/internal/service/crypto"
	"github.com/belliel/crypto-price-aggregator/pkg/sre/log"
)

var options = struct {
	HTTPPort int `long:"http-port" description:"http binding port" env:"HTTP_PORT" default:"8080"`

	CryptoServiceTargetCryptoLastTradedPricesPairs          string        `long:"crypto-service-target-last-traded-prices-pairs" description:"crypto service worker default pairs" env:"CRYPTO_SERVICE_TARGET_CRYPTO_LAST_TRADED_PRICES_PAIRS" default:"BTC/USD,BTC/CHF,BTC/EUR"`
	CryptoServiceTargetCryptoLastTradedPricesWorkerInterval time.Duration `long:"crypto-service-target-last-traded-prices-worker-interval" description:"crypto service worker fetch interval" env:"CRYPTO_SERVICE_TARGET_CRYPTO_LAST_TRADED_PRICES_WORKER_INTERVAL" default:"10s"`
	CryptoServiceCacheTTL                                   time.Duration `long:"crypto-service-cache-ttl" description:"crypto service worker cache ttl" env:"CRYPTO_SERVICE_CACHE_TTL" default:"30s"`

	LogLevel string `long:"log-level" env:"LOG_LEVEL" choice:"debug" choice:"error" choice:"warn" choice:"info" default:"info"`
}{}

func main() {
	if _, err := flags.Parse(&options); err != nil {
		stdlog.Fatal(fmt.Errorf("failed to parse flags: %w", err))
	}

	log.SetupLogger(log.ToSLogLevel(options.LogLevel))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceCrypto := cryptoService.NewService(
		strings.Split(options.CryptoServiceTargetCryptoLastTradedPricesPairs, ","),
		options.CryptoServiceCacheTTL,
		options.CryptoServiceTargetCryptoLastTradedPricesWorkerInterval,
		stdhttp.DefaultClient,
	)

	r := chi.NewRouter()
	{
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(httplog.RequestLogger(httplog.NewLogger("ltp", httplog.Options{
			JSON:             true,
			LogLevel:         log.ToSLogLevel(options.LogLevel),
			Concise:          true,
			RequestHeaders:   true,
			MessageFieldName: "message",
			Tags:             map[string]string{},
			QuietDownRoutes: []string{
				"/",
				"/ping",
			},
			SourceFieldName: "source",
		})))
		r.Use(middleware.Recoverer)

		r.Route("/api/v1", func(router chi.Router) {
			cryptoHandler.NewHandler(serviceCrypto).SetupRouter(router)
		})
	}

	server := http.NewServer(ctx, fmt.Sprintf(":%d", options.HTTPPort), r)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return server.Serve()
	})

	g.Go(func() error {
		return serviceCrypto.LastTradedPricesForDefaultPairsWorker(ctx)
	})

	g.Go(func() error {
		<-gCtx.Done()
		return server.Shutdown()
	})

	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	slog.Info("run and waiting for connections")
	slog.Info(fmt.Sprintf("http server run on: http://localhost:%d", options.HTTPPort))

	if err := g.Wait(); err != nil {
		slog.With("reason", err).Error("exited")
	}
}

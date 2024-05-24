package crypto

import (
	"github.com/belliel/crypto-price-aggregator/internal/service"
	"log/slog"
	"net/http"
	"sort"
	"sync"
	"time"
)

var logger = slog.With("service", "crypto")

var _ service.Crypto = (*Service)(nil)

type Service struct {
	cacheMutex    sync.RWMutex
	cache         map[string]cacheEntry
	cacheEntryTTL time.Duration

	tickerInterval                     time.Duration
	lastTradedCryptoPricesDefaultPairs []string

	httpClient *http.Client
}

func NewService(defaultPairs []string, cacheTTL time.Duration, tickerInterval time.Duration, client *http.Client) *Service {
	sort.Strings(defaultPairs)

	return &Service{
		// min 1
		cache:                              make(map[string]cacheEntry, 1),
		cacheMutex:                         sync.RWMutex{},
		cacheEntryTTL:                      cacheTTL,
		tickerInterval:                     tickerInterval,
		lastTradedCryptoPricesDefaultPairs: defaultPairs,
		httpClient:                         client,
	}
}

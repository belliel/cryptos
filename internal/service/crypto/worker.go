package crypto

import (
	"context"
	"github.com/belliel/crypto-price-aggregator/internal/entities/crypto"
	"strings"
	"time"
)

// lastTradedCryptoPriceAPI - uses kraken API
// https://docs.kraken.com/rest/#tag/Spot-Market-Data/operation/getTickerInformation
const lastTradedCryptoPriceAPI = "https://api.kraken.com/0/public/Ticker?pair=%s"

func (s *Service) LastTradedPricesForDefaultPairsWorker(ctx context.Context) error {
	log := logger.With("method", "lastTradedPricesForDefaultPairsWorker")

	ticker := time.NewTicker(s.tickerInterval)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return context.Canceled
		case <-ticker.C:
			result, err := s.getLastPricesFromSource(ctx, &crypto.GetLastTradedPriceQuery{Pairs: s.lastTradedCryptoPricesDefaultPairs})
			if err != nil {
				log.With("error", err.Error()).Error("something went wrong on worker")
				continue
			}

			cacheKey := strings.Join(s.lastTradedCryptoPricesDefaultPairs, ",")

			s.cacheMutex.Lock()
			s.cache[cacheKey] = cacheEntry{
				Time: time.Now().Unix(),
				Data: result,
			}
			s.cacheMutex.Unlock()
		}
	}
}

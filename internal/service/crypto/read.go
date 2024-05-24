package crypto

import (
	"context"
	"fmt"
	"github.com/belliel/crypto-price-aggregator/internal/entities/crypto"
	"sort"
	"strings"
	"time"
)

func (s *Service) GetLastTradedPrice(ctx context.Context, query *crypto.GetLastTradedPriceQuery) ([]crypto.GetLastTradedPriceResult, error) {
	sort.Strings(query.Pairs)
	if len(query.Pairs) == 0 || query.Pairs[0] == "" {
		query.Pairs = s.lastTradedCryptoPricesDefaultPairs
	}

	now := time.Now().Unix()
	cacheKey := strings.Join(query.Pairs, ",")

	var result []crypto.GetLastTradedPriceResult

	s.cacheMutex.RLock()
	res, ok := s.cache[cacheKey]
	s.cacheMutex.RUnlock()
	if ok && now-res.Time < int64(s.cacheEntryTTL.Seconds()) {
		return res.Data, nil
	}

	result, err := s.getLastPricesFromSource(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get last traded crypto prices for pairs: %w", err)
	}

	s.cacheMutex.Lock()
	s.cache[cacheKey] = cacheEntry{
		Time: time.Now().Add(s.cacheEntryTTL).Unix(),
		Data: result,
	}
	s.cacheMutex.Unlock()

	return result, nil
}

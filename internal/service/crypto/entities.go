package crypto

import "github.com/belliel/crypto-price-aggregator/internal/entities/crypto"

type krakenTickerLTPResult struct {
	CloseWith []string `json:"c"`
}

type krakenTickerResult struct {
	Error  []string                         `json:"error"`
	Result map[string]krakenTickerLTPResult `json:"result"`
}

type cacheEntry struct {
	Time int64
	Data []crypto.GetLastTradedPriceResult
}

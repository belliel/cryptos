package crypto

import "github.com/belliel/crypto-price-aggregator/internal/entities/crypto"

func toLastTradedPriceResponse(results []crypto.GetLastTradedPriceResult) LastTradedPriceResponse {
	ltpEntries := make([]LTP, len(results))
	for i := 0; i < len(results); i++ {
		ltpEntries[i] = LTP{
			Pair:   results[i].Pair,
			Amount: results[i].Amount,
		}
	}

	return LastTradedPriceResponse{LTP: ltpEntries}
}

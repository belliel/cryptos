package crypto

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/belliel/crypto-price-aggregator/internal/entities/crypto"
)

func (s *Service) getLastPricesFromSource(ctx context.Context, query *crypto.GetLastTradedPriceQuery) ([]crypto.GetLastTradedPriceResult, error) {
	log := logger.With("method", "getLastPricesFromSource")

	apiURL := fmt.Sprintf(lastTradedCryptoPriceAPI, strings.Join(query.Pairs, ","))

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request with context: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform last traded crypto price request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.With("status_code", resp.StatusCode).Error("request failed with bad status code")
		return nil, fmt.Errorf("api failed with bad status code: %s: %d", apiURL, resp.StatusCode)
	}

	var apiResult krakenTickerResult
	if err := json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}

	if len(apiResult.Error) != 0 {
		log.With("api_errors", apiResult.Error).Error("api response contains errors")
		return nil, ErrAPIHaveFailed
	}
	if len(apiResult.Result) != len(query.Pairs) {
		resultPairs := make([]string, 0, len(apiResult.Result))
		for key := range apiResult.Result {
			resultPairs = append(resultPairs, key)
		}
		sort.Strings(resultPairs)

		log.With("pairs", query.Pairs, "result_pairs", resultPairs).Warn("query pairs and result pairs have different length")
		return nil, ErrPairsDifferentLength
	}

	var result = make([]crypto.GetLastTradedPriceResult, 0, len(apiResult.Result))
	for i := 0; i < len(query.Pairs); i++ {
		// we don't check for pair existence
		// because it's assumed as resolved before
		res := apiResult.Result[query.Pairs[i]]

		if len(res.CloseWith) < 1 {
			log.With("pair", query.Pairs[i], "close_prices", res.CloseWith).Error(ErrAPIBadDataNoClosePrice.Error())
			return nil, ErrAPIBadDataNoClosePrice
		}

		result = append(result, crypto.GetLastTradedPriceResult{
			Pair:   query.Pairs[i],
			Amount: res.CloseWith[0],
		})
	}

	return result, nil
}

package crypto

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"

	"github.com/belliel/crypto-price-aggregator/internal/entities/crypto"
)

func TestGetLastPricesFromSource(t *testing.T) {
	tests := []struct {
		Name        string
		Pairs       []string
		Expect      []crypto.GetLastTradedPriceResult
		ExpectError bool
		Setup       func(pairs []string) *http.Client
	}{
		{
			Name:  "ok",
			Pairs: []string{"BTC/USD"},
			Expect: []crypto.GetLastTradedPriceResult{
				{
					Pair:   "BTC/USD",
					Amount: "228",
				},
			},
			Setup: func(pairs []string) *http.Client {
				h := &http.Client{}

				readyURL := fmt.Sprintf(lastTradedCryptoPriceAPI, strings.Join(pairs, ","))
				u, _ := url.Parse(readyURL)

				gock.InterceptClient(h)
				gock.New(u.Scheme+"://"+u.Host).Get(u.Path).MatchParam("pair", u.Query().Get("pair")).Reply(http.StatusOK).BodyString(`
{
	"error": [],
	"result": {
		"BTC/USD": {
			"o": ["228", "322"]
		}
	}
}
`)
				return h
			},
		},
		{
			Name:  "ok 2",
			Pairs: []string{"BTC/USD", "BTC/ETC"},
			Expect: []crypto.GetLastTradedPriceResult{
				{
					Pair:   "BTC/USD",
					Amount: "228",
				},
				{
					Pair:   "BTC/ETC",
					Amount: "1.1",
				},
			},
			Setup: func(pairs []string) *http.Client {
				h := &http.Client{}

				readyURL := fmt.Sprintf(lastTradedCryptoPriceAPI, strings.Join(pairs, ","))
				u, _ := url.Parse(readyURL)

				gock.InterceptClient(h)
				gock.New(u.Scheme+"://"+u.Host).Get(u.Path).MatchParam("pair", u.Query().Get("pair")).Reply(http.StatusOK).BodyString(`
{
	"error": [],
	"result": {
		"BTC/USD": {
			"o": ["228", "322"]
		},
		"BTC/ETC": {
			"o": ["1.1", "2"]
		}
	}
}
`)
				return h
			},
		},
		{
			Name:        "bad no",
			Pairs:       []string{"BTC/USD", "BTC/ETC"},
			Expect:      nil,
			ExpectError: true,
			Setup: func(pairs []string) *http.Client {
				h := &http.Client{}

				readyURL := fmt.Sprintf(lastTradedCryptoPriceAPI, strings.Join(pairs, ","))

				u, _ := url.Parse(readyURL)

				gock.InterceptClient(h)
				gock.New(u.Scheme+"://"+u.Host).Get(u.Path).MatchParam("pair", u.Query().Get("pair")).Reply(http.StatusOK).BodyString(`
{
	"error": ["oh no there is an error"],
	"result": {}
}
`)
				return h
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			httpClient := tc.Setup(tc.Pairs)

			result, err := NewService([]string{}, time.Duration(1), time.Duration(1), httpClient).
				getLastPricesFromSource(context.Background(), &crypto.GetLastTradedPriceQuery{Pairs: tc.Pairs})
			if !tc.ExpectError {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.Expect, result)
		})
	}
}

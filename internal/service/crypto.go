package service

import (
	"context"
	"github.com/belliel/crypto-price-aggregator/internal/entities/crypto"
)

type Crypto interface {
	GetLastTradedPrice(ctx context.Context, query *crypto.GetLastTradedPriceQuery) ([]crypto.GetLastTradedPriceResult, error)
}

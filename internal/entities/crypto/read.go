package crypto

type GetLastTradedPriceQuery struct {
	Pairs []string
}

type GetLastTradedPriceResult struct {
	Pair   string
	Amount string
}

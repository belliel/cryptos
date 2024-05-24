package crypto

type LTP struct {
	Pair   string `json:"pair"`
	Amount string `json:"amount"`
}

type LastTradedPriceResponse struct {
	LTP []LTP `json:"ltp"`
}

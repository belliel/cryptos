package crypto

import "errors"

var (
	ErrPairsDifferentLength   = errors.New("query pairs and result pairs have different length")
	ErrAPIHaveFailed          = errors.New("api response contains errors")
	ErrAPIBadDataNoClosePrice = errors.New("api response not contains close price")
)

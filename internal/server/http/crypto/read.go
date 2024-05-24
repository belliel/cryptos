package crypto

import (
	"errors"
	"github.com/belliel/crypto-price-aggregator/internal/entities/crypto"
	"github.com/belliel/crypto-price-aggregator/internal/server/http/api"
	cryptoService "github.com/belliel/crypto-price-aggregator/internal/service/crypto"
	"net/http"
	"strings"
)

func (h *Handler) LastTradedPrice(w http.ResponseWriter, r *http.Request) {
	pairs := strings.Split(r.URL.Query().Get("pairs"), ",")

	result, err := h.ServiceCrypto.GetLastTradedPrice(r.Context(), &crypto.GetLastTradedPriceQuery{Pairs: pairs})
	if err != nil {
		if errors.Is(err, cryptoService.ErrAPIHaveFailed) {
			api.WriteError(w, http.StatusBadRequest, &api.Error{
				Message: "bad request",
				Code:    http.StatusBadRequest,
			})
			return
		}
		api.WriteError(w, http.StatusInternalServerError, &api.Error{
			Message: "internal error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	api.WriteOk(w, toLastTradedPriceResponse(result))
}

package crypto

import (
	"github.com/belliel/crypto-price-aggregator/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	ServiceCrypto service.Crypto
}

func NewHandler(crypto service.Crypto) *Handler {
	return &Handler{ServiceCrypto: crypto}
}

func (h *Handler) SetupRouter(router chi.Router) {
	router.Get("/ltp", h.LastTradedPrice)
}

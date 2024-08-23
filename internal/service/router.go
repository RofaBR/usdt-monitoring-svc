package service

import (
	"github.com/RofaBR/usdt-monitoring-svc/internal/service/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
		),
	)

	handler := handlers.NewHandler(s.storage)

	r.Route("/integrations/usdt-monitoring-svc", func(r chi.Router) {
		r.Get("/transfers", handler.GetTransfers)
	})

	return r
}

package service

import (
	"net"
	"net/http"

	"github.com/RofaBR/usdt-monitoring-svc/internal/config"
	"github.com/RofaBR/usdt-monitoring-svc/internal/storage"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log      *logan.Entry
	copus    types.Copus
	listener net.Listener
	storage  storage.Storage
}

func (s *service) run() error {
	s.log.Info("Service started")
	r := s.router()

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func newService(cfg config.Config) *service {

	db := storage.NewPostgresStorage(cfg.DB().RawDB())

	return &service{
		log:      cfg.Log(),
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
		storage:  db,
	}
}

func Run(cfg config.Config) {
	svc := newService(cfg)

	go svc.GetTransferEvents(cfg)

	if err := svc.run(); err != nil {
		panic(err)
	}
}

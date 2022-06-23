package api

import (
	"context"
	"embed"
	"net"

	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/ignite/backend/api/pb"
)

//go:embed schemas/*
var fsSchemas embed.FS

const (
	DefaultAddress = "127.0.0.1:4040"
)

type Option = func(*Service)

func WithAddress(addr string) Option {
	return func(s *Service) {
		s.addr = addr
	}
}

func NewService(db postgres.Adapter, options ...Option) Service {
	s := Service{
		addr:    DefaultAddress,
		db:      db,
		schemas: postgres.NewSchemas(fsSchemas, "api"),
	}

	for _, o := range options {
		o(&s)
	}

	return s
}

type Service struct {
	addr    string
	db      postgres.Adapter
	schemas postgres.Schemas
}

func (s Service) Run(ctx context.Context) error {
	log.WithFields(log.Fields{
		"backend": s.db.GetType(),
	}).Debug("Updating data backend schema")

	// Make sure the API database schemas are up to date
	if err := s.db.UpdateSchema(ctx, s.schemas); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	srv := grpc.NewServer()

	pb.RegisterMapperServer(srv, &server{db: s.db})

	// Start listening for requests
	errc := make(chan error)
	go func() {
		log.WithFields(log.Fields{
			"address": listener.Addr(),
		}).Info("Listening for requests")

		if err := srv.Serve(listener); err != nil {
			errc <- err
		}
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}

	return nil
}

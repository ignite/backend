package collector

import (
	"context"
	"errors"
	"time"

	"github.com/ignite/cli/ignite/pkg/cosmostxcollector"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultMinHeight = 1
	DefaultTimeout   = 5 * time.Minute
	DefaultGrace     = 5 * time.Millisecond
)

type Option func(*Service)

func WithGrace(t time.Duration) Option {
	return func(s *Service) {
		s.grace = t
	}
}

func WithTimeout(t time.Duration) Option {
	return func(s *Service) {
		s.timeout = t
	}
}

func WithMinHeight(h int64) Option {
	return func(s *Service) {
		s.minHeight = h
	}
}

func NewService(db adapter.Adapter, client cosmostxcollector.TXsCollecter, options ...Option) Service {
	s := Service{
		grace:     DefaultGrace,
		timeout:   DefaultTimeout,
		minHeight: DefaultMinHeight,
		db:        db,
		collector: cosmostxcollector.New(db, client),
	}

	for _, o := range options {
		o(&s)
	}

	return s
}

type Service struct {
	grace     time.Duration
	timeout   time.Duration
	minHeight int64
	db        adapter.Adapter
	collector cosmostxcollector.Collector
}

func (s Service) Run(ctx context.Context) error {
	log.WithFields(log.Fields{
		"backend": s.db.GetType(),
	}).Debugf("Initializing data backend...")

	if err := s.db.Init(ctx); err != nil {
		return err
	}

LOOP:
	for {
		// TODO: implement a retry policy on collection fail (increase grace)?
		if err := s.collect(ctx); err != nil {
			log.Errorf("error collecting TXs: %v", err)
		}

		select {
		case <-time.After(s.grace):
		case <-ctx.Done():
			break LOOP
		}
	}

	return nil
}

func (s Service) collect(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, s.timeout)
	defer cancel()

	fromHeight, err := s.db.GetLatestHeight(ctx)
	if err != nil {
		return err
	}

	if fromHeight < s.minHeight {
		fromHeight = s.minHeight
	} else {
		// Start from the height that follows the latest height
		fromHeight += 1
	}

	err = s.collector.Collect(ctx, fromHeight)
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

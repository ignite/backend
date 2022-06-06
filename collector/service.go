package collector

import (
	"context"
	"time"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosmetric"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosmetric/adapter"
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

func NewService(db adapter.Adapter, client cosmosmetric.TXsCollecter, options ...Option) Service {
	s := Service{
		grace:     DefaultGrace,
		timeout:   DefaultTimeout,
		minHeight: DefaultMinHeight,
		db:        db,
		collector: cosmosmetric.NewCollector(db, client),
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
	collector cosmosmetric.Collector
}

func (s Service) Run(ctx context.Context) error {
	if err := s.db.Init(ctx); err != nil {
		return err
	}

LOOP:
	for {
		if err := s.collect(ctx); err != nil {
			return err
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

	if err := s.collector.Collect(ctx, fromHeight); err != nil {
		return err
	}

	return nil
}

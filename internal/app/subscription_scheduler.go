package app

import (
	"context"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
	"github.com/rs/zerolog"
)

type SubscriptionScheduler struct {
	logger zerolog.Logger
	subs   *SubscriptionService
	repo   ports.SubscriptionRepository

	TickInterval time.Duration
	BatchSize    int
	Enqueue      bool
}

func NewSubscriptionScheduler(logger zerolog.Logger, subs *SubscriptionService, repo ports.SubscriptionRepository) *SubscriptionScheduler {
	return &SubscriptionScheduler{
		logger:       logger,
		subs:         subs,
		repo:         repo,
		TickInterval: 60 * time.Second,
		BatchSize:    10,
		Enqueue:      true,
	}
}

func (sch *SubscriptionScheduler) Run(ctx context.Context) {
	interval := sch.TickInterval
	if interval <= 0 {
		interval = 60 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			sch.logger.Info().Msg("subscription scheduler stopped")
			return
		case <-ticker.C:
			sch.tick(ctx)
		}
	}
}

func (sch *SubscriptionScheduler) tick(ctx context.Context) {
	if sch.subs == nil || sch.repo == nil {
		return
	}
	limit := sch.BatchSize
	if limit <= 0 {
		limit = 10
	}

	due, err := sch.repo.Due(ctx, time.Now().UTC(), limit)
	if err != nil {
		sch.logger.Error().Err(err).Msg("scheduler due query failed")
		return
	}
	if len(due) == 0 {
		return
	}

	for _, sub := range due {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, err := sch.subs.SyncOnce(ctx, sub.ID, sch.Enqueue)
		if err != nil {
			sch.logger.Warn().Err(err).Str("subscription_id", sub.ID).Msg("subscription sync failed")
		}
	}
}

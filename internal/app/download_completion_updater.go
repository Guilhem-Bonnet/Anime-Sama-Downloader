package app

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
	"github.com/rs/zerolog"
)

type DownloadCompletionUpdater struct {
	logger zerolog.Logger
	bus    ports.EventBus
	subs   ports.SubscriptionRepository
}

func NewDownloadCompletionUpdater(logger zerolog.Logger, bus ports.EventBus, subs ports.SubscriptionRepository) *DownloadCompletionUpdater {
	return &DownloadCompletionUpdater{logger: logger, bus: bus, subs: subs}
}

type downloadJobMeta struct {
	SubscriptionID string `json:"subscriptionId"`
	Episode        int    `json:"episode"`
	Source         string `json:"source,omitempty"`
}

func (u *DownloadCompletionUpdater) Run(ctx context.Context) {
	if u == nil || u.bus == nil || u.subs == nil {
		return
	}
	ch, cancel := u.bus.Subscribe()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			u.logger.Info().Msg("download completion updater stopped")
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			u.handleEvent(ctx, evt)
		}
	}
}

func (u *DownloadCompletionUpdater) handleEvent(ctx context.Context, evt ports.Event) {
	if evt.Topic != "job.completed" {
		return
	}

	var job JobDTO
	if err := json.Unmarshal(evt.Payload, &job); err != nil {
		return
	}
	if job.Type != "download" {
		return
	}

	meta := downloadJobMeta{}
	if len(job.Params) > 0 {
		_ = json.Unmarshal(job.Params, &meta)
	}
	meta.SubscriptionID = strings.TrimSpace(meta.SubscriptionID)
	if meta.SubscriptionID == "" || meta.Episode <= 0 {
		return
	}

	updated, err := u.subs.MarkDownloadedEpisodeMax(ctx, meta.SubscriptionID, meta.Episode)
	if err != nil {
		u.logger.Warn().Err(err).Str("subscription_id", meta.SubscriptionID).Msg("failed to mark episode downloaded")
		return
	}

	// Best-effort notification.
	if u.bus != nil {
		b, _ := json.Marshal(toSubscriptionDTO(updated))
		if len(b) > 0 {
			u.bus.Publish("subscription.downloaded", b)
		}
	}
}

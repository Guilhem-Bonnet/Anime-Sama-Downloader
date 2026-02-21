package app

import (
	"context"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

func TestEpisodeResolver_NewEpisodeResolver(t *testing.T) {
	resolver := NewEpisodeResolver()
	if resolver == nil {
		t.Error("NewEpisodeResolver should return non-nil resolver")
	}
}

func TestEpisodeResolver_Resolve_ReturnsStructuredResult(t *testing.T) {
	// This test verifies the EpisodeResolver structure
	// Full integration tests would require a mock server
	resolver := NewEpisodeResolver()

	// Test with invalid URL to verify error handling
	ctx := context.Background()
	_, err := resolver.Resolve(ctx, "not-a-valid-url", "auto")

	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestComputeNextCheck_NewEpisodesAvailable(t *testing.T) {
	// When new episodes are available, should check sooner
	sub := domain.Subscription{
		LastScheduledEpisode: 5,
	}
	maxAvailable := 10

	nextCheck := ComputeNextCheck(sub, maxAvailable)
	expectedInterval := 10 * time.Minute

	now := time.Now().UTC()
	diff := nextCheck.Sub(now)

	// Allow 5 seconds tolerance for test execution time
	if diff < expectedInterval-5*time.Second || diff > expectedInterval+5*time.Second {
		t.Errorf("Expected next check in ~%v, got %v", expectedInterval, diff)
	}
}

func TestComputeNextCheck_AllEpisodesScheduled(t *testing.T) {
	// When all episodes are scheduled, should check less frequently
	sub := domain.Subscription{
		LastScheduledEpisode: 10,
	}
	maxAvailable := 10

	nextCheck := ComputeNextCheck(sub, maxAvailable)
	expectedInterval := 2 * time.Hour

	now := time.Now().UTC()
	diff := nextCheck.Sub(now)

	// Allow 5 seconds tolerance for test execution time
	if diff < expectedInterval-5*time.Second || diff > expectedInterval+5*time.Second {
		t.Errorf("Expected next check in ~%v, got %v", expectedInterval, diff)
	}
}

func TestComputeNextCheck_MoreScheduledThanAvailable(t *testing.T) {
	// Edge case: scheduled more than available (shouldn't happen but handle gracefully)
	sub := domain.Subscription{
		LastScheduledEpisode: 15,
	}
	maxAvailable := 10

	nextCheck := ComputeNextCheck(sub, maxAvailable)
	expectedInterval := 2 * time.Hour

	now := time.Now().UTC()
	diff := nextCheck.Sub(now)

	// Should use the longer interval since not behind
	if diff < expectedInterval-5*time.Second || diff > expectedInterval+5*time.Second {
		t.Errorf("Expected next check in ~%v, got %v", expectedInterval, diff)
	}
}

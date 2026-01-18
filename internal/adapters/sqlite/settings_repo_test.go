package sqlite

import (
	"context"
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

func TestSettingsRepository_DefaultsAndPersist(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := NewSettingsRepository(db.SQL)

	got, err := repo.Get(ctx)
	if err != nil {
		t.Fatalf("Get(default): %v", err)
	}
	if got.Destination == "" {
		t.Fatalf("expected default Destination, got empty")
	}

	want := domain.DefaultSettings()
	want.Destination = "/tmp/videos"
	want.OutputNamingMode = domain.OutputNamingMediaServer
	want.MaxWorkers = 3
	want.MaxConcurrentDownloads = 6

	updated, err := repo.Put(ctx, want)
	if err != nil {
		t.Fatalf("Put: %v", err)
	}
	if updated.Destination != want.Destination {
		t.Fatalf("Destination: want %q, got %q", want.Destination, updated.Destination)
	}
	if updated.OutputNamingMode != want.OutputNamingMode {
		t.Fatalf("OutputNamingMode: want %q, got %q", want.OutputNamingMode, updated.OutputNamingMode)
	}
	if updated.MaxWorkers != want.MaxWorkers {
		t.Fatalf("MaxWorkers: want %d, got %d", want.MaxWorkers, updated.MaxWorkers)
	}
	if updated.MaxConcurrentDownloads != want.MaxConcurrentDownloads {
		t.Fatalf("MaxConcurrentDownloads: want %d, got %d", want.MaxConcurrentDownloads, updated.MaxConcurrentDownloads)
	}

	got2, err := repo.Get(ctx)
	if err != nil {
		t.Fatalf("Get(after Put): %v", err)
	}
	if got2.Destination != want.Destination {
		t.Fatalf("Destination after Put: want %q, got %q", want.Destination, got2.Destination)
	}
}

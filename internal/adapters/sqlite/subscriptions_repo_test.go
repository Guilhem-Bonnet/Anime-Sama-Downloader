package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

func setupTestDB(t *testing.T) *SubscriptionsRepository {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return NewSubscriptionsRepository(db.SQL)
}

func TestSubscriptionsRepository_Create_Success(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	sub := domain.Subscription{
		ID:                    "sub-test-1",
		BaseURL:               "https://anime-sama.si/catalogue/test",
		Label:                 "Test Anime",
		Player:                "default",
		LastScheduledEpisode:  0,
		LastDownloadedEpisode: 0,
		LastAvailableEpisode:  0,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	result, err := repo.Create(ctx, sub)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if result.ID != "sub-test-1" {
		t.Errorf("expected ID 'sub-test-1', got %q", result.ID)
	}
	if result.Label != "Test Anime" {
		t.Errorf("expected label 'Test Anime', got %q", result.Label)
	}
}

func TestSubscriptionsRepository_Create_DuplicateBaseURL(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	sub1 := domain.Subscription{
		ID:        "sub-1",
		BaseURL:   "https://anime-sama.si/catalogue/test",
		Label:     "First",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := repo.Create(ctx, sub1)
	if err != nil {
		t.Fatalf("Create first: %v", err)
	}

	sub2 := domain.Subscription{
		ID:        "sub-2",
		BaseURL:   "https://anime-sama.si/catalogue/test", // same BaseURL
		Label:     "Second",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = repo.Create(ctx, sub2)
	if err != ports.ErrConflict {
		t.Errorf("expected ErrConflict, got %v", err)
	}
}

func TestSubscriptionsRepository_Get_Success(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	sub := domain.Subscription{
		ID:        "sub-get-1",
		BaseURL:   "https://anime-sama.si/catalogue/test",
		Label:     "Get Test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := repo.Create(ctx, sub)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	result, err := repo.Get(ctx, "sub-get-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if result.ID != "sub-get-1" {
		t.Errorf("expected ID 'sub-get-1', got %q", result.ID)
	}
}

func TestSubscriptionsRepository_List_Success(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		sub := domain.Subscription{
			ID:        "sub-" + string(rune(i+'0')),
			BaseURL:   "https://anime-sama.si/catalogue/test" + string(rune(i+'0')),
			Label:     "Test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err := repo.Create(ctx, sub)
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
	}

	result, err := repo.List(ctx, 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 subscriptions, got %d", len(result))
	}
}

func TestSubscriptionsRepository_Update_Success(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	sub := domain.Subscription{
		ID:        "sub-update",
		BaseURL:   "https://anime-sama.si/catalogue/test",
		Label:     "Original",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := repo.Create(ctx, sub)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Update
	sub.Label = "Updated"
	sub.LastDownloadedEpisode = 5
	_, err = repo.Update(ctx, sub)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	result, err := repo.Get(ctx, "sub-update")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if result.Label != "Updated" {
		t.Errorf("expected label 'Updated', got %q", result.Label)
	}
	if result.LastDownloadedEpisode != 5 {
		t.Errorf("expected LastDownloadedEpisode 5, got %d", result.LastDownloadedEpisode)
	}
}

func TestSubscriptionsRepository_Delete_Success(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	sub := domain.Subscription{
		ID:        "sub-delete",
		BaseURL:   "https://anime-sama.si/catalogue/test",
		Label:     "To Delete",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := repo.Create(ctx, sub)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	err = repo.Delete(ctx, "sub-delete")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	// Verify deleted
	result, _ := repo.Get(ctx, "sub-delete")
	if result.ID != "" {
		t.Errorf("expected empty subscription after delete, got %v", result)
	}
}

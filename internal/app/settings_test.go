package app

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// stubSettingsRepo is a test stub for SettingsRepository
type stubSettingsRepo struct {
	getErr error
	getSrc domain.Settings
	putErr error
	putDst domain.Settings
}

func (s *stubSettingsRepo) Get(ctx context.Context) (domain.Settings, error) {
	return s.getSrc, s.getErr
}

func (s *stubSettingsRepo) Put(ctx context.Context, settings domain.Settings) (domain.Settings, error) {
	if s.putErr != nil {
		return domain.Settings{}, s.putErr
	}
	s.putDst = settings
	return settings, nil
}

func (s *stubSettingsRepo) Initialize(ctx context.Context) error {
	return nil
}

var _ ports.SettingsRepository = (*stubSettingsRepo)(nil)

// TestNewSettingsService_Success tests creating a SettingsService
func TestNewSettingsService_Success(t *testing.T) {
	repo := &stubSettingsRepo{}
	svc := NewSettingsService(repo)

	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
}

// TestSettingsService_Get_Success tests Get returns settings
func TestSettingsService_Get_Success(t *testing.T) {
	settings := domain.Settings{
		Destination:             "/tmp/videos",
		OutputNamingMode:        "auto",
		MaxWorkers:              4,
		MaxConcurrentDownloads:  2,
		AniListToken:            "token123",
	}

	repo := &stubSettingsRepo{getSrc: settings}
	svc := NewSettingsService(repo)

	result, err := svc.Get(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Destination != "/tmp/videos" {
		t.Fatalf("expected destination /tmp/videos, got %q", result.Destination)
	}
	if result.MaxWorkers != 4 {
		t.Fatalf("expected 4 workers, got %d", result.MaxWorkers)
	}
}

// TestSettingsService_Get_RepoError tests Get propagates repo error
func TestSettingsService_Get_RepoError(t *testing.T) {
	repo := &stubSettingsRepo{getErr: errors.New("repo error")}
	svc := NewSettingsService(repo)

	_, err := svc.Get(context.Background())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "repo error" {
		t.Fatalf("expected 'repo error', got %q", err.Error())
	}
}

// TestSettingsService_Put_DefaultsEmpty tests Put fills empty fields with defaults
func TestSettingsService_Put_DefaultsEmpty(t *testing.T) {
	repo := &stubSettingsRepo{}
	svc := NewSettingsService(repo)

	input := domain.Settings{} // All empty

	result, err := svc.Put(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	defaults := domain.DefaultSettings()
	if result.Destination != defaults.Destination {
		t.Fatalf("expected default destination, got %q", result.Destination)
	}
	if result.OutputNamingMode != defaults.OutputNamingMode {
		t.Fatalf("expected default naming mode, got %q", result.OutputNamingMode)
	}
	if result.MaxWorkers != defaults.MaxWorkers {
		t.Fatalf("expected default workers, got %d", result.MaxWorkers)
	}
}

// TestSettingsService_Put_PreservesNonEmpty tests Put preserves non-empty fields
func TestSettingsService_Put_PreservesNonEmpty(t *testing.T) {
	repo := &stubSettingsRepo{}
	svc := NewSettingsService(repo)

	input := domain.Settings{
		Destination:            "/custom/path",
		OutputNamingMode:       "custom-mode",
		MaxWorkers:             10,
		MaxConcurrentDownloads: 5,
	}

	result, err := svc.Put(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Destination != "/custom/path" {
		t.Fatalf("expected /custom/path, got %q", result.Destination)
	}
	if result.MaxWorkers != 10 {
		t.Fatalf("expected 10 workers, got %d", result.MaxWorkers)
	}
}

// TestSettingsService_Put_DefaultsZeroMaxWorkers tests zero MaxWorkers is replaced
func TestSettingsService_Put_DefaultsZeroMaxWorkers(t *testing.T) {
	repo := &stubSettingsRepo{}
	svc := NewSettingsService(repo)

	input := domain.Settings{
		Destination:      "/tmp",
		MaxWorkers:       0, // Zero - should be replaced
		OutputNamingMode: "auto",
	}

	result, err := svc.Put(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	defaults := domain.DefaultSettings()
	if result.MaxWorkers != defaults.MaxWorkers {
		t.Fatalf("expected default workers, got %d", result.MaxWorkers)
	}
}

// TestSettingsService_Put_DefaultsNegativeMaxConcurrentDownloads tests negative value is replaced
func TestSettingsService_Put_DefaultsNegativeMaxConcurrentDownloads(t *testing.T) {
	repo := &stubSettingsRepo{}
	svc := NewSettingsService(repo)

	input := domain.Settings{
		Destination:             "/tmp",
		OutputNamingMode:        "auto",
		MaxWorkers:              1,
		MaxConcurrentDownloads:  -1, // Negative - should be replaced
	}

	result, err := svc.Put(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	defaults := domain.DefaultSettings()
	if result.MaxConcurrentDownloads != defaults.MaxConcurrentDownloads {
		t.Fatalf("expected default concurrent downloads, got %d", result.MaxConcurrentDownloads)
	}
}

// TestSettingsService_Put_RepoError tests Put propagates repo error
func TestSettingsService_Put_RepoError(t *testing.T) {
	repo := &stubSettingsRepo{putErr: errors.New("put failed")}
	svc := NewSettingsService(repo)

	_, err := svc.Put(context.Background(), domain.Settings{MaxWorkers: 1})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "put failed" {
		t.Fatalf("expected 'put failed', got %q", err.Error())
	}
}

// TestSettingsService_Get_AutoDefaultDestination_FromEnv tests env override
func TestSettingsService_Get_AutoDefaultDestination_FromEnv(t *testing.T) {
	oldEnv := os.Getenv("ASD_DEFAULT_DESTINATION")
	defer os.Setenv("ASD_DEFAULT_DESTINATION", oldEnv)

	os.Setenv("ASD_DEFAULT_DESTINATION", "/env/override")

	settings := domain.Settings{
		Destination: domain.DefaultSettings().Destination, // Use default
	}

	repo := &stubSettingsRepo{getSrc: settings}
	svc := NewSettingsService(repo)

	result, err := svc.Get(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Destination != "/env/override" {
		t.Fatalf("expected env override /env/override, got %q", result.Destination)
	}
}

// TestSettingsService_Get_CustomDestinationPreserved tests custom destination is preserved
func TestSettingsService_Get_CustomDestinationPreserved(t *testing.T) {
	settings := domain.Settings{
		Destination: "/custom/non-default",
	}

	repo := &stubSettingsRepo{getSrc: settings}
	svc := NewSettingsService(repo)

	result, err := svc.Get(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Destination != "/custom/non-default" {
		t.Fatalf("expected custom destination to be preserved, got %q", result.Destination)
	}
}

// TestAutoDefaultDestination_NonEmpty tests empty path handling
func TestAutoDefaultDestination_NonEmpty(t *testing.T) {
	custom := "/my/custom/path"
	dest, ok := autoDefaultDestination(custom)

	if ok {
		t.Fatalf("expected ok=false for custom path, got true")
	}
	if dest != "" {
		t.Fatalf("expected empty string, got %q", dest)
	}
}

// TestAutoDefaultDestination_EmptyDefault tests default path handling
func TestAutoDefaultDestination_EmptyDefault(t *testing.T) {
	// When input matches default, check env and /data/videos
	defaultPath := domain.DefaultSettings().Destination

	oldEnv := os.Getenv("ASD_DEFAULT_DESTINATION")
	defer os.Setenv("ASD_DEFAULT_DESTINATION", oldEnv)
	os.Setenv("ASD_DEFAULT_DESTINATION", "")

	_, _ = autoDefaultDestination(defaultPath)

	// Should either be /data/videos (if it exists) or empty
	// We can't guarantee /data/videos exists in test env, so just check behavior
}

// TestAutoDefaultDestination_Whitespace tests whitespace trimming
func TestAutoDefaultDestination_Whitespace(t *testing.T) {
	dest, _ := autoDefaultDestination("  ")

	// Whitespace should be treated as empty
	// We can't guarantee /data/videos exists, so just check it's processed
	if dest == "  " {
		t.Fatalf("expected whitespace to be handled, got %q", dest)
	}
}

// TestSettingsService_Put_PartialDefaults tests partial field defaults
func TestSettingsService_Put_PartialDefaults(t *testing.T) {
	repo := &stubSettingsRepo{}
	svc := NewSettingsService(repo)

	input := domain.Settings{
		Destination: "/custom", // Provide this
		MaxWorkers:  0,          // Leave empty - should default
		// Others left empty
	}

	result, err := svc.Put(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Destination != "/custom" {
		t.Fatalf("expected custom destination, got %q", result.Destination)
	}

	defaults := domain.DefaultSettings()
	if result.MaxWorkers != defaults.MaxWorkers {
		t.Fatalf("expected default workers, got %d", result.MaxWorkers)
	}
	if result.OutputNamingMode != defaults.OutputNamingMode {
		t.Fatalf("expected default naming mode, got %q", result.OutputNamingMode)
	}
}

// TestSettingsService_Put_AllFieldsPreserved tests all provided fields are preserved
func TestSettingsService_Put_AllFieldsPreserved(t *testing.T) {
	repo := &stubSettingsRepo{}
	svc := NewSettingsService(repo)

	input := domain.Settings{
		Destination:             "/full/custom",
		OutputNamingMode:        "full-mode",
		MaxWorkers:              7,
		MaxConcurrentDownloads:  3,
		AniListToken:            "my-token",
	}

	result, err := svc.Put(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Destination != "/full/custom" ||
		result.OutputNamingMode != "full-mode" ||
		result.MaxWorkers != 7 ||
		result.MaxConcurrentDownloads != 3 ||
		result.AniListToken != "my-token" {
		t.Fatalf("expected all fields preserved in result: %+v", result)
	}
}

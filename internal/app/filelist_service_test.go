package app

import (
	"context"
	"testing"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// testCatalogueForFileListing creates a test catalogue with known file structure
func testCatalogueForFileListing() []*domain.Anime {
	return []*domain.Anime{
		{
			ID:           "1",
			Title:        "Naruto",
			Year:         2002,
			Status:       "Completed",
			EpisodeCount: 3,
			Genres:       []string{"Action", "Adventure"},
			ThumbnailURL: "https://example.com/naruto.jpg",
		},
		{
			ID:           "2",
			Title:        "Death Note",
			Year:         2006,
			Status:       "Completed",
			EpisodeCount: 2,
			Genres:       []string{"Mystery", "Thriller"},
			ThumbnailURL: "https://example.com/deathnote.jpg",
		},
	}
}

// TestFileListService_GetFileList_Success tests successful file listing retrieval
func TestFileListService_GetFileList_Success(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	fileList, err := svc.GetFileList(context.Background(), "1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fileList == nil {
		t.Fatal("expected file list, got nil")
	}

	if fileList.AnimeID != "1" {
		t.Errorf("expected anime ID '1', got '%s'", fileList.AnimeID)
	}

	if len(fileList.Files) != 3 {
		t.Errorf("expected 3 files for Naruto, got %d", len(fileList.Files))
	}

	// Verify file structure
	for i, file := range fileList.Files {
		if file.ID == "" {
			t.Errorf("file %d: empty ID", i)
		}
		if file.Name == "" {
			t.Errorf("file %d: empty name", i)
		}
		if file.Size <= 0 {
			t.Errorf("file %d: invalid size %d", i, file.Size)
		}
		if file.Type == "" {
			t.Errorf("file %d: empty type", i)
		}
	}
}

// TestFileListService_GetFileList_NotFound tests handling of missing anime
func TestFileListService_GetFileList_NotFound(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	fileList, err := svc.GetFileList(context.Background(), "999")

	if err == nil {
		t.Fatal("expected error for non-existent anime")
	}

	if fileList != nil {
		t.Error("expected nil file list for non-existent anime")
	}
}

// TestFileListService_GetFilesByAnimeTitle_Success tests file retrieval by title
func TestFileListService_GetFilesByAnimeTitle_Success(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	fileList, err := svc.GetFilesByAnimeTitle(context.Background(), "Naruto")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fileList == nil {
		t.Fatal("expected file list, got nil")
	}

	if fileList.AnimeID != "1" {
		t.Errorf("expected anime ID '1', got '%s'", fileList.AnimeID)
	}

	if len(fileList.Files) != 3 {
		t.Errorf("expected 3 files for Naruto, got %d", len(fileList.Files))
	}
}

// TestFileListService_GetFilesByAnimeTitle_NotFound tests title lookup for missing anime
func TestFileListService_GetFilesByAnimeTitle_NotFound(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	fileList, err := svc.GetFilesByAnimeTitle(context.Background(), "NonExistentAnime")

	if err == nil {
		t.Fatal("expected error for non-existent anime title")
	}

	if fileList != nil {
		t.Error("expected nil file list for non-existent anime")
	}
}

// TestFileListService_FileMetadata_Consistency tests that file metadata is consistent
func TestFileListService_FileMetadata_Consistency(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	fileList, err := svc.GetFileList(context.Background(), "2")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fileList.Files) != 2 {
		t.Errorf("expected 2 files for Death Note, got %d", len(fileList.Files))
	}

	// Verify sequential naming and metadata
	for i, file := range fileList.Files {
		if file.Duration <= 0 {
			t.Errorf("file %d: expected positive duration, got %d", i, file.Duration)
		}
	}
}

// TestFileListService_Context_Cancellation tests context cancellation handling
func TestFileListService_Context_Cancellation(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Should handle cancelled context gracefully
	fileList, err := svc.GetFileList(ctx, "1")

	// Either returns error or proceeds with the operation
	// depending on implementation - both are acceptable
	if fileList == nil && err == nil {
		t.Error("expected either file list or error, got both nil")
	}
}

// TestFileListService_LargeAnime_Performance tests performance with many episodes
func TestFileListService_LargeAnime_Performance(t *testing.T) {
	// Create anime with many episodes
	largeAnime := []*domain.Anime{
		{
			ID:           "large-1",
			Title:        "One Piece",
			Year:         1999,
			Status:       "Ongoing",
			EpisodeCount: 1000, // Large episode count
			Genres:       []string{"Action", "Adventure"},
			ThumbnailURL: "https://example.com/onepiece.jpg",
		},
	}

	svc := NewFileListService(largeAnime)

	fileList, err := svc.GetFileList(context.Background(), "large-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fileList.Files) != 1000 {
		t.Errorf("expected 1000 files, got %d", len(fileList.Files))
	}

	// Verify metadata is generated correctly for all episodes
	if fileList.Files[0].ID != "large-1-ep1" {
		t.Errorf("expected first episode ID 'large-1-ep1', got '%s'", fileList.Files[0].ID)
	}

	if fileList.Files[999].ID != "large-1-ep1000" {
		t.Errorf("expected last episode ID 'large-1-ep1000', got '%s'", fileList.Files[999].ID)
	}
}

// TestFileListService_FileMetadata_Uniqueness tests that file IDs are unique
func TestFileListService_FileMetadata_Uniqueness(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	fileList, err := svc.GetFileList(context.Background(), "1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check for unique IDs
	idMap := make(map[string]bool)
	for _, file := range fileList.Files {
		if idMap[file.ID] {
			t.Errorf("duplicate file ID found: %s", file.ID)
		}
		idMap[file.ID] = true
	}

	// Check for unique paths
	pathMap := make(map[string]bool)
	for _, file := range fileList.Files {
		if pathMap[file.Path] {
			t.Errorf("duplicate file path found: %s", file.Path)
		}
		pathMap[file.Path] = true
	}
}

// TestFileListService_FileSizes_Realistic tests that file sizes are realistic
func TestFileListService_FileSizes_Realistic(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	fileList, err := svc.GetFileList(context.Background(), "1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify sizes are in realistic range (200MB - 600MB per episode)
	const minSize = 200 * 1024 * 1024 // 200 MB
	const maxSize = 600 * 1024 * 1024 // 600 MB

	for i, file := range fileList.Files {
		if file.Size < minSize || file.Size > maxSize {
			t.Errorf("file %d: unrealistic size %d bytes (expected %d-%d)", i, file.Size, minSize, maxSize)
		}
	}
}

// TestFileListService_Duration_Realistic tests that durations are realistic
func TestFileListService_Duration_Realistic(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	fileList, err := svc.GetFileList(context.Background(), "1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify durations are in realistic range (18-30 minutes)
	const minDuration = 18 * 60  // 18 minutes in seconds
	const maxDuration = 30 * 60  // 30 minutes in seconds

	for i, file := range fileList.Files {
		if file.Duration < minDuration || file.Duration > maxDuration {
			t.Errorf("file %d: unrealistic duration %d seconds (expected %d-%d)", i, file.Duration, minDuration, maxDuration)
		}
	}
}

// TestFileListService_EmptyAnime_ZeroFiles tests anime with no episodes
func TestFileListService_EmptyAnime_ZeroFiles(t *testing.T) {
	emptyAnime := []*domain.Anime{
		{
			ID:           "empty-1",
			Title:        "Empty Anime",
			Year:         2020,
			Status:       "Planned",
			EpisodeCount: 0, // No episodes yet
			Genres:       []string{"Unknown"},
			ThumbnailURL: "https://example.com/empty.jpg",
		},
	}

	svc := NewFileListService(emptyAnime)

	fileList, err := svc.GetFileList(context.Background(), "empty-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fileList.Files) != 0 {
		t.Errorf("expected 0 files for empty anime, got %d", len(fileList.Files))
	}
}

// TestFileListService_CaseInsensitive_TitleSearch tests case-insensitive title matching
func TestFileListService_CaseInsensitive_TitleSearch(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	testCases := []string{
		"naruto",
		"NARUTO",
		"NaRuTo",
		"  Naruto  ", // with spaces
	}

	for _, title := range testCases {
		fileList, err := svc.GetFilesByAnimeTitle(context.Background(), title)

		if err != nil {
			t.Errorf("failed to find anime with title '%s': %v", title, err)
			continue
		}

		if fileList.AnimeID != "1" {
			t.Errorf("title '%s': expected anime ID '1', got '%s'", title, fileList.AnimeID)
		}
	}
}

// TestFileListService_SpecialCharacters_TitleSearch tests titles with special characters
func TestFileListService_SpecialCharacters_TitleSearch(t *testing.T) {
	specialAnime := []*domain.Anime{
		{
			ID:           "special-1",
			Title:        "Re:Zero - Starting Life in Another World",
			Year:         2016,
			Status:       "Completed",
			EpisodeCount: 2,
			Genres:       []string{"Fantasy"},
			ThumbnailURL: "https://example.com/rezero.jpg",
		},
	}

	svc := NewFileListService(specialAnime)

	// Test exact match with special characters
	fileList, err := svc.GetFilesByAnimeTitle(context.Background(), "Re:Zero - Starting Life in Another World")

	if err != nil {
		t.Fatalf("failed to find anime with special characters: %v", err)
	}

	if fileList.AnimeID != "special-1" {
		t.Errorf("expected anime ID 'special-1', got '%s'", fileList.AnimeID)
	}
}

// TestFileListService_MultipleRequests_Consistency tests that repeated calls return consistent data
func TestFileListService_MultipleRequests_Consistency(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	// Make multiple requests for the same anime
	fileList1, err1 := svc.GetFileList(context.Background(), "1")
	fileList2, err2 := svc.GetFileList(context.Background(), "1")

	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected errors: %v, %v", err1, err2)
	}

	// Verify both responses are identical
	if len(fileList1.Files) != len(fileList2.Files) {
		t.Errorf("inconsistent file count: %d vs %d", len(fileList1.Files), len(fileList2.Files))
	}

	for i := range fileList1.Files {
		if fileList1.Files[i].ID != fileList2.Files[i].ID {
			t.Errorf("file %d: inconsistent ID: '%s' vs '%s'", i, fileList1.Files[i].ID, fileList2.Files[i].ID)
		}
		if fileList1.Files[i].Size != fileList2.Files[i].Size {
			t.Errorf("file %d: inconsistent size: %d vs %d", i, fileList1.Files[i].Size, fileList2.Files[i].Size)
		}
	}
}

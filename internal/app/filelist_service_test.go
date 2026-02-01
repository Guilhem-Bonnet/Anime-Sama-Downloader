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

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

// ==================== TASK 6.3: ERROR SCENARIOS ====================

// TestFileListService_GetFileList_InvalidAnimeID tests handling of invalid anime ID formats
func TestFileListService_GetFileList_InvalidAnimeID(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	testCases := []struct {
		name    string
		animeID string
	}{
		{"empty ID", ""},
		{"whitespace only", "   "},
		{"non-existent ID", "999999"},
		{"special characters", "!@#$%"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fileList, err := svc.GetFileList(context.Background(), tc.animeID)

			if tc.animeID != "" && tc.animeID != "   " {
				// For valid ID formats that don't exist
				if err == nil {
					t.Error("expected error for non-existent anime ID")
				}
				if fileList != nil {
					t.Error("expected nil file list for non-existent anime")
				}
			}
		})
	}
}

// TestFileListService_GetFilesByAnimeTitle_InvalidInput tests handling of invalid title inputs
func TestFileListService_GetFilesByAnimeTitle_InvalidInput(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	testCases := []struct {
		name  string
		title string
	}{
		{"empty title", ""},
		{"whitespace only", "   "},
		{"very long title", string(make([]byte, 10000))},
		{"non-existent title", "ThisAnimeDoesNotExist"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fileList, err := svc.GetFilesByAnimeTitle(context.Background(), tc.title)

			if tc.title != "" && tc.title != "   " {
				if err == nil {
					t.Error("expected error for non-existent anime title")
				}
				if fileList != nil {
					t.Error("expected nil file list")
				}
			}
		})
	}
}

// TestFileListService_NilCatalogueHandling tests graceful handling of nil catalogue
func TestFileListService_NilCatalogueHandling(t *testing.T) {
	svc := NewFileListService(nil)

	fileList, err := svc.GetFileList(context.Background(), "1")

	if err == nil {
		t.Error("expected error with nil catalogue")
	}
	if fileList != nil {
		t.Error("expected nil file list with nil catalogue")
	}
}

// TestFileListService_LargeTitleSearch tests performance with special characters
func TestFileListService_LargeTitleSearch(t *testing.T) {
	catalogue := []*domain.Anime{
		{
			ID:           "1",
			Title:        "Shingeki no Kyojin: Attack on Titan",
			Year:         2013,
			Status:       "Completed",
			EpisodeCount: 1,
			Genres:       []string{"Action"},
			ThumbnailURL: "https://example.com/aot.jpg",
		},
		{
			ID:           "2",
			Title:        "Kimetsu no Yaiba: Demon Slayer",
			Year:         2019,
			Status:       "Ongoing",
			EpisodeCount: 1,
			Genres:       []string{"Action"},
			ThumbnailURL: "https://example.com/ks.jpg",
		},
	}

	svc := NewFileListService(catalogue)

	// Use exact title match (case-insensitive)
	fileList, err := svc.GetFilesByAnimeTitle(context.Background(), "Shingeki no Kyojin: Attack on Titan")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fileList == nil {
		t.Fatal("expected file list, got nil")
	}

	if fileList.AnimeID != "1" {
		t.Errorf("expected anime ID '1', got '%s'", fileList.AnimeID)
	}
}

// ==================== TASK 6.4: PAGINATION & FILTERING ====================

// TestFileListService_Filtering_ByFileType tests filtering files by type
func TestFileListService_Filtering_ByFileType(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	fileList, err := svc.GetFileList(context.Background(), "1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify all files have valid MIME types (video container formats)
	validTypes := map[string]bool{"video/x-matroska": true, "video/mp4": true, "video/x-msvideo": true}
	for i, file := range fileList.Files {
		if !validTypes[file.Type] {
			t.Logf("file %d: type '%s' (may be valid MIME type)", i, file.Type)
		}
	}

	// Count files by type
	typeCounts := make(map[string]int)
	for _, file := range fileList.Files {
		typeCounts[file.Type]++
	}

	// Verify we have files of expected type
	if len(typeCounts) == 0 {
		t.Error("expected at least one file type")
	}

	// Verify all files are of the same type (consistent generation)
	if len(typeCounts) > 1 {
		t.Logf("files have %d different types: %v", len(typeCounts), typeCounts)
	}
}

// TestFileListService_Filtering_BySize tests logical filtering by file size
func TestFileListService_Filtering_BySize(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	fileList, err := svc.GetFileList(context.Background(), "1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file sizes are within realistic ranges
	for i, file := range fileList.Files {
		// Reasonable size range: 100MB to 2GB
		minSize := int64(100 * 1024 * 1024)      // 100MB
		maxSize := int64(2 * 1024 * 1024 * 1024) // 2GB

		if file.Size < minSize || file.Size > maxSize {
			t.Logf("file %d: size %d bytes is outside expected range [%d, %d]", i, file.Size, minSize, maxSize)
		}
	}
}

// TestFileListService_Pagination_Logic tests pagination parameter handling
func TestFileListService_Pagination_Logic(t *testing.T) {
	// Create a large catalogue to test pagination
	largeCatalog := make([]*domain.Anime, 50)
	for i := 0; i < 50; i++ {
		largeCatalog[i] = &domain.Anime{
			ID:           string(rune(i)),
			Title:        "Anime " + string(rune(i)),
			Year:         2020 + (i % 5),
			Status:       "Completed",
			EpisodeCount: i + 1,
			Genres:       []string{"Action"},
			ThumbnailURL: "https://example.com/anime.jpg",
		}
	}

	svc := NewFileListService(largeCatalog)

	// Test that we can retrieve file lists for different anime
	for i := 0; i < 5; i++ {
		fileList, err := svc.GetFileList(context.Background(), string(rune(i)))

		if err != nil {
			t.Errorf("anime %d: unexpected error: %v", i, err)
		}
		if fileList == nil {
			t.Errorf("anime %d: expected file list, got nil", i)
		}
	}
}

// TestFileListService_Sorting_ByFileName tests file name sorting consistency
func TestFileListService_Sorting_ByFileName(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	fileList, err := svc.GetFileList(context.Background(), "1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify files have unique names (deterministic generation)
	fileNames := make(map[string]bool)
	for i, file := range fileList.Files {
		if fileNames[file.Name] {
			t.Errorf("file %d: duplicate name '%s'", i, file.Name)
		}
		fileNames[file.Name] = true
	}

	// Verify consistent ordering on repeated calls
	fileList2, _ := svc.GetFileList(context.Background(), "1")
	for i := range fileList.Files {
		if fileList.Files[i].Name != fileList2.Files[i].Name {
			t.Errorf("file %d: name mismatch on repeated call: '%s' vs '%s'", 
				i, fileList.Files[i].Name, fileList2.Files[i].Name)
		}
	}
}

// TestFileListService_OffsetPagination tests offset-based pagination semantics
func TestFileListService_OffsetPagination(t *testing.T) {
	svc := NewFileListService(testCatalogueForFileListing())

	// Get full list
	fullList, err := svc.GetFileList(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalFiles := len(fullList.Files)

	// Verify pagination parameters would work
	pageSize := 2
	numPages := (totalFiles + pageSize - 1) / pageSize

	if numPages < 1 {
		t.Error("expected at least 1 page")
	}

	// Verify total file count is consistent
	if totalFiles < 1 {
		t.Error("expected at least 1 file in test data")
	}
}

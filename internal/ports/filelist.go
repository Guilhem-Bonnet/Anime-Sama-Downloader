package ports

import (
	"context"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// FileListService defines the interface for retrieving file listings
type FileListService interface {
	// GetFileList returns the list of files for a given anime
	// Returns error if anime not found or network failure occurs
	GetFileList(ctx context.Context, animeID string) (*domain.FileList, error)

	// GetFilesByAnimeTitle returns files by anime title
	GetFilesByAnimeTitle(ctx context.Context, title string) (*domain.FileList, error)
}

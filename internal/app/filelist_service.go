package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

// FileListServiceImpl implements the FileListService interface
type FileListServiceImpl struct {
	catalogue []*domain.Anime
}

// NewFileListService creates a new file listing service
func NewFileListService(catalogue []*domain.Anime) ports.FileListService {
	return &FileListServiceImpl{
		catalogue: catalogue,
	}
}

// GetFileList returns the list of files for a given anime ID
func (s *FileListServiceImpl) GetFileList(ctx context.Context, animeID string) (*domain.FileList, error) {
	// Find anime in catalogue
	var anime *domain.Anime
	for _, a := range s.catalogue {
		if a.ID == animeID {
			anime = a
			break
		}
	}

	if anime == nil {
		return nil, fmt.Errorf("anime not found: %s", animeID)
	}

	return s.generateFileList(anime), nil
}

// GetFilesByAnimeTitle returns the list of files for a given anime title
func (s *FileListServiceImpl) GetFilesByAnimeTitle(ctx context.Context, title string) (*domain.FileList, error) {
	// Find anime by title (case-insensitive)
	normalizedTitle := strings.ToLower(strings.TrimSpace(title))
	var anime *domain.Anime
	for _, a := range s.catalogue {
		if strings.ToLower(a.Title) == normalizedTitle {
			anime = a
			break
		}
	}

	if anime == nil {
		return nil, fmt.Errorf("anime not found: %s", title)
	}

	return s.generateFileList(anime), nil
}

// generateFileList generates a file list for an anime with realistic metadata
func (s *FileListServiceImpl) generateFileList(anime *domain.Anime) *domain.FileList {
	files := make([]domain.File, anime.EpisodeCount)

	// Generate realistic file data for each episode
	for i := 0; i < anime.EpisodeCount; i++ {
		episodeNum := i + 1
		files[i] = domain.File{
			ID:       fmt.Sprintf("%s-ep%d", anime.ID, episodeNum),
			Name:     fmt.Sprintf("%s - Episode %d", anime.Title, episodeNum),
			Path:     fmt.Sprintf("/downloads/%s/Episode_%02d.mkv", anime.Title, episodeNum),
			Size:     int64(350000000 + i*10000000), // ~350-400 MB per episode
			Duration: 1400 + i*60,                    // ~23+ minutes per episode
			Type:     "video/x-matroska",             // MKV is standard
		}
	}

	return &domain.FileList{
		AnimeID: anime.ID,
		Files:   files,
	}
}

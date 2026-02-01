package app

import (
	"context"
	"log/slog"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// SearchService provides anime search across multiple sources.
type SearchService struct {
	resolvers []domain.IResolver
	logger    *slog.Logger
}

// NewSearchService creates a new SearchService.
func NewSearchService(resolvers []domain.IResolver, logger *slog.Logger) *SearchService {
	return &SearchService{
		resolvers: resolvers,
		logger:    logger,
	}
}

// SearchAnime searches for anime using registered resolvers with fallback.
func (s *SearchService) SearchAnime(ctx context.Context, query string) (*domain.SearchResult, error) {
	for _, resolver := range s.resolvers {
		result, err := resolver.Resolve(ctx, query)
		if err != nil {
			s.logger.Warn("resolver failed", slog.String("resolver", resolver.Name()), slog.String("error", err.Error()))
			continue
		}
		if result != nil {
			s.logger.Info("search succeeded", slog.String("resolver", resolver.Name()), slog.String("anime_id", result.AnimeID))
			return result, nil
		}
	}
	return nil, domain.NewAppError(domain.ErrSearchFailed, "could not find anime in any source")
}

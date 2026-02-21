package app

import (
	"context"
	"fmt"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// MockAnimeDetailService provides hardcoded anime details for testing.
// Real scraping service will be implemented in later stories.
type MockAnimeDetailService struct {
	fixtures map[string]domain.AnimeDetail
}

// NewMockAnimeDetailService creates a new mock service with test fixtures.
func NewMockAnimeDetailService() *MockAnimeDetailService {
	return &MockAnimeDetailService{
		fixtures: createTestFixtures(),
	}
}

// GetDetail returns anime details by ID from fixtures.
func (s *MockAnimeDetailService) GetDetail(ctx context.Context, id string) (domain.AnimeDetail, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return domain.AnimeDetail{}, ctx.Err()
	default:
	}

	detail, ok := s.fixtures[id]
	if !ok {
		return domain.AnimeDetail{}, fmt.Errorf("anime not found: %s", id)
	}

	return detail, nil
}

// createTestFixtures returns hardcoded test data for development.
// IDs MUST match devCatalogue() in cmd/asd-server/main.go so that
// search → detail navigation works end-to-end.
func createTestFixtures() map[string]domain.AnimeDetail {
	placeholder := "/assets/cover-placeholder.svg"

	return map[string]domain.AnimeDetail{
		"mushishi": {
			ID:           "mushishi",
			Title:        "Mushishi",
			ThumbnailURL: placeholder,
			Synopsis:     "Ginko, un mushishi — un expert des créatures primitives appelées mushi — voyage à travers un Japon rural et mystique pour aider ceux qui souffrent de phénomènes surnaturels liés aux mushi.",
			Year:         2005,
			Status:       "completed",
			Genres:       []string{"Drama", "Mystery", "Supernatural"},
			EpisodeCount: 26,
			Seasons: []domain.Season{
				{
					Number: 1,
					Name:   "Saison 1",
					Episodes: []domain.Episode{
						{Number: 1, Title: "La Vallée Verte", SeasonNumber: 1, URL: "https://example.com/mushishi/s1e1"},
						{Number: 2, Title: "La Lumière de la Paupière", SeasonNumber: 1, URL: "https://example.com/mushishi/s1e2"},
						{Number: 3, Title: "Tendre Cornes", SeasonNumber: 1, URL: "https://example.com/mushishi/s1e3"},
						{Number: 4, Title: "L'Oreiller d'Herbe", SeasonNumber: 1, URL: "https://example.com/mushishi/s1e4"},
						{Number: 5, Title: "Le Voyageur dans la Mer de l'Écriture", SeasonNumber: 1, URL: "https://example.com/mushishi/s1e5"},
					},
				},
			},
		},
		"mononoke": {
			ID:           "mononoke",
			Title:        "Mononoke",
			ThumbnailURL: placeholder,
			Synopsis:     "Un mystérieux vendeur de médicaments voyage à travers le Japon féodal, traquant et exorcisant les mononoke — des esprits maléfiques — en découvrant leur Forme, leur Vérité et leur Raison d'être.",
			Year:         2007,
			Status:       "completed",
			Genres:       []string{"Horror", "Mystery"},
			EpisodeCount: 12,
			Seasons: []domain.Season{
				{
					Number: 1,
					Name:   "Saison 1",
					Episodes: []domain.Episode{
						{Number: 1, Title: "Zashiki Warashi — Partie 1", SeasonNumber: 1, URL: "https://example.com/mononoke/s1e1"},
						{Number: 2, Title: "Zashiki Warashi — Partie 2", SeasonNumber: 1, URL: "https://example.com/mononoke/s1e2"},
						{Number: 3, Title: "Umi Bozu — Partie 1", SeasonNumber: 1, URL: "https://example.com/mononoke/s1e3"},
					},
				},
			},
		},
		"natsume-yuujinchou": {
			ID:           "natsume-yuujinchou",
			Title:        "Natsume Yuujinchou",
			ThumbnailURL: placeholder,
			Synopsis:     "Natsume Takashi hérite du Livre des Amis de sa grand-mère — un carnet contenant les noms d'esprits yokai qu'elle avait vaincus. Accompagné de Nyanko-sensei, il entreprend de libérer les yokai de leur contrat.",
			Year:         2008,
			Status:       "ongoing",
			Genres:       []string{"Slice of Life", "Supernatural"},
			EpisodeCount: 13,
			Seasons: []domain.Season{
				{
					Number: 1,
					Name:   "Saison 1",
					Episodes: []domain.Episode{
						{Number: 1, Title: "Le Chat et le Livre des Amis", SeasonNumber: 1, URL: "https://example.com/natsume/s1e1"},
						{Number: 2, Title: "Rosée au Sanctuaire", SeasonNumber: 1, URL: "https://example.com/natsume/s1e2"},
						{Number: 3, Title: "Le Yokai brûlé", SeasonNumber: 1, URL: "https://example.com/natsume/s1e3"},
						{Number: 4, Title: "Le Petit Renard", SeasonNumber: 1, URL: "https://example.com/natsume/s1e4"},
					},
				},
			},
		},
		"samurai-champloo": {
			ID:           "samurai-champloo",
			Title:        "Samurai Champloo",
			ThumbnailURL: placeholder,
			Synopsis:     "Fuu, une serveuse, recrute deux samouraïs aux styles opposés — le turbulent Mugen et le stoïque Jin — pour l'aider à trouver « le samouraï qui sent le tournesol ». Un road-trip dans un Japon Edo hip-hop.",
			Year:         2004,
			Status:       "completed",
			Genres:       []string{"Action", "Adventure"},
			EpisodeCount: 26,
			Seasons: []domain.Season{
				{
					Number: 1,
					Name:   "Saison 1",
					Episodes: []domain.Episode{
						{Number: 1, Title: "Tempestuous Temperaments", SeasonNumber: 1, URL: "https://example.com/samurai-champloo/s1e1"},
						{Number: 2, Title: "Redeye Reprisal", SeasonNumber: 1, URL: "https://example.com/samurai-champloo/s1e2"},
						{Number: 3, Title: "Hellhounds for Hire (Part 1)", SeasonNumber: 1, URL: "https://example.com/samurai-champloo/s1e3"},
						{Number: 4, Title: "Hellhounds for Hire (Part 2)", SeasonNumber: 1, URL: "https://example.com/samurai-champloo/s1e4"},
						{Number: 5, Title: "Artistic Anarchy", SeasonNumber: 1, URL: "https://example.com/samurai-champloo/s1e5"},
					},
				},
			},
		},
		"dororo": {
			ID:           "dororo",
			Title:        "Dororo",
			ThumbnailURL: placeholder,
			Synopsis:     "Hyakkimaru est né sans membres, peau ni organes — offerts aux démons par son père seigneur de guerre. Il parcourt le Japon féodal pour vaincre ces démons et reconquérir son corps, accompagné du jeune voleur Dororo.",
			Year:         2019,
			Status:       "completed",
			Genres:       []string{"Action", "Drama"},
			EpisodeCount: 24,
			Seasons: []domain.Season{
				{
					Number: 1,
					Name:   "Saison 1",
					Episodes: []domain.Episode{
						{Number: 1, Title: "L'histoire de Daigo", SeasonNumber: 1, URL: "https://example.com/dororo/s1e1"},
						{Number: 2, Title: "L'histoire de Bandai", SeasonNumber: 1, URL: "https://example.com/dororo/s1e2"},
						{Number: 3, Title: "L'histoire de Jukai", SeasonNumber: 1, URL: "https://example.com/dororo/s1e3"},
						{Number: 4, Title: "L'histoire de l'Épée Maudite", SeasonNumber: 1, URL: "https://example.com/dororo/s1e4"},
					},
				},
			},
		},
		"spice-and-wolf": {
			ID:           "spice-and-wolf",
			Title:        "Spice and Wolf",
			ThumbnailURL: placeholder,
			Synopsis:     "Kraft Lawrence, un marchand itinérant, rencontre Holo — une divinité-louve incarnée sous forme humaine. Ensemble, ils voyagent de ville en ville, négociant et spéculant, tandis que Holo cherche à rejoindre sa terre natale du Nord.",
			Year:         2008,
			Status:       "completed",
			Genres:       []string{"Fantasy", "Romance"},
			EpisodeCount: 13,
			Seasons: []domain.Season{
				{
					Number: 1,
					Name:   "Saison 1",
					Episodes: []domain.Episode{
						{Number: 1, Title: "Le Loup et les Meilleures Épices", SeasonNumber: 1, URL: "https://example.com/spice-and-wolf/s1e1"},
						{Number: 2, Title: "Le Loup et le Voyage Lointain", SeasonNumber: 1, URL: "https://example.com/spice-and-wolf/s1e2"},
						{Number: 3, Title: "Le Loup et la Négociation Habile", SeasonNumber: 1, URL: "https://example.com/spice-and-wolf/s1e3"},
					},
				},
			},
		},
	}
}

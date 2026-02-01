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
func createTestFixtures() map[string]domain.AnimeDetail {
	return map[string]domain.AnimeDetail{
		"naruto": {
			ID:           "naruto",
			Title:        "Naruto",
			ThumbnailURL: "https://cdn.anime-sama.si/naruto.jpg",
			Synopsis:     "Naruto Uzumaki is a young ninja who seeks recognition from his peers and dreams of becoming the Hokage, the leader of his village. The story follows his journey as he faces various challenges and makes friends along the way.",
			Year:         2002,
			Status:       "completed",
			Genres:       []string{"Action", "Adventure", "Shonen"},
			EpisodeCount: 220,
			Seasons: []domain.Season{
				{
					Number: 1,
					Name:   "Season 1",
					Episodes: []domain.Episode{
						{Number: 1, Title: "Enter: Naruto Uzumaki!", SeasonNumber: 1, URL: "https://anime-sama.si/naruto/s1e1"},
						{Number: 2, Title: "My Name is Konohamaru!", SeasonNumber: 1, URL: "https://anime-sama.si/naruto/s1e2"},
						{Number: 3, Title: "Sasuke and Sakura: Friends or Foes?", SeasonNumber: 1, URL: "https://anime-sama.si/naruto/s1e3"},
						{Number: 4, Title: "Pass or Fail: Survival Test", SeasonNumber: 1, URL: "https://anime-sama.si/naruto/s1e4"},
						{Number: 5, Title: "You Failed! Kakashi's Final Decision", SeasonNumber: 1, URL: "https://anime-sama.si/naruto/s1e5"},
						{Number: 6, Title: "A Dangerous Mission! Journey to the Land of Waves!", SeasonNumber: 1, URL: "https://anime-sama.si/naruto/s1e6"},
						{Number: 7, Title: "The Assassin of the Mist!", SeasonNumber: 1, URL: "https://anime-sama.si/naruto/s1e7"},
						{Number: 8, Title: "The Oath of Pain", SeasonNumber: 1, URL: "https://anime-sama.si/naruto/s1e8"},
					},
				},
			},
		},
		"naruto-shippuden": {
			ID:           "naruto-shippuden",
			Title:        "Naruto Shippuden",
			ThumbnailURL: "https://cdn.anime-sama.si/naruto-shippuden.jpg",
			Synopsis:     "Naruto Uzumaki returns after two years of training with a new goal: to save his friend Sasuke Uchiha from the evil Orochimaru. With the Akatsuki organization hunting for the Nine-Tails sealed within him, Naruto must become stronger while protecting those he cares about.",
			Year:         2007,
			Status:       "completed",
			Genres:       []string{"Action", "Adventure", "Shonen"},
			EpisodeCount: 500,
			Seasons: []domain.Season{
				{
					Number: 1,
					Name:   "Season 1",
					Episodes: []domain.Episode{
						{Number: 1, Title: "Homecoming", SeasonNumber: 1, URL: "https://anime-sama.si/naruto-shippuden/s1e1"},
						{Number: 2, Title: "The Akatsuki Makes Its Move", SeasonNumber: 1, URL: "https://anime-sama.si/naruto-shippuden/s1e2"},
						{Number: 3, Title: "The Results of Training", SeasonNumber: 1, URL: "https://anime-sama.si/naruto-shippuden/s1e3"},
						{Number: 4, Title: "The Jinchuriki of the Sand", SeasonNumber: 1, URL: "https://anime-sama.si/naruto-shippuden/s1e4"},
					},
				},
				{
					Number: 2,
					Name:   "Season 2",
					Episodes: []domain.Episode{
						{Number: 1, Title: "A New Target", SeasonNumber: 2, URL: "https://anime-sama.si/naruto-shippuden/s2e1"},
						{Number: 2, Title: "The Mysterious Mission", SeasonNumber: 2, URL: "https://anime-sama.si/naruto-shippuden/s2e2"},
					},
				},
			},
		},
		"one-piece": {
			ID:           "one-piece",
			Title:        "One Piece",
			ThumbnailURL: "https://cdn.anime-sama.si/one-piece.jpg",
			Synopsis:     "Monkey D. Luffy sets off on an adventure to find the legendary treasure known as One Piece and become the Pirate King. Along the way, he gathers a diverse crew of pirates, each with their own dreams and abilities.",
			Year:         1999,
			Status:       "ongoing",
			Genres:       []string{"Action", "Adventure", "Comedy", "Shonen"},
			EpisodeCount: 1100,
			Seasons: []domain.Season{
				{
					Number: 1,
					Name:   "East Blue Saga",
					Episodes: []domain.Episode{
						{Number: 1, Title: "I'm Luffy! The Man Who Will Become Pirate King!", SeasonNumber: 1, URL: "https://anime-sama.si/one-piece/s1e1"},
						{Number: 2, Title: "The Great Swordsman Appears! Pirate Hunter, Roronoa Zoro!", SeasonNumber: 1, URL: "https://anime-sama.si/one-piece/s1e2"},
						{Number: 3, Title: "Morgan vs. Luffy! Who's This Beautiful Young Girl?", SeasonNumber: 1, URL: "https://anime-sama.si/one-piece/s1e3"},
					},
				},
			},
		},
		"attack-on-titan": {
			ID:           "attack-on-titan",
			Title:        "Attack on Titan",
			ThumbnailURL: "https://cdn.anime-sama.si/attack-on-titan.jpg",
			Synopsis:     "Centuries ago, mankind was slaughtered to near extinction by monstrous humanoid creatures called Titans. Now, the remaining humans live within three concentric walls protecting them from the Titans. Eren Yeager vows to eliminate all Titans after witnessing his mother's death.",
			Year:         2013,
			Status:       "completed",
			Genres:       []string{"Action", "Dark Fantasy", "Drama"},
			EpisodeCount: 87,
			Seasons: []domain.Season{
				{
					Number: 1,
					Name:   "Season 1",
					Episodes: []domain.Episode{
						{Number: 1, Title: "To You, in 2000 Years", SeasonNumber: 1, URL: "https://anime-sama.si/aot/s1e1"},
						{Number: 2, Title: "That Day", SeasonNumber: 1, URL: "https://anime-sama.si/aot/s1e2"},
					},
				},
			},
		},
		"demon-slayer": {
			ID:           "demon-slayer",
			Title:        "Demon Slayer: Kimetsu no Yaiba",
			ThumbnailURL: "https://cdn.anime-sama.si/demon-slayer.jpg",
			Synopsis:     "Tanjiro Kamado becomes a demon slayer after his family is slaughtered and his younger sister Nezuko is turned into a demon. He embarks on a journey to find a cure for Nezuko and avenge his family.",
			Year:         2019,
			Status:       "ongoing",
			Genres:       []string{"Action", "Dark Fantasy", "Adventure"},
			EpisodeCount: 55,
			Seasons: []domain.Season{
				{
					Number: 1,
					Name:   "Season 1",
					Episodes: []domain.Episode{
						{Number: 1, Title: "Cruelty", SeasonNumber: 1, URL: "https://anime-sama.si/demon-slayer/s1e1"},
						{Number: 2, Title: "Trainer Sakonji Urokodaki", SeasonNumber: 1, URL: "https://anime-sama.si/demon-slayer/s1e2"},
					},
				},
			},
		},
	}
}

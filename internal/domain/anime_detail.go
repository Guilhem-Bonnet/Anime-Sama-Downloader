package domain

// AnimeDetail represents complete information about an anime.
type AnimeDetail struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Synopsis     string   `json:"synopsis"`
	Year         int      `json:"year"`
	Status       string   `json:"status"`       // "ongoing", "completed", "planning"
	Genres       []string `json:"genres"`       // ["Action", "Adventure", "Shonen"]
	EpisodeCount int      `json:"episode_count"`
	Seasons      []Season `json:"seasons"`
}

// Season represents a season of an anime.
type Season struct {
	Number   int       `json:"number"`   // 1, 2, 3...
	Name     string    `json:"name"`     // "Season 1", "Part 2", etc.
	Episodes []Episode `json:"episodes"`
}

// Episode represents a single episode within a season.
type Episode struct {
	Number       int    `json:"number"`        // Episode number within season
	Title        string `json:"title"`         // Episode title (optional)
	SeasonNumber int    `json:"season_number"` // Which season this belongs to
	URL          string `json:"url"`           // Download URL (placeholder for now)
}

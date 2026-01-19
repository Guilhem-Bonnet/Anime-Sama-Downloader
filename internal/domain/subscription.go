package domain

import "time"

type Subscription struct {
	ID string

	// BaseURL pointe vers une saison/langue (ex: https://anime-sama.si/catalogue/xxx/saison1/vostfr/)
	BaseURL string

	// Label est un nom libre pour affichage/chemins.
	Label string

	// Player: "auto" ou un nom exact (ex: "Player 1").
	Player string

	LastScheduledEpisode  int
	LastDownloadedEpisode int
	LastAvailableEpisode  int

	NextCheckAt  time.Time
	LastCheckedAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

package domain

type OutputNamingMode string

const (
	OutputNamingLegacy      OutputNamingMode = "legacy"
	OutputNamingMediaServer OutputNamingMode = "media-server"
)

type Settings struct {
	// Chemin racine de destination.
	Destination string `json:"destination"`

	// Nommage des fichiers/chemins.
	OutputNamingMode OutputNamingMode `json:"outputNamingMode"`
	SeparateLang     bool             `json:"separateLang"`

	// Concurrence (sera utilisée plus tard pour les workers).
	MaxWorkers             int `json:"maxWorkers"`
	MaxConcurrentDownloads int `json:"maxConcurrentDownloads"`

	// Media servers.
	JellyfinURL    string `json:"jellyfinUrl"`
	JellyfinAPIKey string `json:"jellyfinApiKey"`
	PlexURL        string `json:"plexUrl"`
	PlexToken      string `json:"plexToken"`
	PlexSectionID  string `json:"plexSectionId"`
}

func DefaultSettings() Settings {
	return Settings{
		Destination:            "videos",
		OutputNamingMode:       OutputNamingLegacy,
		SeparateLang:           false,
		MaxWorkers:             2,
		MaxConcurrentDownloads: 4,
	}
}

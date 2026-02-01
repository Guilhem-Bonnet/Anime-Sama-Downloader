package domain

import "time"

// Download represents a media download in the system.
type Download struct {
	DownloadID  string    `json:"download_id"`
	JobID       string    `json:"job_id"`
	AnimeID     string    `json:"anime_id"`
	EpisodeNum  int       `json:"episode_number"`
	Metadata    string    `json:"metadata"` // JSON-encoded metadata
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IsCompleted returns true if download is finished.
func (d *Download) IsCompleted() bool {
	// Will check job status via repository
	return false
}

// GetFilePath returns the expected file path once download is complete.
func (d *Download) GetFilePath() string {
	// Will construct path from anime_id + episode_num
	return ""
}

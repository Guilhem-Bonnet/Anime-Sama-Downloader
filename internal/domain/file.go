package domain

// File represents a downloadable file (episode) with metadata
type File struct {
	ID       string // Unique identifier for the file
	Name     string // Display name (e.g., "Episode 1 - The Beginning")
	Path     string // Full path on storage
	Size     int64  // Size in bytes
	Duration int    // Duration in seconds (for video files)
	Type     string // File type (e.g., "video/mp4", "video/mkv")
}

// FileList represents a collection of files for an anime
type FileList struct {
	AnimeID string
	Files   []File
}

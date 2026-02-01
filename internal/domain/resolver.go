package domain

import "context"

// SearchResult represents anime search results from a resolver.
type SearchResult struct {
	AnimeID  string `json:"anime_id"`
	Title    string `json:"title"`
	Episodes int    `json:"episodes"`
	Source   string `json:"source"` // AnimeSama, MangaDex, etc.
}

// IResolver defines the interface for search source adapters.
type IResolver interface {
	// Resolve searches for anime in this source.
	Resolve(ctx context.Context, query string) (*SearchResult, error)
	// Name returns the resolver's display name.
	Name() string
}

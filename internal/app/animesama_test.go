package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestBestPlayer_PrefersSupportedHostsOverEmbed4Me(t *testing.T) {
	players := map[string][]string{
		"Player 1": {"https://lpayer.embed4me.com/#ia5qn", "https://lpayer.embed4me.com/#cufsv"},
		"Player 2": {"https://oneupload.net/abcd", "https://oneupload.net/efgh"},
	}

	got := BestPlayer(players)
	if got != "Player 2" {
		t.Fatalf("expected %q, got %q", "Player 2", got)
	}
}

func TestBestPlayer_PrefersDirectMP4(t *testing.T) {
	players := map[string][]string{
		"Player 1": {"https://example.com/stream.m3u8"},
		"Player 2": {"https://example.com/video.mp4"},
	}

	got := BestPlayer(players)
	if got != "Player 2" {
		t.Fatalf("expected %q, got %q", "Player 2", got)
	}
}

// TestParseEpisodesJS_Success tests parsing valid episodes JS
func TestParseEpisodesJS_Success(t *testing.T) {
	jsText := `var eps1 = ['https://example.com/ep1', 'https://example.com/ep2', 'https://example.com/ep3'];
var eps2 = ['https://example.com/alt1', 'https://example.com/alt2'];`

	result, err := ParseEpisodesJS(jsText)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(result.Players))
	}

	if player1, ok := result.Players["Player 1"]; ok {
		if len(player1) != 3 {
			t.Fatalf("expected 3 episodes for Player 1, got %d", len(player1))
		}
		if player1[0] != "https://example.com/ep1" {
			t.Fatalf("unexpected episodes for Player 1: %v", player1)
		}
	} else {
		t.Fatalf("Player 1 not found in result")
	}

	if player2, ok := result.Players["Player 2"]; ok {
		if len(player2) != 2 {
			t.Fatalf("expected 2 episodes for Player 2, got %d", len(player2))
		}
	} else {
		t.Fatalf("Player 2 not found in result")
	}
}

// TestParseEpisodesJS_NoMatches tests parsing with no episode arrays
func TestParseEpisodesJS_NoMatches(t *testing.T) {
	jsText := `// No episodes here
var someVar = 123;`

	_, err := ParseEpisodesJS(jsText)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "no episodes arrays found" {
		t.Fatalf("expected 'no episodes arrays found', got %q", err.Error())
	}
}

// TestParseEpisodesJS_EmptyEpisodeList tests parsing returns error when no plausible URLs
func TestParseEpisodesJS_EmptyEpisodeList(t *testing.T) {
	jsText := `var eps1 = [];`

	_, err := ParseEpisodesJS(jsText)
	if err == nil {
		t.Fatalf("expected error for no plausible URLs, got nil")
	}
}

// TestParseEpisodesJS_SingleElement tests parsing single episode
func TestParseEpisodesJS_SingleElement(t *testing.T) {
	jsText := `var eps1 = ['https://example.com/single'];`

	result, err := ParseEpisodesJS(jsText)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if player1, ok := result.Players["Player 1"]; ok {
		if len(player1) != 1 || player1[0] != "https://example.com/single" {
			t.Fatalf("expected [https://example.com/single], got %v", player1)
		}
	} else {
		t.Fatalf("Player 1 not found")
	}
}

// TestParseEpisodesJS_WithSpaces tests parsing with various whitespace
func TestParseEpisodesJS_WithSpaces(t *testing.T) {
	jsText := `var   eps1   =   [  'https://example.com/ep1'  ,  'https://example.com/ep2'  ];`

	result, err := ParseEpisodesJS(jsText)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if player1, ok := result.Players["Player 1"]; ok {
		if len(player1) != 2 {
			t.Fatalf("expected 2 episodes, got %d", len(player1))
		}
	} else {
		t.Fatalf("Player 1 not found")
	}
}

// TestParseEpisodesJS_HighPlayerNumbers tests parsing with higher player numbers
func TestParseEpisodesJS_HighPlayerNumbers(t *testing.T) {
	jsText := `var eps10 = ['https://example.com/p10ep1'];
var eps25 = ['https://example.com/p25ep1'];`

	result, err := ParseEpisodesJS(jsText)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(result.Players))
	}

	if player10, ok := result.Players["Player 10"]; ok {
		if len(player10) != 1 {
			t.Fatalf("expected 1 episode for Player 10, got %d", len(player10))
		}
	} else {
		t.Fatalf("Player 10 not found")
	}
}

// TestFetchEpisodesJS_Success tests successful fetch with valid server
func TestFetchEpisodesJS_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/episodes.js" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/javascript")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`var eps1 = ['test.html'];`))
	}))
	defer srv.Close()

	result, err := FetchEpisodesJS(context.Background(), srv.URL+"/")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(result, "var eps1") {
		t.Fatalf("expected episodes.js content, got %q", result)
	}
}

// TestFetchEpisodesJS_NotFound tests fetch with 404 response
func TestFetchEpisodesJS_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	_, err := FetchEpisodesJS(context.Background(), srv.URL+"/")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "episodes.js http error") {
		t.Fatalf("expected http error message, got %q", err.Error())
	}
}

// TestFetchEpisodesJS_InvalidBaseURL tests fetch with invalid base URL
func TestFetchEpisodesJS_InvalidBaseURL(t *testing.T) {
	_, err := FetchEpisodesJS(context.Background(), "not a valid url")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// TestFetchEpisodesJS_ContextCanceled tests fetch respects context cancellation
func TestFetchEpisodesJS_ContextCanceled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow server
		<-r.Context().Done()
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := FetchEpisodesJS(ctx, srv.URL+"/")
	if err == nil {
		t.Fatalf("expected error for canceled context, got nil")
	}
}

// TestCanonicalizeAnimeSamaBaseURL_Success tests valid URLs are canonicalized
func TestCanonicalizeAnimeSamaBaseURL_Success(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://anime-sama.si/", "https://anime-sama.si/"},
		{"https://anime-sama.si", "https://anime-sama.si/"},
		{"https://www.anime-sama.si/", "https://anime-sama.si/"},
		{"https://anime-sama.tv/", "https://anime-sama.si/"},
		{"https://anime-sama.fr/", "https://anime-sama.si/"},
		{"https://anime-sama.org/", "https://anime-sama.si/"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := CanonicalizeAnimeSamaBaseURL(tt.input)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if result != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestCanonicalizeAnimeSamaBaseURL_AddTrailingSlash tests trailing slash is added
func TestCanonicalizeAnimeSamaBaseURL_AddTrailingSlash(t *testing.T) {
	result, err := CanonicalizeAnimeSamaBaseURL("https://anime-sama.si/anime/test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.HasSuffix(result, "/") {
		t.Fatalf("expected trailing slash, got %q", result)
	}
}

// TestCanonicalizeAnimeSamaBaseURL_InvalidURL tests invalid URLs are rejected
func TestCanonicalizeAnimeSamaBaseURL_InvalidURL(t *testing.T) {
	tests := []string{
		"not a url",
		"",
		"    ",
		"://invalid",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := CanonicalizeAnimeSamaBaseURL(input)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", input)
			}
		})
	}
}

// TestCanonicalizeAnimeSamaBaseURL_PreservesPath tests path is preserved
func TestCanonicalizeAnimeSamaBaseURL_PreservesPath(t *testing.T) {
	result, err := CanonicalizeAnimeSamaBaseURL("https://anime-sama.si/path/to/anime")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(result, "/path/to/anime") {
		t.Fatalf("expected path to be preserved, got %q", result)
	}
}

// TestCanonicalizeAnimeSamaBaseURL_WithPort tests URL with port
func TestCanonicalizeAnimeSamaBaseURL_WithPort(t *testing.T) {
	result, err := CanonicalizeAnimeSamaBaseURL("https://anime-sama.si:8443/")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(result, ":8443") {
		t.Fatalf("expected port to be preserved, got %q", result)
	}
}

// TestCanonicalizeAnimeSamaBaseURL_WithWhitespace tests whitespace is trimmed
func TestCanonicalizeAnimeSamaBaseURL_WithWhitespace(t *testing.T) {
	result, err := CanonicalizeAnimeSamaBaseURL("  https://anime-sama.si/  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(result, "anime-sama.si") {
		t.Fatalf("expected valid URL after trim, got %q", result)
	}
}

// TestCanonicalizeAnimeSamaBaseURL_CaseInsensitive tests case-insensitive domain matching
func TestCanonicalizeAnimeSamaBaseURL_CaseInsensitive(t *testing.T) {
	result, err := CanonicalizeAnimeSamaBaseURL("https://ANIME-SAMA.TV/")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(result, "anime-sama.si") {
		t.Fatalf("expected lowercase anime-sama.si, got %q", result)
	}
}

// TestParseEpisodesJS_ManyPlayers tests parsing many players
func TestParseEpisodesJS_ManyPlayers(t *testing.T) {
	jsText := `var eps1 = ['https://example.com/1_1', 'https://example.com/1_2'];
var eps2 = ['https://example.com/2_1', 'https://example.com/2_2'];
var eps3 = ['https://example.com/3_1', 'https://example.com/3_2'];
var eps4 = ['https://example.com/4_1', 'https://example.com/4_2'];
var eps5 = ['https://example.com/5_1', 'https://example.com/5_2'];`

	result, err := ParseEpisodesJS(jsText)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Players) != 5 {
		t.Fatalf("expected 5 players, got %d", len(result.Players))
	}

	for i := 1; i <= 5; i++ {
		playerName := "Player " + strconv.Itoa(i)
		if _, ok := result.Players[playerName]; !ok {
			t.Fatalf("Player %d not found", i)
		}
	}
}

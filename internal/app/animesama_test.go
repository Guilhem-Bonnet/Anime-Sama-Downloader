package app

import "testing"

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

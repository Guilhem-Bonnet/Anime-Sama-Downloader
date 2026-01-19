package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResolveDirectMediaURL_PrefersMP4FromHTML(t *testing.T) {
	var videoURL string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/page":
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte("<html><body><source src=\"" + videoURL + "\" type=\"video/mp4\"></body></html>"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()
	videoURL = ts.URL + "/video.mp4"

	got, err := ResolveDirectMediaURL(context.Background(), ts.URL+"/page")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got != videoURL {
		t.Fatalf("expected %q, got %q", videoURL, got)
	}
}

func TestResolveDirectMediaURL_FollowsIframeAndResolvesRelative(t *testing.T) {
	var videoURL string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/outer":
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte("<html><body><iframe src=\"/inner\"></iframe></body></html>"))
		case "/inner":
			w.Header().Set("Content-Type", "text/html")
			// relative media url
			_, _ = w.Write([]byte("<html><body><source src=\"/video.mp4\" type=\"video/mp4\"></body></html>"))
		case "/video.mp4":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()
	videoURL = ts.URL + "/video.mp4"

	got, err := ResolveDirectMediaURL(context.Background(), ts.URL+"/outer")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got != videoURL {
		t.Fatalf("expected %q, got %q", videoURL, got)
	}
}

func TestResolveDirectMediaURL_DecodesEscapedSlashes(t *testing.T) {
	var escaped string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html><script>var u=\"" + escaped + "\";</script></html>"))
	}))
	defer ts.Close()

	// build something like https:\/\/host\/video.m3u8
	trimmed := strings.TrimPrefix(ts.URL, "http://")
	escaped = "http:\\/\\/" + trimmed + "\\/video.m3u8"

	got, err := ResolveDirectMediaURL(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got != "http://"+trimmed+"/video.m3u8" {
		t.Fatalf("expected %q, got %q", "http://"+trimmed+"/video.m3u8", got)
	}
}

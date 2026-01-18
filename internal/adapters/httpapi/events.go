package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Petit heartbeat tant qu'on n'a pas encore de bus.
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	fmt.Fprintf(w, "event: hello\ndata: {\"status\":\"connected\"}\n\n")
	flusher.Flush()

	var sub <-chan ports.Event
	var cancel func()
	if s.bus != nil {
		sub, cancel = s.bus.Subscribe()
		defer cancel()
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case evt := <-sub:
			payload := evt.Payload
			if !json.Valid(payload) {
				payload = []byte(`{"error":"invalid payload"}`)
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", evt.Topic, payload)
			flusher.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, "event: ping\ndata: {}\n\n")
			flusher.Flush()
		}
	}
}

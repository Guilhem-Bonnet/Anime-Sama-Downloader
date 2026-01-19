package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AnimeSamaEpisodes struct {
	// Players maps "Player 1" -> episode URLs slice (index 0 => ep1). Empty string means missing/unavailable.
	Players map[string][]string
}

func CanonicalizeAnimeSamaBaseURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid baseUrl")
	}
	// normalize known domains to anime-sama.si
	host := strings.ToLower(u.Host)
	host = strings.TrimPrefix(host, "www.")
	if host == "anime-sama.tv" || host == "anime-sama.fr" || host == "anime-sama.org" || host == "anime-sama.si" {
		u.Host = "anime-sama.si"
	}
	// ensure trailing slash
	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}
	return u.String(), nil
}

func FetchEpisodesJS(ctx context.Context, baseURL string) (string, error) {
	canon, err := CanonicalizeAnimeSamaBaseURL(baseURL)
	if err != nil {
		return "", err
	}
	jsURL := strings.TrimRight(canon, "/") + "/episodes.js"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jsURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "asd-server")
	req.Header.Set("Accept", "text/javascript,*/*;q=0.1")
	req.Header.Set("Referer", canon)

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("episodes.js http error: %s", resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

var reEpisodesArray = regexp.MustCompile(`var\s+(eps\d+)\s*=\s*\[([^\]]*)\];`)
var reQuoted = regexp.MustCompile(`'([^']*)'`)

func ParseEpisodesJS(jsText string) (AnimeSamaEpisodes, error) {
	matches := reEpisodesArray.FindAllStringSubmatch(jsText, -1)
	if len(matches) == 0 {
		return AnimeSamaEpisodes{}, errors.New("no episodes arrays found")
	}

	players := map[string][]string{}
	for _, m := range matches {
		name := m[1] // eps1
		content := m[2]
		num := extractDigits(name)
		if num == "" {
			continue
		}
		playerName := "Player " + num

		items := reQuoted.FindAllStringSubmatch(content, -1)
		urls := make([]string, 0, len(items))
		for _, it := range items {
			u := strings.TrimSpace(it[1])
			if isPlausibleEpisodeURL(u) {
				urls = append(urls, u)
			} else {
				urls = append(urls, "")
			}
		}
		players[playerName] = urls
	}

	hasAny := false
	for _, urls := range players {
		for _, u := range urls {
			if u != "" {
				hasAny = true
				break
			}
		}
		if hasAny {
			break
		}
	}
	if !hasAny {
		return AnimeSamaEpisodes{}, errors.New("no plausible episode urls")
	}

	return AnimeSamaEpisodes{Players: players}, nil
}

func extractDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func isPlausibleEpisodeURL(u string) bool {
	u = strings.TrimSpace(u)
	if u == "" {
		return false
	}
	if !(strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")) {
		return false
	}
	if strings.HasSuffix(u, "=") {
		return false
	}
	if strings.HasSuffix(u, "/embed/") {
		return false
	}

	lu := strings.ToLower(u)
	// VK placeholder often has empty oid/id
	if strings.Contains(lu, "vk.com/video_ext.php") {
		parsed, err := url.Parse(u)
		if err != nil {
			return false
		}
		q := parsed.Query()
		oid := q.Get("oid")
		id := q.Get("id")
		if oid == "" || id == "" {
			return false
		}
		if _, err := strconv.Atoi(strings.TrimPrefix(oid, "-")); err != nil {
			return false
		}
		if _, err := strconv.Atoi(id); err != nil {
			return false
		}
	}

	// Sibnet placeholder: missing videoid
	if strings.Contains(lu, "video.sibnet.ru") && strings.Contains(lu, "shell.php") {
		parsed, err := url.Parse(u)
		if err != nil {
			return false
		}
		videoid := parsed.Query().Get("videoid")
		if videoid == "" {
			return false
		}
		if _, err := strconv.Atoi(videoid); err != nil {
			return false
		}
	}

	// Vidmoly placeholder: embed-.html without id
	if strings.Contains(lu, "vidmoly") && strings.Contains(lu, "/embed-") {
		parsed, err := url.Parse(u)
		if err != nil {
			return false
		}
		p := parsed.Path
		if !strings.HasPrefix(p, "/embed-") || !strings.HasSuffix(p, ".html") {
			return false
		}
		mid := strings.TrimSuffix(strings.TrimPrefix(p, "/embed-"), ".html")
		if strings.TrimSpace(mid) == "" {
			return false
		}
	}

	// SendVid placeholder: /embed/<id> with empty id
	if strings.Contains(lu, "sendvid.com") && strings.Contains(lu, "/embed/") {
		parsed, err := url.Parse(u)
		if err != nil {
			return false
		}
		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		if len(parts) < 2 || parts[0] != "embed" || strings.TrimSpace(parts[1]) == "" {
			return false
		}
	}

	return true
}

func BestPlayer(players map[string][]string) string {
	best := ""
	bestCount := -1
	for name, urls := range players {
		c := 0
		for _, u := range urls {
			if u != "" {
				c++
			}
		}
		if c > bestCount {
			bestCount = c
			best = name
		}
	}
	if best == "" {
		return "auto"
	}
	return best
}

func MaxAvailableEpisode(urls []string) int {
	max := 0
	for i := len(urls) - 1; i >= 0; i-- {
		if urls[i] != "" {
			max = i + 1
			break
		}
	}
	return max
}

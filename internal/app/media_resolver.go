package app

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	reMediaAbs = regexp.MustCompile(`(?i)(https?:)?//[^\s"']+\.(mp4|m3u8)(\?[^\s"']*)?`)
	reMediaRel = regexp.MustCompile(`(?i)(/[^\s"']+\.(mp4|m3u8)(\?[^\s"']*)?)`)

	reIframeSrc    = regexp.MustCompile(`(?i)<iframe[^>]+src=['"]([^'"]+)['"]`)
	reMetaRefresh  = regexp.MustCompile(`(?i)<meta[^>]+http-equiv=['"]refresh['"][^>]+content=['"][^'"]*url=([^'">\s]+)[^'"]*['"]`)
	reLocationHref = regexp.MustCompile(`(?i)location\.(href|replace)\s*\(\s*['"]([^'"]+)['"]\s*\)`)
)

func looksLikeDirectMediaURL(u string) bool {
	lu := strings.ToLower(strings.TrimSpace(u))
	return strings.Contains(lu, ".mp4") || strings.Contains(lu, ".m3u8")
}

func normalizeTextForURLScan(s string) string {
	// Common JS/JSON escapes seen in embed pages.
	s = strings.ReplaceAll(s, `\/`, "/")
	s = strings.ReplaceAll(s, `\u0026`, "&")
	s = strings.ReplaceAll(s, `\u002F`, "/")
	s = strings.ReplaceAll(s, `\u003A`, ":")
	// Common HTML entity.
	s = strings.ReplaceAll(s, "&amp;", "&")
	return s
}

func resolveRef(base *url.URL, ref string) (string, bool) {
	ref = strings.TrimSpace(ref)
	if ref == "" || base == nil {
		return "", false
	}
	if strings.HasPrefix(ref, "//") {
		return base.Scheme + ":" + ref, true
	}
	uu, err := url.Parse(ref)
	if err != nil {
		return "", false
	}
	return base.ResolveReference(uu).String(), true
}

func pickDirectMediaCandidate(base *url.URL, text string) (string, bool) {
	text = normalizeTextForURLScan(text)

	matches := reMediaAbs.FindAllString(text, -1)
	if len(matches) > 0 {
		fix := func(m string) string {
			m = strings.TrimSpace(m)
			if strings.HasPrefix(m, "//") {
				return base.Scheme + ":" + m
			}
			return m
		}
		for _, m := range matches {
			m = fix(m)
			if strings.Contains(strings.ToLower(m), ".mp4") {
				return m, true
			}
		}
		return fix(matches[0]), true
	}

	// Relative URLs like /video.mp4
	matches = reMediaRel.FindAllString(text, -1)
	if len(matches) > 0 {
		best := ""
		for _, m := range matches {
			abs, ok := resolveRef(base, m)
			if !ok {
				continue
			}
			if strings.Contains(strings.ToLower(abs), ".mp4") {
				return abs, true
			}
			if best == "" {
				best = abs
			}
		}
		if best != "" {
			return best, true
		}
	}

	return "", false
}

func resolveDirectMediaURL(ctx context.Context, raw string, depth int) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("empty url")
	}
	if looksLikeDirectMediaURL(raw) {
		return raw, nil
	}
	if depth <= 0 {
		return raw, nil
	}

	base, err := url.Parse(raw)
	if err != nil {
		return raw, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return raw, nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) Gecko/20100101 Firefox/120.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return raw, nil
	}
	defer resp.Body.Close()

	ct := strings.ToLower(resp.Header.Get("Content-Type"))
	if !strings.Contains(ct, "text/html") && !strings.Contains(ct, "application/xhtml") {
		return raw, nil
	}

	b, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	text := string(b)

	if u, ok := pickDirectMediaCandidate(base, text); ok {
		return u, nil
	}

	// Follow common indirections (iframe / meta refresh / location.href)
	for _, re := range []*regexp.Regexp{reMetaRefresh, reLocationHref, reIframeSrc} {
		m := re.FindStringSubmatch(text)
		if len(m) >= 2 {
			ref := m[len(m)-1]
			next, ok := resolveRef(base, ref)
			if !ok {
				continue
			}
			res, _ := resolveDirectMediaURL(ctx, next, depth-1)
			if strings.TrimSpace(res) != "" && res != raw {
				return res, nil
			}
		}
	}

	return raw, nil
}

// ResolveDirectMediaURL tries to turn a host/embed/page URL into a direct media URL.
// It prefers mp4 over m3u8 when both are present.
func ResolveDirectMediaURL(ctx context.Context, raw string) (string, error) {
	return resolveDirectMediaURL(ctx, raw, 3)
}

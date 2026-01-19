package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/transform"
	"golang.org/x/text/runes"
	"unicode"
)

type AnimeSamaCandidate struct {
	CatalogueURL string  `json:"catalogueUrl"`
	Slug         string  `json:"slug"`
	MatchedTitle string  `json:"matchedTitle"`
	Score        float64 `json:"score"`
}

type AnimeSamaCatalogueResolver struct {
	BaseURL string
	Client  *http.Client
}

func NewAnimeSamaCatalogueResolver() *AnimeSamaCatalogueResolver {
	return &AnimeSamaCatalogueResolver{
		BaseURL: "https://anime-sama.si",
		Client: &http.Client{
			Timeout: 12 * time.Second,
		},
	}
}

func (r *AnimeSamaCatalogueResolver) WithBaseURL(base string) *AnimeSamaCatalogueResolver {
	if strings.TrimSpace(base) != "" {
		r.BaseURL = strings.TrimRight(strings.TrimSpace(base), "/")
	}
	return r
}

var reHyphens = regexp.MustCompile(`-+`)

var (
	reParensContent = regexp.MustCompile(`\([^\)]*\)`)
	reBrackets      = regexp.MustCompile(`\[[^\]]*\]`)
	reSaisonSuffix  = regexp.MustCompile(`(?i)\b(saison|season|cour|part|partie)\s*(\d+)\b`)
	reNthSeason     = regexp.MustCompile(`(?i)\b(\d+)(st|nd|rd|th)\s+season\b`)
)

func titleVariantsForProbe(title string) []string {
	t := strings.TrimSpace(title)
	if t == "" {
		return nil
	}
	variants := []string{t}
	// Drop bracketed content: "(TV)", "[Dub]", etc.
	clean := reParensContent.ReplaceAllString(t, " ")
	clean = reBrackets.ReplaceAllString(clean, " ")
	clean = strings.Join(strings.Fields(clean), " ")
	if clean != "" && clean != t {
		variants = append(variants, clean)
	}
	// Remove explicit season/cour/part suffix.
	clean2 := reSaisonSuffix.ReplaceAllString(clean, " ")
	clean2 = reNthSeason.ReplaceAllString(clean2, " ")
	clean2 = strings.Join(strings.Fields(clean2), " ")
	if clean2 != "" && clean2 != clean {
		variants = append(variants, clean2)
	}

	seen := map[string]struct{}{}
	out := make([]string, 0, len(variants))
	for _, v := range variants {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		k := strings.ToLower(v)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, v)
		if len(out) >= 6 {
			break
		}
	}
	return out
}

func slugifyAnimeSamaTitle(title string) string {
	s := strings.TrimSpace(strings.ToLower(title))
	if s == "" {
		return ""
	}

	// Remove accents (NFD -> remove Mn -> NFC).
	tr := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	if out, _, err := transform.String(tr, s); err == nil {
		s = out
	}

	// Replace common punctuation/separators with spaces.
	replacer := strings.NewReplacer("'", " ", "/", " ", ":", " ", "+", " ", "&", " ")
	s = replacer.Replace(s)

	// Keep [a-z0-9 -] only.
	b := strings.Builder{}
	b.Grow(len(s))
	for _, ch := range s {
		switch {
		case ch >= 'a' && ch <= 'z':
			b.WriteRune(ch)
		case ch >= '0' && ch <= '9':
			b.WriteRune(ch)
		case ch == ' ' || ch == '-':
			b.WriteRune(' ')
		default:
			// skip
		}
	}
	s = strings.TrimSpace(b.String())
	s = strings.ReplaceAll(s, " ", "-")
	s = reHyphens.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func slugVariants(slug string) []string {
	base := strings.Trim(strings.ToLower(strings.TrimSpace(slug)), "-")
	if base == "" {
		return nil
	}

	variants := []string{base}
	variants = append(variants, reHyphens.ReplaceAllString(base, "-"))

	// apostrophe-s like: hells-... <-> hell-s-...
	variants = append(variants, regexp.MustCompile(`([a-z])s-([a-z])`).ReplaceAllString(base, `$1-s-$2`))
	variants = append(variants, regexp.MustCompile(`([a-z])-s-([a-z])`).ReplaceAllString(base, `$1s-$2`))

	seen := map[string]struct{}{}
	out := make([]string, 0, len(variants))
	for _, v := range variants {
		v = strings.Trim(v, "-")
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
		if len(out) >= 6 {
			break
		}
	}
	return out
}

func (r *AnimeSamaCatalogueResolver) ProbeSlug(ctx context.Context, slug string) (string, error) {
	if r == nil {
		return "", errors.New("nil resolver")
	}
	slug = strings.Trim(strings.TrimSpace(slug), "/")
	if slug == "" {
		return "", errors.New("empty slug")
	}
	base := strings.TrimRight(strings.TrimSpace(r.BaseURL), "/")
	url := fmt.Sprintf("%s/catalogue/%s/", base, slug)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) Gecko/20100101 Firefox/120.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3")
	client := r.Client
	if client == nil {
		client = &http.Client{Timeout: 12 * time.Second}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("not found")
	}
	return url, nil
}

func (r *AnimeSamaCatalogueResolver) ResolveCandidates(ctx context.Context, titles []string, maxCandidates int) ([]AnimeSamaCandidate, error) {
	if maxCandidates <= 0 {
		maxCandidates = 3
	}

	out := []AnimeSamaCandidate{}
	seenSlug := map[string]struct{}{}

	for ti, raw := range titles {
		for tv, t := range titleVariantsForProbe(raw) {
			slug := slugifyAnimeSamaTitle(t)
			for vi, v := range slugVariants(slug) {
			if _, ok := seenSlug[v]; ok {
				continue
			}
			seenSlug[v] = struct{}{}

			url, err := r.ProbeSlug(ctx, v)
			if err != nil {
				continue
			}

				score := 0.8
				if ti == 0 {
					score = 1.0
				}
				if tv > 0 {
					score -= 0.05
				}
				if vi > 0 {
					score -= 0.1
				}
				out = append(out, AnimeSamaCandidate{CatalogueURL: url, Slug: v, MatchedTitle: t, Score: score})
				if len(out) >= maxCandidates {
					return out, nil
				}
			}
		}
	}

	return out, nil
}

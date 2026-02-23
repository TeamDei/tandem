package latin

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const wiktionaryBase = "https://en.wiktionary.org/api/rest_v1/page/definition/"

// tagRe strips HTML tags from Wiktionary definition strings.
var tagRe = regexp.MustCompile(`<[^>]+>`)

// Definition holds a part-of-speech heading and its English glosses.
type Definition struct {
	PartOfSpeech string
	Glosses      []string
}

// wiktEntry is the per-language entry shape returned by the Wiktionary API.
type wiktEntry struct {
	PartOfSpeech string `json:"partOfSpeech"`
	Definitions  []struct {
		Definition string `json:"definition"`
	} `json:"definitions"`
}

// fetchWiktionary queries the Wiktionary REST API for word and returns
// any Latin-language entries found.
func fetchWiktionary(word string) ([]wiktEntry, error) {
	word = strings.ToLower(punctRe.ReplaceAllString(word, ""))
	if word == "" {
		return nil, nil
	}

	url := wiktionaryBase + word
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("wiktionary: building request: %w", err)
	}

	// Wikimedia requires a descriptive User-Agent; bare Go agents get 403'd.
	req.Header.Set("User-Agent", "tandem/1.0 (https://github.com/teamdei/tandem; latin-reader) go-net/http")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wiktionary: GET %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // word simply not in Wiktionary
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wiktionary: unexpected status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("wiktionary: reading body: %w", err)
	}

	// The response is a JSON object keyed by language code, e.g. {"la": [...], "en": [...]}
	var raw map[string][]wiktEntry
	if err = json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("wiktionary: parsing JSON: %w", err)
	}

	// Latin entries live under "la".
	entries, ok := raw["la"]
	if !ok || len(entries) == 0 {
		return nil, nil // word simply not in Wiktionary
	}
	return entries, nil
}

// formOfRe matches the lemma link inside a Wiktionary form-of definition.
var formOfRe = regexp.MustCompile(`class="form-of-definition-link"[^>]*>.*?<a[^>]*href="/wiki/([^"#]+)(?:#[^"]*)?"`)

// LookupDefinitions queries the Wiktionary REST API. If the word only contains
// inflected "form-of" definitions, it automatically resolves and looks up the
// base lemma instead. It returns the actual word looked up (which may be the lemma)
// and its definitions.
func LookupDefinitions(word string) (string, []Definition, error) {
	entries, err := fetchWiktionary(word)
	if err != nil || len(entries) == 0 {
		return word, nil, err
	}

	// Check if we need to resolve a form-of.
	// We only follow form-of links if ALL definitions are purely inflections,
	// or if we want to be safe, we just grab the first form-of link we find.
	// Let's see if we have ANY non-form-of definitions.
	hasSubstantive := false
	var firstFormOf string

	for _, entry := range entries {
		for _, def := range entry.Definitions {
			if strings.Contains(def.Definition, "form-of-definition") {
				if firstFormOf == "" {
					if m := formOfRe.FindStringSubmatch(def.Definition); len(m) > 1 {
						firstFormOf = m[1]
					}
				}
			} else {
				hasSubstantive = true
			}
		}
	}

	// If the word has no substantive definitions of its own, fetch its lemma.
	if !hasSubstantive && firstFormOf != "" {
		resolved, err := fetchWiktionary(firstFormOf)
		if err == nil && len(resolved) > 0 {
			entries = resolved
			word = firstFormOf
		}
	}

	var defs []Definition
	for _, entry := range entries {
		d := Definition{PartOfSpeech: strings.Title(strings.ToLower(entry.PartOfSpeech))} //nolint:staticcheck
		for _, def := range entry.Definitions {
			gloss := cleanGloss(def.Definition)
			if gloss != "" {
				d.Glosses = append(d.Glosses, gloss)
			}
		}
		if len(d.Glosses) > 0 {
			defs = append(defs, d)
		}
	}
	return word, defs, nil
}

// styleRe removes <style>...</style> blocks and their content.
var styleRe = regexp.MustCompile(`(?s)<style[^>]*>.*?</style>`)

// cleanGloss strips HTML tags and decodes HTML entities from a Wiktionary
// definition string.
func cleanGloss(s string) string {
	s = styleRe.ReplaceAllString(s, "")

	// Wikimedia often nests spans. Since regex can't handle arbitrary nesting well,
	// we iteratively remove innermost spans until none are left. This safely kills
	// <span class="usage-label-sense">.mw-parser-output...</span> without taking
	// the whole definition down with it.
	for {
		prev := s
		s = regexp.MustCompile(`<span[^>]*>[^<]*</span>`).ReplaceAllString(s, "")
		if s == prev {
			break
		}
	}

	s = tagRe.ReplaceAllString(s, "")
	s = html.UnescapeString(s)

	// Some Wiktionary entries have leftover CSS class text if it wasn't
	// inside a <style> tag, but usually the <style> tag houses it.
	// We also replace multiple spaces with a single space.
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimSpace(s)
}

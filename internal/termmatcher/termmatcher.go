package termmatcher

import (
	"log"
	"net/url"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/caneroj1/stemmer"
)

const minTitleLength = 5 // Set the minimum character limit as needed

// Related checks if the URL title matches any of the provided search terms.
func Related(href string, searchTerms []string) bool {
	title := extractTitleFromURL(href)
	if title == "" {
		log.Println("Title is empty for URL:", href)
		return false
	}

	processedTitle := processTitle(title)
	if processedTitle == "" {
		log.Println("Processed title is empty for URL:", href)
		return false
	}

	// Check if the title meets the minimum character limit
	if len(processedTitle) < minTitleLength {
		log.Println("Processed title is shorter than minimum length for URL:", href)
		return false
	}

	isMatch := matchSearchTerms(processedTitle, searchTerms)
	if isMatch {
		log.Println("Match found for URL:", href)
	} else {
		log.Println("No match found for URL:", href)
	}

	return isMatch
}

// extractTitleFromURL extracts the title part from a URL.
func extractTitleFromURL(urlString string) string {
	u, err := url.Parse(urlString)
	if err != nil {
		// Handle the error, e.g., log it or return an error string
		return ""
	}

	// Check if the URL has a path component
	if u.Path == "" || u.Path == "/" {
		// If there's no path component, return an empty string
		return ""
	}

	// Split the path and get the last segment as the title
	segments := strings.Split(u.Path, "/")
	title := segments[len(segments)-1]

	return title
}

// processTitle processes the title by removing '-', stopwords, and stemming.
func processTitle(title string) string {
	log.Println("Original title:", title)
	title = strings.ReplaceAll(title, "-", " ")
	title = stopwords.CleanString(title, "en", false)
	title = strings.TrimSpace(title)

	words := strings.Split(title, " ")
	words = stemmer.StemMultiple(words)

	// Lemmatize (if needed)
	// ...

	processedTitle := strings.Join(words, " ")
	log.Println("Processed title:", processedTitle)

	return processedTitle
}

// matchSearchTerms checks if the processed title matches any of the search terms.
func matchSearchTerms(title string, searchTerms []string) bool {
	swg := metrics.NewSmithWatermanGotoh()
	swg.CaseSensitive = false
	swg.GapPenalty = -0.1
	swg.Substitution = metrics.MatchMismatch{
		Match:    1,
		Mismatch: -0.5,
	}

	title = strings.ToLower(title)

	for _, term := range searchTerms {
		term = strings.ToLower(term)
		term = stemmer.Stem(term)
		similarity := strutil.Similarity(term, title, swg)
		if similarity == 1 {
			return true
		}
	}

	return false
}

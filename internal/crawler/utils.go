package crawler

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
)

// extractHostFromURL extracts the host from the given URL.
// It uses the GetHostFromURL function to parse the URL and retrieve the host.
// If the URL cannot be parsed, it logs an error and returns an empty string along with the error.
// Parameters:
// - url: The URL from which to extract the host.
// Returns:
// - string: The extracted host from the URL.
// - error: An error if the URL cannot be parsed.
func (cm *CrawlManager) extractHostFromURL(url string) (string, error) {
	host, err := GetHostFromURL(url, cm.Logger())
	if err != nil {
		cm.LoggerField.Error(fmt.Sprintf("Failed to parse URL: url: %v, error: %v", url, err), nil)
		return "", err
	}
	cm.LoggerField.Debug(fmt.Sprintf("Extracted host from URL: %s", host))
	return host, nil
}

func (cm *CrawlManager) GetAnchorElementHandler() func(*colly.HTMLElement) error {
	return func(e *colly.HTMLElement) error {
		start := time.Now()
		href := cm.getHref(e)
		if href == "" || strings.HasSuffix(href, ".png") || strings.HasSuffix(href, ".jpg") || strings.HasSuffix(href, ".jpeg") {
			return nil // Skip if href is empty or is an image file
		}

		cm.processLink(e, href)
		err := cm.visitWithColly(href)
		if err != nil {
			elapsed := time.Since(start)
			cm.LoggerField.Debug(fmt.Sprintf("[GetAnchorElementHandler] Error visiting URL: %s, Error: %v, Elapsed time: %s", href, err, elapsed))
			return err // Return the error to propagate it
		}
		elapsed := time.Since(start)
		cm.LoggerField.Debug(fmt.Sprintf("[GetAnchorElementHandler] Visited URL: %s, Elapsed time: %s", href, elapsed))
		return nil // No error occurred
	}
}

// getHref retrieves the href attribute from the given HTML element.
// It logs a debug message if the href is empty or if the anchor element has no href attribute.
// Parameters:
// - e: The HTML element to retrieve the href from.
// Returns:
// - string: The href attribute value.
func (cm *CrawlManager) getHref(e *colly.HTMLElement) string {
	href := e.Request.AbsoluteURL(e.Attr("href"))
	if href == "" {
		cm.LoggerField.Debug("Found anchor element with no href attribute")
	} else {
		cm.LoggerField.Debug(fmt.Sprintf("Processing link: %s", href))
	}
	return href
}

// incrementTotalLinks increments the total links count in the StatsManager.
// It logs a debug message indicating the total links count has been incremented.
func (cm *CrawlManager) incrementTotalLinks() {
	cm.StatsManager.LinkStats.IncrementTotalLinks()
	cm.LoggerField.Debug("Incremented total links count")
}

// logCurrentURL logs the current URL being crawled.
// Parameters:
// - e: The HTML element representing the current page.
func (cm *CrawlManager) logCurrentURL(e *colly.HTMLElement) {
	cm.LoggerField.Debug(fmt.Sprintf("Current URL being crawled: %s", e.Request.URL.String()))
}

// createPageData creates a PageData instance with the given href.
// Parameters:
// - href: The URL to create the PageData instance for.
// Returns:
// - PageData: A PageData instance with the URL set.
func (cm *CrawlManager) createPageData(href string) PageData {
	return PageData{
		URL: href,
	}
}

// logSearchTerms logs the search terms used for crawling.
// Parameters:
// - options: The CrawlOptions containing the search terms.
func (cm *CrawlManager) logSearchTerms() {
	cm.LoggerField.Debug(fmt.Sprintf("Search terms: %v", cm.Options.SearchTerms))
}

// getMatchingTerms retrieves the matching terms from the given href and anchor text.
// Parameters:
// - href: The URL to check for matching terms.
// - anchorText: The text of the anchor element.
// - options: The CrawlOptions containing the search terms.
// Returns:
// - []string: A slice of strings representing the matching terms.
func (cm *CrawlManager) getMatchingTerms(href string, anchorText string, options *CrawlOptions) []string {
	return termmatcher.GetMatchingTerms(href, anchorText, options.SearchTerms, cm.Logger().(*loggo.Logger))
}

// handleMatchingTerms processes the matching terms and updates the stats.
// Parameters:
// - options: The CrawlOptions containing the search terms.
// - currentURL: The current URL being crawled.
// - pageData: The PageData instance for the current URL.
// - matchingTerms: A slice of strings representing the matching terms.
func (cm *CrawlManager) handleMatchingTerms(options *CrawlOptions, currentURL string, pageData PageData, matchingTerms []string) {
	cm.ProcessMatchingLink(currentURL, pageData, matchingTerms)
	cm.UpdateStats(options, matchingTerms)
}

func (cm *CrawlManager) processLink(e *colly.HTMLElement, href string) {
	start := time.Now()
	cm.incrementTotalLinks()
	cm.logCurrentURL(e)
	pageData := cm.createPageData(href)
	cm.logSearchTerms()
	matchingTerms := cm.getMatchingTerms(href, e.Text, cm.Options)
	cm.handleMatchingTerms(cm.Options, e.Request.URL.String(), pageData, matchingTerms)
	elapsed := time.Since(start)
	cm.LoggerField.Debug(fmt.Sprintf("[processLink] Processed link: %s, Elapsed time: %s", href, elapsed))
}

// ProcessMatchingLink processes a matching link by updating the page data and appending the result.
// Parameters:
// - href: The URL of the matching link.
// - pageData: The PageData instance for the matching link.
// - matchingTerms: A slice of strings representing the matching terms.
func (cm *CrawlManager) ProcessMatchingLink(href string, pageData PageData, matchingTerms []string) {
	if href == "" {
		cm.LoggerField.Error("Missing URL for matching link", nil)
		return
	}

	pageData.UpdatePageData(href, matchingTerms)
	cm.AppendResult(pageData)
}

// UpdateStats updates the stats for matched and non-matched links.
// Parameters:
// - options: The CrawlOptions containing the search terms.
// - matchingTerms: A slice of strings representing the matching terms.
func (cm *CrawlManager) UpdateStats(_ *CrawlOptions, matchingTerms []string) {
	if len(matchingTerms) > 0 {
		cm.incrementMatchedLinks()
		cm.LoggerField.Debug("Incremented matched links count")
	} else {
		cm.incrementNonMatchedLinkCount()
		cm.LoggerField.Debug("Incremented not matched links count")
	}
}

// incrementMatchedLinks increments the matched links count in the StatsManager.
func (cm *CrawlManager) incrementMatchedLinks() {
	cm.StatsManager.LinkStats.IncrementMatchedLinks()
}

// incrementNonMatchedLinkCount increments the non-matched links count in the StatsManager.
func (cm *CrawlManager) incrementNonMatchedLinkCount() {
	cm.StatsManager.LinkStats.IncrementNotMatchedLinks()
	cm.LoggerField.Debug("Incremented not matched links count")
}

// GetHostFromURL extracts the host from the given URL.
// Parameters:
// - inputURL: The URL to parse and extract the host from.
// - appLogger: The loggo instance to log any errors.
// Returns:
// - string: The extracted host from the URL.
// - error: An error if the URL cannot be parsed or if the host cannot be extracted.
func GetHostFromURL(inputURL string, appLogger loggo.LoggerInterface) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		appLogger.Error(fmt.Sprintf("Failed to parse URL: %s, Error: %v", inputURL, err), nil)
		return "", err
	}

	host := parsedURL.Hostname()
	if host == "" {
		appLogger.Error("failed to extract host from URL", nil)
		return "", errors.New("failed to extract host from URL")
	}

	return host, nil
}

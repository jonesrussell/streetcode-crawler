package crawler

import (
	"strings"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
)

func (cs *CrawlManager) getAnchorElementHandler(options *CrawlOptions) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		href := cs.getHref(e)
		if href == "" {
			return
		}

		cs.processLink(e, href, options)
		err := cs.visitWithColly(href)
		if err != nil {
			cs.Debug("[getAnchorElementHandler] Error visiting URL", "url", href, "error", err)
		}
	}
}

func (cs *CrawlManager) getHref(e *colly.HTMLElement) string {
	href := e.Request.AbsoluteURL(e.Attr("href"))
	if href == "" {
		cs.Debug("Found anchor element with no href attribute")
	} else {
		cs.Debug("Processing link", "href", href)
	}
	return href
}

func (cs *CrawlManager) incrementTotalLinks(options *CrawlOptions) {
	cs.StatsManager.LinkStats.IncrementTotalLinks()
	cs.Debug("Incremented total links count")
}

func (cs *CrawlManager) logCurrentURL(e *colly.HTMLElement) {
	cs.Debug("Current URL being crawled", "url", e.Request.URL.String())
}

func (cs *CrawlManager) createPageData(href string) PageData {
	return PageData{
		URL: href,
	}
}

func (cs *CrawlManager) logSearchTerms(options *CrawlOptions) {
	cs.Debug("Search terms", "terms", options.SearchTerms)
}

func (cs *CrawlManager) getMatchingTerms(href string, anchorText string, options *CrawlOptions) []string {
	return termmatcher.GetMatchingTerms(href, anchorText, options.SearchTerms, cs.Logger())
}

func (cs *CrawlManager) handleMatchingTerms(options *CrawlOptions, currentURL string, pageData PageData, matchingTerms []string) {
	if len(matchingTerms) > 0 {
		cs.ProcessMatchingLinkAndUpdateStats(options, currentURL, pageData, matchingTerms)
	} else {
		cs.incrementNonMatchedLinkCount(options)
		cs.Debug("Link does not match search terms", "link", pageData.URL)
	}
}

func (cs *CrawlManager) processLink(e *colly.HTMLElement, href string, options *CrawlOptions) {
	cs.incrementTotalLinks(options)
	cs.logCurrentURL(e)
	pageData := cs.createPageData(href)
	cs.logSearchTerms(options)
	matchingTerms := cs.getMatchingTerms(href, e.Text, options)
	cs.handleMatchingTerms(options, e.Request.URL.String(), pageData, matchingTerms)
}

// handleMatchingLinks is responsible for handling the links that match the search criteria during crawling.
func (cs *CrawlManager) handleMatchingLinks(href string) error {
	cs.Debug("[handleMatchingLinks] Start handling matching links", "url", href)

	cs.Debug("End handling matching links", "url", href)
	return nil
}

func (cs *CrawlManager) handleSetupError(err error) error {
	cs.Error("Error setting up crawling logic", "error", err)
	return err
}

func (cs *CrawlManager) ProcessMatchingLinkAndUpdateStats(options *CrawlOptions, href string, pageData PageData, matchingTerms []string) {
	if cs == nil {
		cs.Error("CrawlManager instance is nil")
		return
	}

	if href == "" {
		cs.Error("Missing URL for matching link")
		return
	}

	cs.StatsManager.LinkStats.IncrementMatchedLinks()
	cs.Debug("Incremented matched links count")

	if err := cs.MatchedLinkProcessor.HandleMatchingLinks(href); err != nil {
		cs.Error("[ProcessMatchingLinkAndUpdateStats] Error handling matching links", "error", err)
		return
	}

	cs.MatchedLinkProcessor.UpdatePageData(&pageData, href, matchingTerms)
	cs.MatchedLinkProcessor.AppendResult(options, pageData)
}

func (cs *CrawlManager) incrementMatchedLinks(options *CrawlOptions, linkStats *stats.Stats) {
	linkStats.IncrementMatchedLinks()
}

func (cs *CrawlManager) updatePageData(pageData *PageData, href string, matchingTerms []string) {
	pageData.MatchingTerms = matchingTerms
	pageData.ParentURL = href // Store the parent URL
}

func (cs *CrawlManager) appendResult(options *CrawlOptions, pageData PageData) {
	*options.Results = append(*options.Results, pageData)
}

func (cs *CrawlManager) incrementNonMatchedLinkCount(options *CrawlOptions) {
	cs.StatsManager.LinkStats.IncrementNotMatchedLinks()
	cs.Debug("Incremented not matched links count")
}

func (cs *CrawlManager) createLimitRule() *colly.LimitRule {
	return &colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: DefaultParallelism,
		Delay:       DefaultDelay,
	}
}

func (cs *CrawlManager) splitSearchTerms(searchTerms string) []string {
	terms := strings.Split(searchTerms, ",")
	var validTerms []string
	for _, term := range terms {
		if term != "" {
			validTerms = append(validTerms, term)
		}
	}
	return validTerms
}

func (cs *CrawlManager) createStartCrawlingOptions(crawlSiteID string, searchTerms []string, debug bool) *CrawlOptions {
	var results []PageData
	return NewCrawlOptions(crawlSiteID, searchTerms, debug, &results, stats.NewStats())
}

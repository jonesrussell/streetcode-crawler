package crawler

import (
	"context"
	"errors"
	"strings"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/models"
)

func (cm *CrawlManager) getHref(e *colly.HTMLElement) (string, error) {
	href := e.Request.AbsoluteURL(e.Attr("href"))
	if href == "" {
		return "", errors.New("Found anchor element with no href attribute")
	}
	return href, nil
}

func (cm *CrawlManager) createPageData(href string) models.PageData {
	return models.PageData{
		URL: href,
	}
}

func (cm *CrawlManager) handleMatchingTerms(options *CrawlOptions, currentURL string, pageData models.PageData, matchingTerms []string) error {
	cm.Logger.Debug("handleMatchingTerms called")

	// Calculate the similarity score
	similarityScore := cm.TermMatcher.CompareTerms(currentURL, strings.Join(matchingTerms, " "))

	pageData.UpdatePageData(matchingTerms, similarityScore) // Update the PageData with the similarity score

	// Append the PageData directly to Results.Pages
	cm.Results.Pages = append(cm.Results.Pages, pageData)

	cm.UpdateStats(options, matchingTerms)

	// Save the result to Redis
	ctx := context.Background() // Or use a context from your application
	key := options.CrawlSiteID

	err := cm.DBManager.SaveResultsToRedis(ctx, []models.PageData{pageData}, key)

	if err != nil {
		cm.Logger.Error("Error saving result to Redis: ", err)
		return err
	}

	return nil
}

func (cm *CrawlManager) processLink(e *colly.HTMLElement, href string) error {
	cm.StatsManager.LinkStats.IncrementTotalLinks()
	pageData := cm.createPageData(href)
	matchingTerms := cm.TermMatcher.GetMatchingTerms(href, e.Text, cm.Options.SearchTerms)
	if len(matchingTerms) > 0 {
		err := cm.handleMatchingTerms(cm.Options, e.Request.URL.String(), pageData, matchingTerms)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cm *CrawlManager) UpdateStats(_ *CrawlOptions, matchingTerms []string) {
	if len(matchingTerms) > 0 {
		cm.StatsManager.LinkStats.IncrementMatchedLinks()
	} else {
		cm.StatsManager.LinkStats.IncrementNotMatchedLinks()
	}
}
package crawler

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/internal/stats"
)

type CrawlManager struct {
	LoggerField    logger.Logger
	Client         prowlredis.ClientInterface
	MongoDBWrapper mongodbwrapper.MongoDBInterface
	Collector      *colly.Collector
	CrawlingMu     *sync.Mutex
	StatsManager   *StatsManager
}

// NewCrawlManager creates a new instance of CrawlManager.
func NewCrawlManager(
	loggerField logger.Logger,
	client prowlredis.ClientInterface,
	mongoDBWrapper mongodbwrapper.MongoDBInterface,
) *CrawlManager {
	return &CrawlManager{
		LoggerField:    loggerField,
		Client:         client,
		MongoDBWrapper: mongoDBWrapper,
		CrawlingMu:     &sync.Mutex{},
	}
}

type StatsManager struct {
	LinkStats   *stats.Stats
	LinkStatsMu sync.RWMutex
}

// NewStatsManager creates a new StatsManager with initialized fields.
func NewStatsManager() *StatsManager {
	return &StatsManager{
		LinkStats:   &stats.Stats{},
		LinkStatsMu: sync.RWMutex{},
	}
}

// Crawl starts the crawling process for a given URL with the provided options.
// It logs the URL being crawled, sets up the crawling logic, visits the URL, and returns the crawling results.
// Parameters:
// - url: The URL to start crawling.
// - options: The CrawlOptions containing configuration for the crawling process.
// Returns:
// - []PageData: A slice of PageData representing the crawling results.
// - error: An error if the crawling process encounters any issues.
func (cm *CrawlManager) Crawl(url string, options *CrawlOptions) ([]PageData, error) {
	cm.LoggerField.Debug(fmt.Sprintf("CrawlURL: %s", url))
	err := cm.SetupCrawlingLogic(options)
	if err != nil {
		return nil, err
	}

	err = cm.CrawlURL(url)
	if err != nil {
		return nil, err
	}

	return *options.Results, nil
}

// HandleVisitError handles the error occurred during the visit of a URL.
// It logs the error and returns it.
// Parameters:
// - url: The URL that encountered an error during the visit.
// - err: The error that occurred during the visit.
// Returns:
// - error: The error that was logged and returned.
func (cm *CrawlManager) HandleVisitError(url string, err error) error {
	cm.LoggerField.Error(fmt.Sprintf("Error visiting URL: url: %s, error: %v", url, err))
	return err
}

// StartCrawling initiates the crawling process with the given parameters.
// It validates the input parameters, configures the collector, and starts the crawling process.
// It returns an error if the crawling process fails to start.
// Parameters:
// - ctx: The context for the crawling operation.
// - url: The URL to start crawling from.
// - searchTerms: The search terms to match against the crawled content.
// - crawlSiteID: The ID of the site to crawl.
// - maxDepth: The maximum depth to crawl.
// - debug: A flag indicating whether to enable debug mode for the crawling process.
func (cm *CrawlManager) StartCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool) error {
	if err := cm.validateParameters(url, searchTerms, crawlSiteID, maxDepth); err != nil {
		return err
	}

	cm.initializeStatsManager()

	host, err := cm.extractHostFromURL(url)
	if err != nil {
		return err
	}

	if err := cm.configureCollector(host, maxDepth); err != nil {
		return err
	}

	options := cm.createCrawlingOptions(crawlSiteID, searchTerms, debug)

	return cm.performCrawling(ctx, url, options)
}

func (cm *CrawlManager) ConfigureCollector(allowedDomains []string, maxDepth int) error {
	cm.Collector = colly.NewCollector(
		colly.Async(false),
		colly.MaxDepth(maxDepth),
		colly.Debugger(cm.LoggerField),
	)

	cm.LoggerField.Debug(fmt.Sprintf("Allowed Domains: %v", allowedDomains))
	cm.Collector.AllowedDomains = allowedDomains

	limitRule := cm.createLimitRule()
	if err := cm.Collector.Limit(limitRule); err != nil {
		cm.LoggerField.Error(fmt.Sprintf("Failed to set limit rule: %v", err))
		return err
	}

	// Respect robots.txt
	cm.Collector.AllowURLRevisit = false
	cm.Collector.IgnoreRobotsTxt = false

	// Register OnScraped callback
	cm.Collector.OnScraped(func(r *colly.Response) {
		cm.LoggerField.Debug(fmt.Sprintf("[OnScraped] Page scraped: %s", r.Request.URL.String()))
		cm.StatsManager.LinkStatsMu.Lock()
		defer cm.StatsManager.LinkStatsMu.Unlock()
		cm.StatsManager.LinkStats.IncrementTotalPages()
	})

	return nil
}

func (cm *CrawlManager) logCrawlingStatistics() {
	report := cm.StatsManager.LinkStats.Report()
	infoMessage := fmt.Sprintf("Crawling statistics: TotalLinks=%v, MatchedLinks=%v, NotMatchedLinks=%v, TotalPages=%v",
		report["TotalLinks"], report["MatchedLinks"], report["NotMatchedLinks"], report["TotalPages"])
	cm.LoggerField.Info(infoMessage)
}

func (cm *CrawlManager) visitWithColly(url string) error {
	cm.LoggerField.Debug(fmt.Sprintf("[visitWithColly] Visiting URL with Colly: %v", url))

	err := cm.Collector.Visit(url)
	if err != nil {
		switch {
		case errors.Is(err, colly.ErrAlreadyVisited):
			errorMessage := fmt.Sprintf("[visitWithColly] URL already visited: %v", url)
			cm.LoggerField.Debug(errorMessage)
		case errors.Is(err, colly.ErrForbiddenDomain):
			errorMessage := fmt.Sprintf("[visitWithColly] Forbidden domain - Skipping visit: %v", url)
			cm.LoggerField.Debug(errorMessage)
		default:
			errorMessage := fmt.Sprintf("[visitWithColly] Error visiting URL: url=%v, error=%v", url, err)
			cm.LoggerField.Error(errorMessage)
		}
		return nil
	}

	successMessage := fmt.Sprintf("[visitWithColly] Successfully visited URL: %v", url)
	cm.LoggerField.Debug(successMessage)
	return nil
}

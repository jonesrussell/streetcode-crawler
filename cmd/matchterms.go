package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/crawlresult"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var matchtermsCmd = &cobra.Command{
	Use:   "matchterms",
	Short: "Crawl websites and extract information",
	Long: `Crawl is a CLI tool designed to perform web scraping and data extraction from websites.
           It allows users to specify parameters such as depth of crawl and target elements to extract.`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("debug") {
			fmt.Println("\nFlags:")
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				fmt.Printf("  %-12s : %v\n", flag.Name, flag.Value)
			})

			fmt.Println("\nRedis Environment Variables:")
			fmt.Printf("  %-12s : %s\n", "REDIS_HOST", viper.GetString("REDIS_HOST"))
			fmt.Printf("  %-12s : %s\n", "REDIS_PORT", viper.GetString("REDIS_PORT"))
			fmt.Printf("  %-12s : %s\n", "REDIS_AUTH", viper.GetString("REDIS_AUTH"))
		}

		ctx := context.Background()
		crawlerService, err := initializeManager(ctx, viper.GetBool("debug"))
		if err != nil {
			fmt.Println("Failed to initialize Crawl Manager", "error", err)
			os.Exit(1)
		}

		startCrawling(ctx, viper.GetString("url"), viper.GetString("searchterms"), viper.GetString("crawlsiteid"), viper.GetInt("maxdepth"), viper.GetBool("debug"), crawlerService)
	},
}

func init() {
	matchtermsCmd.Flags().String("url", "", "URL to crawl")
	matchtermsCmd.Flags().String("searchterms", "", "Search terms for crawling")
	matchtermsCmd.Flags().Int("maxdepth", 1, "Maximum depth for crawling")

	viper.BindPFlag("url", matchtermsCmd.Flags().Lookup("url"))
	viper.BindPFlag("searchterms", matchtermsCmd.Flags().Lookup("searchterms"))
	viper.BindPFlag("maxdepth", matchtermsCmd.Flags().Lookup("maxdepth"))

	rootCmd.AddCommand(matchtermsCmd)
}

func startCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool, crawlerService *crawler.CrawlManager) error {
	splitSearchTerms := strings.Split(searchTerms, ",")
	host, err := crawler.GetHostFromURL(url, crawlerService.Logger)
	if err != nil {
		crawlerService.Logger.Error("Failed to parse URL", "url", url, "error", err)
		return err
	}

	collector := crawler.ConfigureCollector([]string{host}, maxDepth)
	if collector == nil {
		crawlerService.Logger.Fatal("Failed to configure collector")
		return errors.New("failed to configure collector")
	}

	var results []crawlresult.PageData

	options := crawler.CrawlOptions{
		CrawlSiteID: crawlSiteID,
		Collector:   collector,
		SearchTerms: splitSearchTerms,
		Results:     &results,
		LinkStats:   stats.NewStats(),
		Debug:       debug,
	}
	crawlerService.SetupCrawlingLogic(ctx, &options)

	crawlerService.Logger.Info("Crawler started...")
	if err := collector.Visit(url); err != nil {
		crawlerService.Logger.Error("Error visiting URL", "url", url, "error", err)
		return err
	}

	collector.Wait()

	crawlerService.Logger.Info("Crawling completed.")

	err = saveResultsToRedis(ctx, crawlerService, results)
	if err != nil {
		return err
	}
	printResults(crawlerService, results)

	return nil
}

func saveResultsToRedis(ctx context.Context, crawlerService *crawler.CrawlManager, results []crawlresult.PageData) error {
	for _, result := range results {
		_, err := crawlerService.RedisClient.SAdd(ctx, "yourKeyHere", result)
		if err != nil {
			crawlerService.Logger.Error("Error occurred during saving to Redis", "error", err)
			return err
		}
	}
	return nil
}

func printResults(crawlerService *crawler.CrawlManager, results []crawlresult.PageData) {
	jsonData, err := json.Marshal(results)
	if err != nil {
		crawlerService.Logger.Error("Error occurred during marshaling", "error", err)
		return
	}

	fmt.Println(string(jsonData))
}

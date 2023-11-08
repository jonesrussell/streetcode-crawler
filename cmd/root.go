// Package cmd contains the command-line commands for the crawler application.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jonesrussell/crawler/internal/crawler"
	"github.com/jonesrussell/crawler/internal/crawlresult"
	"github.com/jonesrussell/crawler/internal/rediswrapper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl websites and extract information",
	Long: `Crawl is a CLI tool designed to perform web scraping and data extraction from websites.
           It allows users to specify parameters such as depth of crawl and target elements to extract.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		startCrawling(ctx, viper.GetString("url"), viper.GetString("searchterms"), viper.GetString("crawlsiteid"), viper.GetInt("maxdepth"), viper.GetBool("debug"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("url", "", "URL to crawl")
	rootCmd.PersistentFlags().String("searchterms", "", "Comma-separated search terms")
	rootCmd.PersistentFlags().String("crawlsiteid", "", "CrawlSite ID")
	rootCmd.PersistentFlags().Int("maxdepth", 1, "Maximum depth for the crawler")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.MarkPersistentFlagRequired("searchterms")
	rootCmd.MarkPersistentFlagRequired("crawlsiteid")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("searchterms", rootCmd.PersistentFlags().Lookup("searchterms"))
	viper.BindPFlag("crawlsiteid", rootCmd.PersistentFlags().Lookup("crawlsiteid"))
	viper.BindPFlag("maxdepth", rootCmd.PersistentFlags().Lookup("maxdepth"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv() // Automatically override values from the .env file with those from the environment.

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error while reading config file", err)
	}
}

func startCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool) {
	logger, err := crawler.InitializeLogger(debug)
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}

	// Use Viper to get Redis configuration directly
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisAuth := viper.GetString("REDIS_AUTH")

	redisWrapper, err := rediswrapper.NewRedisWrapper(ctx, redisHost, redisPort, redisAuth, logger)
	if err != nil {
		logger.Errorf("Failed to initialize Redis: %v", err)
		os.Exit(1)
	}

	splitSearchTerms := strings.Split(searchTerms, ",")
	collector := crawler.ConfigureCollector([]string{crawler.GetHostFromURL(url)}, maxDepth)
	if collector == nil {
		logger.Fatal("Failed to configure collector")
		return
	}

	var results []crawlresult.PageData
	crawler.SetupCrawlingLogic(ctx, crawlSiteID, collector, splitSearchTerms, &results, logger, redisWrapper)

	logger.Info("Crawler started...")
	if err := collector.Visit(url); err != nil {
		logger.Error("Error visiting URL", zap.Error(err))
		return
	}

	collector.Wait()

	jsonData, err := json.Marshal(results)
	if err != nil {
		logger.Error("Error occurred during marshaling", zap.Error(err))
		return
	}

	fmt.Println(string(jsonData))
	logger.Info("Crawling completed.")
}

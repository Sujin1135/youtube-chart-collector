package main

import (
	"channel-collector/internal/chart"
	"channel-collector/internal/collector"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	// Common flags
	mode       = flag.String("mode", "chart", "Collection mode: 'channel' or 'chart'")
	outputFile = flag.String("output", "", "Output file path (JSON format). If empty, prints to stdout")

	// Channel collection flags
	channelHandles = flag.String("channels", "", "Comma-separated YouTube channel handles (e.g., @channel1,@channel2)")

	// Chart collection flags
	chartType = flag.String("chart-type", "TrendingVideos", "Chart type: TrendingVideos, TopVideos, TopArtists, TopSongs")
	region    = flag.String("region", "kr", "Region code: kr, us, global")
	timeRange = flag.String("time-range", "RightNow", "Time range: RightNow, Daily, Weekly")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	runChartCollection()
}

func runChartCollection() {
	var configs []*chart.ChartConfig

	chartTypeEnum := parseChartType(*chartType)
	regionEnum := parseRegion(*region)
	timeRangeEnum := parseTimeRange(*timeRange)

	configs = []*chart.ChartConfig{
		{
			Type:      chartTypeEnum,
			Region:    regionEnum,
			TimeRange: timeRangeEnum,
		},
	}
	log.Printf("Starting chart collection: %s/%s/%s\n", *chartType, *region, *timeRange)

	resultCh := make(chan *chart.Chart, len(configs))
	errCh := make(chan error, len(configs))

	var wg sync.WaitGroup
	wg.Add(1)

	chartCollector := collector.NewChartCollector()
	go chartCollector.Collect(configs, resultCh, errCh, &wg)

	wg.Wait()

	// Collect results
	var charts []*chart.Chart
	for chartData := range resultCh {
		charts = append(charts, chartData)
		log.Printf("Collected chart: %s/%s/%s with %d videos\n",
			chartData.Type, chartData.Region, chartData.TimeRange, len(chartData.Videos))
	}

	// Collect errors
	var errors []string
	for err := range errCh {
		errors = append(errors, err.Error())
		log.Printf("Error: %v\n", err)
	}

	// Output results
	result := map[string]interface{}{
		"mode":   "chart",
		"total":  len(charts),
		"charts": charts,
		"errors": errors,
	}

	outputResult(result)
	log.Printf("Chart collection completed. Collected: %d charts, Errors: %d\n", len(charts), len(errors))
}

func outputResult(result interface{}) {
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling result: %v\n", err)
	}

	if *outputFile != "" {
		// Write to file
		err := os.WriteFile(*outputFile, jsonData, 0644)
		if err != nil {
			log.Fatalf("Error writing to file: %v\n", err)
		}
		log.Printf("Results written to: %s\n", *outputFile)
	} else {
		// Print to stdout
		fmt.Println("\n=== Results ===")
		fmt.Println(string(jsonData))
	}
}

func parseChartType(s string) chart.ChartType {
	switch strings.ToLower(s) {
	case "trendingvideos":
		return chart.TrendingVideos
	case "topvideos":
		return chart.TopVideos
	case "topartists":
		return chart.TopArtists
	case "topsongs":
		return chart.TopSongs
	default:
		log.Fatalf("Unknown chart type: %s", s)
		return chart.TrendingVideos
	}
}

func parseRegion(s string) chart.Region {
	switch strings.ToLower(s) {
	case "kr":
		return chart.RegionKR
	case "us":
		return chart.RegionUS
	case "global":
		return chart.RegionGlobal
	default:
		log.Fatalf("Unknown region: %s", s)
		return chart.RegionKR
	}
}

func parseTimeRange(s string) chart.TimeRange {
	switch strings.ToLower(s) {
	case "rightnow":
		return chart.RightNow
	case "daily":
		return chart.Daily
	case "weekly":
		return chart.Weekly
	default:
		log.Fatalf("Unknown time range: %s", s)
		return chart.RightNow
	}
}

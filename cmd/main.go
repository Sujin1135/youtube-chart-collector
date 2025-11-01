package main

import (
	"channel-collector/internal/collector"
	"flag"
	"fmt"
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	runChartCollection()
}

func runChartCollection() {
	videoCollector := collector.NewChartCollector()

	topVideos, err := videoCollector.CollectTopVideos()
	if err != nil {
		log.Printf("Error collecting top videos: %v", err)
	} else {
		for _, video := range topVideos {
			fmt.Printf("Rank: %d | Title: %s | Artist: %s | URL: %s | Peak: %d | Weeks: %d\n",
				video.Rank, video.Title, video.Artist, video.VideoURL, video.PeakPosition, video.WeeksOnChart)
		}
	}

	fmt.Println("\n=== Trending Videos (Right Now) ===")
	trendingVideos, err := videoCollector.CollectTrendingVideos()
	if err != nil {
		log.Fatal(err)
	}

	for _, video := range trendingVideos {
		fmt.Printf("Rank: %d | Title: %s | Artist: %s | URL: %s\n",
			video.Rank, video.Title, video.Artist, video.VideoURL)
	}
}

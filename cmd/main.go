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

	// Collect Weekly Top Songs (주간 인기곡)
	fmt.Println("=== Weekly Top Songs ===")
	weeklyTopSongs, err := videoCollector.CollectWeeklyTopSongs()
	if err != nil {
		log.Printf("Error collecting weekly top songs: %v", err)
	} else {
		for _, song := range weeklyTopSongs {
			fmt.Printf("Rank: %d | Title: %s | Artist: %s | Views: %d | URL: %s | Peak: %d | Weeks: %d\n",
				song.Rank, song.Title, song.Artist, song.ViewCount, song.VideoURL, song.PeakPosition, song.WeeksOnChart)
		}
	}
}

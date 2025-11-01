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

	// Collect Daily Top Shorts Songs (일간 SHORTS 인기곡)
	fmt.Println("=== Daily Top Shorts Songs ===")
	dailyTopShortsSongs, err := videoCollector.CollectDailyTopShortsSongs()
	if err != nil {
		log.Printf("Error collecting daily top shorts songs: %v", err)
	} else {
		for _, song := range dailyTopShortsSongs {
			fmt.Printf("Rank: %d | Title: %s | Artist: %s | Views: %d | URL: %s | Peak: %d | Weeks: %d\n",
				song.Rank, song.Title, song.Artist, song.ViewCount, song.VideoURL, song.PeakPosition, song.WeeksOnChart)
		}
	}
}

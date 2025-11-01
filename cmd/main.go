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

	// Collect Weekly Top Artists (주간 인기 아티스트)
	fmt.Println("=== Weekly Top Artists ===")
	weeklyTopArtists, err := videoCollector.CollectWeeklyTopArtists()
	if err != nil {
		log.Printf("Error collecting weekly top artists: %v", err)
	} else {
		for _, artist := range weeklyTopArtists {
			fmt.Printf("Rank: %d | Artist: %s | Views: %d | Peak: %d | Weeks: %d\n",
				artist.Rank, artist.ArtistName, artist.ViewCount, artist.PeakPosition, artist.WeeksOnChart)
		}
	}
}

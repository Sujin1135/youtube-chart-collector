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
	trendingVideos, err := videoCollector.CollectTrendingVideos()
	if err != nil {
		log.Fatal(err)
	}

	for _, video := range trendingVideos {
		fmt.Println(video)
	}
}

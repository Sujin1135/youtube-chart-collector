package main

import (
	"channel-collector/internal/collector"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	{
		v1 := router.Group("/v1")
		{
			chart := v1.Group("/charts")
			chart.GET("/top_videos", func(c *gin.Context) {
				chartCollector := collector.NewChartCollector()
				videos, err := chartCollector.CollectTopVideos()
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
				}
				c.JSON(200, gin.H{
					"videos": videos,
				})
			})
			chart.GET("/trending_videos", func(c *gin.Context) {
				chartCollector := collector.NewChartCollector()
				videos, err := chartCollector.CollectTrendingVideos()
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
				}
				c.JSON(200, gin.H{
					"videos": videos,
				})
			})
			chart.GET("/weekly_top_artists", func(c *gin.Context) {
				chartCollector := collector.NewChartCollector()
				artists, err := chartCollector.CollectWeeklyTopArtists()
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
				}
				c.JSON(200, gin.H{
					"artists": artists,
				})
			})
			chart.GET("/weekly_top_songs", func(c *gin.Context) {
				chartCollector := collector.NewChartCollector()
				songs, err := chartCollector.CollectWeeklyTopSongs()
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
				}
				c.JSON(200, gin.H{
					"songs": songs,
				})
			})
			chart.GET("/daily_top_shorts", func(c *gin.Context) {
				chartCollector := collector.NewChartCollector()
				shorts, err := chartCollector.CollectDailyTopShortsSongs()
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
				}
				c.JSON(200, gin.H{
					"shorts": shorts,
				})
			})
		}
	}

	router.Run()
}

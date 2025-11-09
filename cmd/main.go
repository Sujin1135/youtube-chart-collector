package main

import (
	"channel-collector/internal/chart"
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
			chartGroup := v1.Group("/charts")
			chartGroup.GET("/top_videos/:timeRange", func(c *gin.Context) {
				chartCollector := collector.NewChartCollector()
				timeRangeParam := c.Param("timeRange")

				var timeRange chart.TimeRange
				switch timeRangeParam {
				case "RightNow":
					timeRange = chart.RightNow
				case "daily":
					timeRange = chart.Daily
				case "weekly":
					timeRange = chart.Weekly
				default:
					c.JSON(400, gin.H{"error": "invalid timeRange. must be one of: RightNow, daily, weekly"})
					return
				}

				videos, err := chartCollector.CollectTopVideos(timeRange)
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				c.JSON(200, gin.H{
					"videos": videos,
				})
			})
			chartGroup.GET("/trending_videos", func(c *gin.Context) {
				chartCollector := collector.NewChartCollector()
				videos, err := chartCollector.CollectTrendingVideos()
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
				}
				c.JSON(200, gin.H{
					"videos": videos,
				})
			})
			chartGroup.GET("/weekly_top_artists", func(c *gin.Context) {
				chartCollector := collector.NewChartCollector()
				artists, err := chartCollector.CollectWeeklyTopArtists()
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
				}
				c.JSON(200, gin.H{
					"artists": artists,
				})
			})
			chartGroup.GET("/weekly_top_songs", func(c *gin.Context) {
				chartCollector := collector.NewChartCollector()
				songs, err := chartCollector.CollectWeeklyTopSongs()
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
				}
				c.JSON(200, gin.H{
					"songs": songs,
				})
			})
			chartGroup.GET("/daily_top_shorts", func(c *gin.Context) {
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

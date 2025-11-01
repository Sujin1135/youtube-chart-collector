package chart

import "time"

// ChartType represents different types of YouTube charts
type ChartType string

const (
	TrendingVideos ChartType = "TrendingVideos"
	TopVideos      ChartType = "TopVideos"
	TopArtists     ChartType = "TopArtists"
	TopSongs       ChartType = "TopSongs"
)

// Region represents the country/region code for charts
type Region string

const (
	RegionKR     Region = "kr"
	RegionUS     Region = "us"
	RegionGlobal Region = "global"
)

// TimeRange represents the time period for the chart
type TimeRange string

const (
	RightNow TimeRange = "RightNow"
	Daily    TimeRange = "Daily"
	Weekly   TimeRange = "Weekly"
)

// ChartConfig defines the configuration for fetching a specific chart
type ChartConfig struct {
	Type      ChartType
	Region    Region
	TimeRange TimeRange
}

// Video represents a single video entry in the chart
type Video struct {
	Rank          int       `json:"rank"`
	VideoID       string    `json:"videoId"`
	Title         string    `json:"title"`
	Artist        string    `json:"artist"`
	ViewCount     int64     `json:"viewCount"`
	ThumbnailURL  string    `json:"thumbnailUrl"`
	ChannelName   string    `json:"channelName"`
	ChannelID     string    `json:"channelId"`
	RankChange    int       `json:"rankChange"` // Positive for up, negative for down
	IsNew         bool      `json:"isNew"`
	PublishedAt   time.Time `json:"publishedAt"`
}

// Chart represents a complete chart with its metadata and entries
type Chart struct {
	Type        ChartType   `json:"type"`
	Region      Region      `json:"region"`
	TimeRange   TimeRange   `json:"timeRange"`
	CollectedAt time.Time   `json:"collectedAt"`
	Videos      []Video     `json:"videos"`
}

// GenChartURL generates the URL for a specific chart configuration
func (c *ChartConfig) GenChartURL() string {
	return "https://charts.youtube.com/charts/" +
		string(c.Type) + "/" +
		string(c.Region) + "/" +
		string(c.TimeRange)
}
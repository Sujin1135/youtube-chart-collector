package chart

import "time"

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

type ChartConfig struct {
	Type      ChartType
	Region    Region
	TimeRange TimeRange
}

type Chart struct {
	Type        ChartType       `json:"type"`
	Region      Region          `json:"region"`
	TimeRange   TimeRange       `json:"timeRange"`
	CollectedAt time.Time       `json:"collectedAt"`
	Videos      []TrendingVideo `json:"videos"`
}

// GenChartURL generates the URL for a specific chart configuration
func (c *ChartConfig) GenChartURL() string {
	return "https://charts.youtube.com/charts/" +
		string(c.Type) + "/" +
		string(c.Region) + "/" +
		string(c.TimeRange)
}

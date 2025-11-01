package collector

import (
	"channel-collector/internal/chart"
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

const (
	rowContainerElemSelector = "div.ytmc-entry-row-container"
)

type ChartCollector struct {
	chromedpOptions []chromedp.ExecAllocatorOption
}

func NewChartCollector() *ChartCollector {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	return &ChartCollector{
		chromedpOptions: opts,
	}
}

func (c *ChartCollector) CollectTrendingVideos() ([]chart.TrendingVideo, error) {
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), c.chromedpOptions...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(string, ...interface{}) {}))
	defer cancel()

	config := &chart.ChartConfig{Type: chart.TrendingVideos, Region: chart.RegionKR, TimeRange: chart.RightNow}
	chartURL := config.GenChartURL()

	var videos []chart.TrendingVideo

	err := chromedp.Run(ctx,
		chromedp.Navigate(chartURL),
		chromedp.WaitVisible(rowContainerElemSelector, chromedp.ByQuery),
		chromedp.Evaluate(chart.GetTrendingExtractScript(), &videos),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to extract chart data: %w", err)
	}

	if len(videos) == 0 {
		return nil, fmt.Errorf("no data extracted - \n")
	}

	return videos, nil
}

func (c *ChartCollector) CollectTopVideos() ([]chart.TopVideo, error) {
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), c.chromedpOptions...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(string, ...interface{}) {}))
	defer cancel()

	config := &chart.ChartConfig{Type: chart.TopVideos, Region: chart.RegionKR, TimeRange: chart.Daily}
	chartURL := config.GenChartURL()

	var videos []chart.TopVideo

	err := chromedp.Run(ctx,
		chromedp.Navigate(chartURL),
		chromedp.WaitVisible(rowContainerElemSelector, chromedp.ByQuery),
		chromedp.Evaluate(chart.GetTopVideosExtractScript(), &videos),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to extract chart data: %w", err)
	}

	if len(videos) == 0 {
		return nil, fmt.Errorf("no data extracted - \n")
	}

	return videos, nil
}

func (c *ChartCollector) CollectWeeklyTopSongs() ([]chart.WeeklyTopSong, error) {
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), c.chromedpOptions...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(string, ...interface{}) {}))
	defer cancel()

	config := &chart.ChartConfig{Type: chart.TopSongs, Region: chart.RegionKR, TimeRange: chart.Weekly}
	chartURL := config.GenChartURL()

	var songs []chart.WeeklyTopSong

	err := chromedp.Run(ctx,
		chromedp.Navigate(chartURL),
		chromedp.WaitVisible(rowContainerElemSelector, chromedp.ByQuery),
		chromedp.Evaluate(chart.GetWeeklyTopSongExtractScript(), &songs),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to extract chart data: %w", err)
	}

	if len(songs) == 0 {
		return nil, fmt.Errorf("no data extracted - \n")
	}

	return songs, nil
}

func (c *ChartCollector) CollectWeeklyTopArtists() ([]chart.WeeklyTopArtist, error) {
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), c.chromedpOptions...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(string, ...interface{}) {}))
	defer cancel()

	config := &chart.ChartConfig{Type: chart.TopArtists, Region: chart.RegionKR, TimeRange: chart.Weekly}
	chartURL := config.GenChartURL()

	var artists []chart.WeeklyTopArtist

	err := chromedp.Run(ctx,
		chromedp.Navigate(chartURL),
		chromedp.WaitVisible(rowContainerElemSelector, chromedp.ByQuery),
		chromedp.Evaluate(chart.GetWeeklyTopArtistExtractScript(), &artists),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to extract chart data: %w", err)
	}

	if len(artists) == 0 {
		return nil, fmt.Errorf("no data extracted - \n")
	}

	return artists, nil
}

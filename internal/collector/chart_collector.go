package collector

import (
	"channel-collector/internal/chart"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

const (
	chartExtractScript = `
		(function() {
			try {
				const videos = [];

				// Shadow DOM을 재귀적으로 탐색하는 함수
				function querySelectorAllDeep(selector, root = document) {
					const results = [];

					// 현재 레벨에서 검색
					const elements = root.querySelectorAll(selector);
					results.push(...elements);

					// 모든 요소의 Shadow Root 탐색
					const allElements = root.querySelectorAll('*');
					for (const element of allElements) {
						if (element.shadowRoot) {
							results.push(...querySelectorAllDeep(selector, element.shadowRoot));
						}
					}

					return results;
				}

				// Shadow DOM 내부에서 chart rows 찾기
				let chartRows = querySelectorAllDeep('div.ytmc-entry-row-container');

				// 만약 못 찾으면 다른 셀렉터 시도
				if (chartRows.length === 0) {
					chartRows = querySelectorAllDeep('ytmc-entry-row');
				}

				console.log('Found', chartRows.length, 'chart rows in Shadow DOM');

				chartRows.forEach((row, index) => {
					try {
						// Rank 추출
						const rankEl = row.querySelector('#rank');
						const rank = rankEl ? parseInt(rankEl.textContent.trim()) : index + 1;

						// Thumbnail & Video ID 추출
						const thumbnailEl = row.querySelector('img.video-thumbnail#thumbnail');
						let videoId = '';
						let thumbnailUrl = '';

						if (thumbnailEl) {
							thumbnailUrl = thumbnailEl.getAttribute('src') || '';
							// src에서 video ID 추출: https://i.ytimg.com/vi/VIDEO_ID/mqdefault.jpg
							const srcMatch = thumbnailUrl.match(/\/vi\/([^\/]+)\//);
							if (srcMatch) {
								videoId = srcMatch[1];
							}

							// endpoint 속성에서도 시도
							if (!videoId) {
								const endpoint = thumbnailEl.getAttribute('endpoint');
								if (endpoint) {
									try {
										const endpointObj = JSON.parse(endpoint);
										const url = endpointObj?.urlEndpoint?.url || '';
										const urlMatch = url.match(/[?&]v=([^&]+)/);
										if (urlMatch) {
											videoId = urlMatch[1];
										}
									} catch (e) {}
								}
							}
						}

						// Title 추출
						const titleEl = row.querySelector('div.title#entity-title');
						const title = titleEl ? titleEl.textContent.trim() : '';

						// Artist 추출
						const artistEl = row.querySelector('span.artistName.clickable');
						const artist = artistEl ? artistEl.textContent.trim() : '';

						// Release Date 추출 (첫 번째 metric)
						const metricEls = row.querySelectorAll('div.metric.content.center');
						let releaseDate = '';
						if (metricEls.length > 0) {
							releaseDate = metricEls[0].textContent.trim();
						}

						// Rank Trend 추출
						const rankTrendEl = row.querySelector('#rank-trend');
						let rankChange = 0;
						let isNew = false;

						if (rankTrendEl && !rankTrendEl.hasAttribute('hidden')) {
							const trendText = rankTrendEl.textContent.trim();
							if (trendText.includes('▲') || trendText.includes('▼')) {
								// rank-up or rank-down 클래스 확인
								if (rankTrendEl.classList.contains('rank-up')) {
									rankChange = 1; // 상승
								} else if (rankTrendEl.classList.contains('rank-down')) {
									rankChange = -1; // 하락
								}
							}
						}

						// Artist ID 추출 (endpoint에서)
						let artistId = '';
						if (artistEl) {
							const endpoint = artistEl.getAttribute('endpoint');
							if (endpoint) {
								try {
									const endpointObj = JSON.parse(endpoint);
									const query = endpointObj?.browseEndpoint?.query || '';
									const queryObj = JSON.parse(query);
									artistId = queryObj?.artistParamsId || '';
								} catch (e) {}
							}
						}

						videos.push({
							rank: rank,
							videoId: videoId,
							title: title,
							artist: artist,
							viewCount: 0,
							thumbnailUrl: thumbnailUrl,
							channelName: artist,
							channelId: artistId,
							rankChange: rankChange,
							isNew: isNew,
							publishedAt: releaseDate
						});
					} catch (e) {
						console.error('Error parsing row:', e);
					}
				});

				return videos;
			} catch (e) {
				console.error('Error extracting chart data:', e);
				return null;
			}
		})();
	`
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

func (c *ChartCollector) Collect(
	configs []*chart.ChartConfig,
	resultCh chan<- *chart.Chart,
	errCh chan<- error,
	wg *sync.WaitGroup,
) {
	defer close(resultCh)
	defer close(errCh)
	defer wg.Done()

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), c.chromedpOptions...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(string, ...interface{}) {}))
	defer cancel()

	for _, config := range configs {
		chartURL := config.GenChartURL()
		log.Printf("Starting to collect chart: %s\n", chartURL)

		chartData, err := c.collectTrendingChart(ctx, config)
		if err != nil {
			errCh <- fmt.Errorf("failed to collect chart %s: %w", chartURL, err)
			continue
		}

		if chartData != nil {
			resultCh <- chartData
			log.Printf("Successfully collected chart: %s with %d videos\n", chartURL, len(chartData.Videos))
		} else {
			errCh <- fmt.Errorf("no data found for chart: %s", chartURL)
		}
	}
}

func (c *ChartCollector) collectTrendingChart(ctx context.Context, config *chart.ChartConfig) (*chart.Chart, error) {
	chartURL := config.GenChartURL()

	var videos []struct {
		Rank         int    `json:"rank"`
		VideoID      string `json:"videoId"`
		Title        string `json:"title"`
		Artist       string `json:"artist"`
		ViewCount    int64  `json:"viewCount"`
		ThumbnailURL string `json:"thumbnailUrl"`
		ChannelName  string `json:"channelName"`
		ChannelID    string `json:"channelId"`
		RankChange   int    `json:"rankChange"`
		IsNew        bool   `json:"isNew"`
		PublishedAt  string `json:"publishedAt"`
	}

	err := chromedp.Run(ctx,
		chromedp.Navigate(chartURL),
		chromedp.WaitVisible("div.ytmc-entry-row-container", chromedp.ByQuery),
		chromedp.Evaluate(chartExtractScript, &videos),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to extract chart data: %w", err)
	}

	if len(videos) == 0 {
		return nil, fmt.Errorf("no data extracted - \n")
	}

	chartVideos := make([]chart.TrendingVideo, 0, len(videos))
	for _, v := range videos {
		chartVideos = append(chartVideos, chart.TrendingVideo{
			Rank:         v.Rank,
			VideoID:      v.VideoID,
			Title:        v.Title,
			Artist:       v.Artist,
			ViewCount:    v.ViewCount,
			ThumbnailURL: v.ThumbnailURL,
			ChannelName:  v.ChannelName,
			ChannelID:    v.ChannelID,
			RankChange:   v.RankChange,
			IsNew:        v.IsNew,
		})
	}

	return &chart.Chart{
		Type:        config.Type,
		Region:      config.Region,
		TimeRange:   config.TimeRange,
		CollectedAt: time.Now(),
		Videos:      chartVideos,
	}, nil
}

package chart

import "time"

type TopVideo struct {
	Rank         int       `json:"rank"`
	VideoID      string    `json:"videoId"`
	VideoURL     string    `json:"videoUrl"`
	Title        string    `json:"title"`
	Artist       string    `json:"artist"`
	ViewCount    int64     `json:"viewCount"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	ChannelName  string    `json:"channelName"`
	ChannelID    string    `json:"channelId"`
	RankChange   int       `json:"rankChange"` // Positive for up, negative for down, 0 for unchanged
	PublishedAt  time.Time `json:"publishedAt"`
	PeakPosition int       `json:"peakPosition"` // 최고 순위
	WeeksOnChart int       `json:"weeksOnChart"` // 차트 진입 주 수
}

func GetTopVideosExtractScript() string {
	return `
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

						// Title 및 Video URL 추출
						const titleEl = row.querySelector('div.title#entity-title');
						const title = titleEl ? titleEl.textContent.trim() : '';

						// Video URL 추출 (#entity-title의 endpoint 속성에서)
						let videoUrl = '';
						if (titleEl) {
							const endpoint = titleEl.getAttribute('endpoint');
							if (endpoint) {
								try {
									const endpointObj = JSON.parse(endpoint);
									videoUrl = endpointObj?.urlEndpoint?.url || '';

									// videoUrl에서 video ID도 추출 (thumbnail에서 추출 못했을 경우 대비)
									if (!videoId && videoUrl) {
										const urlMatch = videoUrl.match(/[?&]v=([^&]+)/);
										if (urlMatch) {
											videoId = urlMatch[1];
										}
									}
								} catch (e) {}
							}
						}

						// Artist 추출
						const artistEl = row.querySelector('span.artistName.clickable');
						const artist = artistEl ? artistEl.textContent.trim() : '';

						// Metrics 추출
						const metricEls = row.querySelectorAll('div.metric.content.center');
						let releaseDate = '';
						let peakPosition = 0;
						let weeksOnChart = 0;
						let viewCount = 0;

						// metricEls[0]: Release Date (hidden일 수 있음)
						if (metricEls.length > 0) {
							const dateStr = metricEls[0].textContent.trim();
							// "Sep 15, 2025"를 파싱 후 KST(+09:00) 기준으로 변환
							const date = new Date(dateStr + ' GMT+0900');
							if (!isNaN(date.getTime())) {
								releaseDate = date.toISOString();
							}
						}

						// metricEls[1]: Peak Position
						if (metricEls.length > 1) {
							peakPosition = parseInt(metricEls[1].textContent.trim()) || 0;
						}

						// metricEls[2]: Weeks on Chart
						if (metricEls.length > 2) {
							weeksOnChart = parseInt(metricEls[2].textContent.trim()) || 0;
						}

						// metricEls[3]: View Count (hidden일 수 있음)
						if (metricEls.length > 3) {
							const viewStr = metricEls[3].textContent.trim().replace(/,/g, '');
							viewCount = parseInt(viewStr) || 0;
						}

						// Rank Trend 추출
						const rankTrendEl = row.querySelector('#rank-trend');
						let rankChange = 0;

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
							// '●'는 변동 없음 (0)
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
							videoUrl: videoUrl,
							title: title,
							artist: artist,
							viewCount: viewCount,
							thumbnailUrl: thumbnailUrl,
							channelName: artist,
							channelId: artistId,
							rankChange: rankChange,
							publishedAt: releaseDate,
							peakPosition: peakPosition,
							weeksOnChart: weeksOnChart
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
}

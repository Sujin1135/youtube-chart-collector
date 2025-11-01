package chart

type WeeklyTopArtist struct {
	Rank         int    `json:"rank"`
	ArtistName   string `json:"artistName"`
	ArtistID     string `json:"artistId"`
	ThumbnailURL string `json:"thumbnailUrl"`
	ViewCount    int64  `json:"viewCount"`
	RankChange   int    `json:"rankChange"`   // Positive for up, negative for down, 0 for unchanged
	PeakPosition int    `json:"peakPosition"` // 최고 순위
	WeeksOnChart int    `json:"weeksOnChart"` // 차트 진입 주 수
}

func GetWeeklyTopArtistExtractScript() string {
	return `
		(function() {
			try {
				const artists = [];

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

						// Thumbnail 추출 (주간 인기 아티스트는 artists-thumbnail 클래스 사용)
						const thumbnailEl = row.querySelector('img.artists-thumbnail#thumbnail');
						let thumbnailUrl = '';

						if (thumbnailEl) {
							thumbnailUrl = thumbnailEl.getAttribute('src') || '';
						}

						// Artist Name 추출 (title 컨테이너 내부의 artistName)
						const titleContainer = row.querySelector('div.title-container');
						let artistName = '';
						let artistId = '';

						if (titleContainer) {
							// title 컨테이너 안의 첫 번째 artistName 요소
							const artistEl = titleContainer.querySelector('span.artistName.clickable');
							if (artistEl) {
								artistName = artistEl.textContent.trim();

								// Artist ID 추출 (endpoint에서)
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
						}

						// Metrics 추출
						const metricEls = row.querySelectorAll('div.metric.content.center');
						let peakPosition = 0;
						let weeksOnChart = 0;
						let viewCount = 0;

						// metricEls[0]: Release Date (hidden)
						// metricEls[1]: Peak Position
						if (metricEls.length > 1) {
							peakPosition = parseInt(metricEls[1].textContent.trim()) || 0;
						}

						// metricEls[2]: Weeks on Chart
						if (metricEls.length > 2) {
							weeksOnChart = parseInt(metricEls[2].textContent.trim()) || 0;
						}

						// metricEls[3]: View Count
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

						artists.push({
							rank: rank,
							artistName: artistName,
							artistId: artistId,
							thumbnailUrl: thumbnailUrl,
							viewCount: viewCount,
							rankChange: rankChange,
							peakPosition: peakPosition,
							weeksOnChart: weeksOnChart
						});
					} catch (e) {
						console.error('Error parsing row:', e);
					}
				});

				return artists;
			} catch (e) {
				console.error('Error extracting chart data:', e);
				return null;
			}
		})();
	`
}
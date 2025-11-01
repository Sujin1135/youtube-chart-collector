# YouTube Chart Collector

YouTube 인기차트 정보와 차트 데이터를 수집하는 Go 기반 웹 크롤러입니다.

## Requirements
|                | Main version  |
|----------------|---------------|
| Go             | 1.24 or More  |
| Chrome Browser | 134.x         |

## Features

- **YouTube Chart 카테고리별 비디오 수집**: 각 카테고리별 비디오 데이터
- **확장 가능한 구조**: 다양한 차트 타입과 지역을 지원

## Project Structure

```
channel-collector/
├── cmd/
│   └── main.go                    # 메인 실행 파일
├── internal/
│   ├── channel/
│   │   └── channel.go             # 채널 데이터 모델
│   ├── chart/
│   │   └── chart.go               # 차트 데이터 모델
└── └── collector/
        └── chart_collector.go     # 차트 수집 로직
```

## Installation

```shell
# Install dependencies
go mod download

# Build
go build -o collector cmd/main.go
```

## Chart Types

### ChartType
- `TrendingVideos`: 인기 급상승 동영상
- `TopVideos`: 인기 동영상
- `TopArtists`: 인기 아티스트
- `TopSongs`: 인기 곡

### Region
- `RegionKR`: 한국
- `RegionUS`: 미국
- `RegionGlobal`: 전세계

### TimeRange
- `RightNow`: 실시간
- `Daily`: 일간
- `Weekly`: 주간

## Extending

새로운 차트 타입을 추가하려면:

1. `internal/chart/chart.go`에 새로운 `ChartType` 상수 추가
2. 필요시 `chart_collector.go`의 JavaScript 스크립트 수정
3. 새로운 `ChartConfig` 생성하여 사용

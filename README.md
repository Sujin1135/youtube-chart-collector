# YouTube Chart Collector

YouTube 채널 정보와 차트 데이터를 수집하는 Go 기반 웹 크롤러입니다.

## Requirements
|                | Main version  |
|----------------|---------------|
| Go             | 1.24 or More  |
| Chrome Browser | 134.x         |

## Features

- **YouTube 채널 정보 수집**: 채널 메타데이터, 구독자 수, 조회수 등을 수집
- **YouTube Charts 수집**: 인기 급상승 뮤직비디오, 인기 동영상 등 차트 데이터 수집
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
│   └── collector/
│       ├── channel_collector.go   # 채널 수집 로직
│       └── chart_collector.go     # 차트 수집 로직
└── examples/
    └── chart_example.go           # 차트 수집 예제
```

## Installation

```shell
# Install dependencies
go mod download

# Build
go build -o collector cmd/main.go
```

## Usage

### CLI Commands

프로그램은 두 가지 모드를 지원합니다: `chart` (차트 수집)과 `channel` (채널 수집)

#### 1. YouTube Charts 수집

**기본 사용법 (단일 차트)**
```shell
# 한국 인기 급상승 뮤직비디오 (기본값)
go run cmd/main.go -mode=chart

# 미국 인기 급상승 뮤직비디오
go run cmd/main.go -mode=chart -region=us

# 한국 주간 인기 동영상
go run cmd/main.go -mode=chart -chart-type=TopVideos -time-range=Weekly

# 결과를 파일로 저장
go run cmd/main.go -mode=chart -output=result.json
```

**여러 차트 동시 수집**
```shell
# 사전 정의된 여러 차트를 한번에 수집
go run cmd/main.go -mode=chart -multi-chart
```

**사용 가능한 플래그:**
- `-mode`: 수집 모드 (`chart` 또는 `channel`, 기본값: `chart`)
- `-chart-type`: 차트 타입 (`TrendingVideos`, `TopVideos`, `TopArtists`, `TopSongs`, 기본값: `TrendingVideos`)
- `-region`: 지역 코드 (`kr`, `us`, `global`, 기본값: `kr`)
- `-time-range`: 시간 범위 (`RightNow`, `Daily`, `Weekly`, 기본값: `RightNow`)
- `-multi-chart`: 여러 차트 동시 수집 (boolean)
- `-output`: 결과 저장 파일 경로 (JSON 형식, 비어있으면 stdout 출력)

#### 2. YouTube 채널 수집

```shell
# 단일 채널 수집
go run cmd/main.go -mode=channel -channels=@channel_name

# 여러 채널 동시 수집 (쉼표로 구분)
go run cmd/main.go -mode=channel -channels=@channel1,@channel2,@channel3

# 결과를 파일로 저장
go run cmd/main.go -mode=channel -channels=@channel1 -output=channels.json
```

**사용 가능한 플래그:**
- `-mode`: 수집 모드 (`channel`)
- `-channels`: 쉼표로 구분된 YouTube 채널 핸들 (예: `@channel1,@channel2`)
- `-output`: 결과 저장 파일 경로 (JSON 형식)

### 예제 실행

```shell
# 차트 수집 예제 실행
go run examples/chart_example.go
```

### 프로그래밍 방식 사용

```go
package main

import (
    "channel-collector/internal/chart"
    "channel-collector/internal/collector"
    "sync"
)

func main() {
    // 차트 설정 생성
    configs := []*chart.ChartConfig{
        {
            Type:      chart.TrendingVideos,  // 차트 타입
            Region:    chart.RegionKR,        // 지역 (한국)
            TimeRange: chart.RightNow,        // 시간 범위
        },
    }

    // 결과 채널 생성
    resultCh := make(chan *chart.Chart, len(configs))
    errCh := make(chan error, len(configs))

    var wg sync.WaitGroup
    wg.Add(1)

    // Collector 실행
    chartCollector := collector.NewChartCollector()
    go chartCollector.Collect(configs, resultCh, errCh, &wg)

    // 결과 처리
    wg.Wait()
    for chartData := range resultCh {
        // 차트 데이터 처리
    }
}
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

## Flow

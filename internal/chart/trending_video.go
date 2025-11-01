package chart

import "time"

type TrendingVideo struct {
	Rank         int       `json:"rank"`
	VideoID      string    `json:"videoId"`
	Title        string    `json:"title"`
	Artist       string    `json:"artist"`
	ViewCount    int64     `json:"viewCount"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	ChannelName  string    `json:"channelName"`
	ChannelID    string    `json:"channelId"`
	RankChange   int       `json:"rankChange"` // Positive for up, negative for down
	IsNew        bool      `json:"isNew"`
	PublishedAt  time.Time `json:"publishedAt"`
}

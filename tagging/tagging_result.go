package tagging

import (
	"time"

	"github.com/AndreyMarkinPPC/media-tagging-go-client/media"
)

type TaggingOutput int

const (
	TaggingOutputTags TaggingOutput = iota
	TaggingOutputDescription
)

type Tag struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type Description struct {
	Text string `json:"text"`
}

type TaggingResult struct {
	Identifier     string
	Content        Description
	Type           media.MediaType
	ProcessedAt    time.Time
	TaggingDetails map[string]string
}

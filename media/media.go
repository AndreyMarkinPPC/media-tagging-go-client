package media

import (
	"errors"
)

var (
	ErrorMediaNotFound = errors.New("No media found")
)

type MediaType string

const (
	MediaTypeImage        MediaType = "Image"
	MediaTypeVideo        MediaType = "Video"
	MediaTypeYouTubeVideo MediaType = "Youtube"
)

type Media struct {
	Path string
	Type MediaType
	Name string
}

func NewMedia(path string, mediaType MediaType) Media {
	return Media{Path: path, Type: mediaType}
}

func FromPaths(paths []string, mediaType MediaType) ([]Media, error) {
	if len(paths) == 0 {
		return []Media{}, ErrorMediaNotFound
	}
	data := make([]Media, len(paths))
	for _, m := range paths {
		data = append(data, NewMedia(m, mediaType))
	}
	return data, nil
}

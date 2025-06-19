package tagging

import (
	"github.com/AndreyMarkinPPC/media-tagging-go-client/media"
)

type taggingOptions struct {
	nTags        int
	tags         []string
	customPrompt string
}

type MediaTaggingRequest struct {
	TaggerType     string
	MediaType      media.MediaType
	MediaPaths     []string
	TaggingOptions taggingOptions
}

type MediaTaggingResponse struct {
	Results []TaggingResult
}

func TagMedia(tagger Tagger, request MediaTaggingRequest) (MediaTaggingResponse, error) {
	data, err := media.FromPaths(request.MediaPaths, request.MediaType)
	if err != nil {
		return MediaTaggingResponse{}, err
	}
	res, err := tagger.Tag(data, request.TaggingOptions)
	if err != nil {
		return MediaTaggingResponse{}, err
	}
	return MediaTaggingResponse{Results: res}, nil
}

func DescribeMedia(tagger Tagger, request MediaTaggingRequest) (MediaTaggingResponse, error) {
	data, err := media.FromPaths(request.MediaPaths, request.MediaType)
	if err != nil {
		return MediaTaggingResponse{}, err
	}
	res, err := tagger.Describe(data, request.TaggingOptions)
	if err != nil {
		return MediaTaggingResponse{}, err
	}
	return MediaTaggingResponse{Results: res}, nil
}

package tagging

import "github.com/AndreyMarkinPPC/media-tagging-go-client/media"

type Tagger interface {
	Tag(media []media.Media, options taggingOptions) ([]TaggingResult, error)
	Describe(
		media []media.Media,
		options taggingOptions,
	) ([]TaggingResult, error)
}

package tagging

import (
	"context"
	"errors"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/AndreyMarkinPPC/media-tagging-go-client/media"
	"google.golang.org/genai"
)

const PROMPT string = "What this media is about"
const MODEL string = "gemini-2.5-flash"

type GeminiTagger struct {
	apiKey string
	model  string
}

func New(model, apiKey string) (GeminiTagger, error) {
	if apiKey == "" {
		apiKey = os.Getenv("GOOGLE_API_KEY")
	}
	if apiKey == "" {
		return GeminiTagger{}, errors.New("No ApiKey found.")
	}
	if model == "" {
		model = MODEL
	}
	return GeminiTagger{model: model, apiKey: apiKey}, nil
}

func (t GeminiTagger) Tag(
	media []media.Media,
	options taggingOptions,
) ([]TaggingResult, error) {
	ctx := context.Background()
	client, err := t.client(ctx)
	if err != nil {
		log.Fatal(err)
	}
	results, err := t.processMedia(
		ctx,
		client,
		media,
		"image/jpeg",
		PROMPT,
	)

	return results, err
}

func (t GeminiTagger) Describe(
	media []media.Media,
	options taggingOptions,
) ([]TaggingResult, error) {
	return []TaggingResult{}, nil
}

func (t GeminiTagger) client(ctx context.Context) (*genai.Client, error) {
	client, err := genai.NewClient(ctx,
		&genai.ClientConfig{
			APIKey:  t.apiKey,
			Backend: genai.BackendGeminiAPI,
		})

	return client, err
}

func (t GeminiTagger) processMedia(
	ctx context.Context,
	client *genai.Client,
	media []media.Media,
	mimeType string,
	prompt string,
) ([]TaggingResult, error) {
	consolidate := make([]TaggingResult, 0)
	resCh := make(chan TaggingResult)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	filesCh := make(chan string)
	wg := sync.WaitGroup{}

	go func() {
		defer close(filesCh)
		for _, f := range media {
			filesCh <- f.Path
		}
	}()

	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range filesCh {
				result, err := t.processMedium(
					ctx,
					client,
					f,
					mimeType,
					prompt,
				)
				if err != nil {
					errCh <- err
				}
				resCh <- result
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return nil, err
		case results := <-resCh:
			consolidate = append(consolidate, results)
		case <-doneCh:
			return consolidate, nil
		}
	}

}

func (t GeminiTagger) processMedium(
	ctx context.Context,
	client *genai.Client,
	mediaPath, mimeType string,
	prompt string,
) (TaggingResult, error) {
	file, err := client.Files.UploadFromPath(
		ctx, mediaPath,
		&genai.UploadFileConfig{MIMEType: mimeType},
	)
	if err != nil {
		return TaggingResult{}, err
	}
	if prompt == "" {
		prompt = PROMPT
	}

	parts := []*genai.Part{
		genai.NewPartFromURI(file.URI, file.MIMEType),
		genai.NewPartFromText("\n"),
		genai.NewPartFromText(prompt),
	}
	request := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(
		ctx, t.model, request, nil)
	if err != nil {
		return TaggingResult{}, err
	}
	return TaggingResult{
		Content: Description{Text: result.Candidates[0].Content.Parts[0].Text},
		Type:    media.MediaTypeImage,
	}, nil
}

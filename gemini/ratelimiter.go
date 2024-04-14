package gemini

import (
	"context"
	"gemini-webscrapping/models"
	"github.com/google/generative-ai-go/genai"
	"golang.org/x/time/rate"
	"time"
)

type RateLimitedExtractor struct {
	limiter *rate.Limiter
}

func newRateLimitedExtractor(interval time.Duration, burstLimit int) *RateLimitedExtractor {
	return &RateLimitedExtractor{
		limiter: rate.NewLimiter(rate.Every(interval), burstLimit),
	}
}

func (rle *RateLimitedExtractor) extractReviewsWithRateLimit(client *genai.Client, strHTML string, prompt string) (models.Reviews, error) {
	if err := rle.limiter.Wait(context.Background()); err != nil {
		return models.Reviews{}, err
	}

	reviews, err := extractReviews(client, strHTML, prompt)
	if err != nil {
		return models.Reviews{}, err
	}
	return reviews, nil
}

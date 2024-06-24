package gemini_review_scrapping

import (
	"context"
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

func (rle *RateLimitedExtractor) extractReviewWithRateLimit(client *genai.Client, strHTML string, prompt string) (string, error) {
	if err := rle.limiter.Wait(context.Background()); err != nil {
		return "", err
	}
	return extractReview(client, strHTML, prompt)
}

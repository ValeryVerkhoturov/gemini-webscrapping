package common

import (
	"gemini-webscrapping/bbb_org"
	"gemini-webscrapping/gemini"
	"gemini-webscrapping/models"
	"gemini-webscrapping/youtube_com"
	"strings"
)

func ScrapReviews(url string) (models.Reviews, error) {
	if strings.HasPrefix(url, "https://www.bbb.org/") && len(strings.Split(url, "/")) >= 9 {
		return bbb_org.ScrapBBB(url)
	} else if strings.HasPrefix(url, "https://www.youtube.com/") || strings.HasPrefix(url, "https://youtu.be/") {
		return youtube_com.ScrapYoutube(url)
	} else {
		return gemini.ScrapGemini(url)
	}
}

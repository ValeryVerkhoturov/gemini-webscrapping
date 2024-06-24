package common_review_scrapping

import (
	"gemini-webscrapping/bbb_org_review_scrapping"
	"gemini-webscrapping/gemini_review_scrapping"
	"gemini-webscrapping/youtube_com_review_scrapping"
	"regexp"
	"strings"
)

func ScrapReview(urlsString string) (string, []error) {
	urls := findURLs(urlsString)
	reviews, errors := make([]string, 0), make([]error, 0)
	var review string
	var err error
	for _, url := range urls {
		review, err = scrapReview(url)
		if err != nil {
			errors = append(errors, err)
		}
		reviews = append(reviews, review)
	}
	if len(reviews) > 1 {
		review, err = gemini_review_scrapping.SummarizeReviews(reviews...)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return review, errors
}

func findURLs(input string) []string {
	re := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:\/~+#-]*[\w@?^=%&/~+#-])`)
	return re.FindAllString(input, -1)
}

func scrapReview(url string) (string, error) {
	if strings.HasPrefix(url, "https://www.bbb.org/") && len(strings.Split(url, "/")) >= 9 {
		return bbb_org_review_scrapping.ScrapBBBReview(url)
	} else if strings.HasPrefix(url, "https://www.youtube.com/") || strings.HasPrefix(url, "https://youtu.be/") {
		return youtube_com_review_scrapping.ScrapYoutubeReview(url)
	} else {
		return gemini_review_scrapping.ScrapNewsSiteReview(url)
	}
}

package bbb_org

import (
	"encoding/json"
	"errors"
	"fmt"
	"gemini-webscrapping/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	CustomerReviewsUrl = "https://www.bbb.org/api/businessprofile/customerreviews"
)

type bBBReview struct {
	ReviewStarRating float32 `json:"reviewStarRating"`
	DisplayName      string  `json:"displayName"`
	Text             string  `json:"text"`
	Date             struct {
		Day   string `json:"day"`
		Month string `json:"month"`
		Year  string `json:"year"`
	} `json:"date"`
}

type response struct {
	Items      []bBBReview `json:"items"`
	BusinessId string      `json:"businessId"`
	BBBId      string      `json:"bbbId"`
	NumFound   int         `json:"numFound"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
	Sort       string      `json:"sort"`
}

func getReviews(businessId string, bbbId string, page int, pageSize int) (models.Reviews, int, error) {
	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("pageSize", strconv.Itoa(pageSize))
	params.Add("businessId", businessId)
	params.Add("bbbId", bbbId)
	params.Add("sort", "reviewDate desc, id desc")
	fullURL := fmt.Sprintf("%s?%s", CustomerReviewsUrl, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, 0, errors.New("Status code of getting BBB reviews " + strconv.Itoa(res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	var jsonResponse response
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, 0, err
	}

	reviews := models.Reviews{}
	for _, responseItem := range jsonResponse.Items {
		date := fmt.Sprintf("%s-%s-%s", responseItem.Date.Year, responseItem.Date.Month, responseItem.Date.Day)
		reviews = append(reviews, &models.Review{
			Name:    responseItem.DisplayName,
			Message: responseItem.Text,
			Mark:    &responseItem.ReviewStarRating,
			Date:    date,
		})
	}

	return reviews, jsonResponse.TotalPages, nil
}

func ScrapBBB(url string) (models.Reviews, error) {
	splitUrl := strings.Split(url, "/")
	splitUrl = strings.Split(splitUrl[8], "-")
	businessId := splitUrl[len(splitUrl)-1]
	bbbId := splitUrl[len(splitUrl)-2]

	reviews, totalPages, err := getReviews(businessId, bbbId, 1, 10)
	if err != nil {
		return nil, err
	}
	if totalPages > 1 {
		for page := 2; page <= totalPages; page++ {
			newReviews, _, err := getReviews(businessId, bbbId, page, 10)
			if err != nil {
				return nil, err
			}
			reviews = append(reviews, newReviews...)
		}
	}
	return reviews, nil
}

package youtube_com

import (
	"encoding/json"
	"errors"
	"fmt"
	"gemini-webscrapping/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	YoutubeCommentsListUrl = "https://www.googleapis.com/youtube/v3/commentThreads"
)

type YoutubeCommentsResponse struct {
	Items []struct {
		Snippet struct {
			TopLevelComment struct {
				Snippet struct {
					ChannelID             string `json:"channelId"`
					TextDisplay           string `json:"textDisplay"`
					TextOriginal          string `json:"textOriginal"`
					AuthorDisplayName     string `json:"authorDisplayName"`
					AuthorProfileImageUrl string `json:"authorProfileImageUrl"`
					AuthorChannelUrl      string `json:"authorChannelUrl"`
					LikeCount             int    `json:"likeCount"`
					PublishedAt           string `json:"publishedAt"` // "2023-11-09T23:48:13Z"
					UpdatedAt             string `json:"updatedAt"`   // "2023-11-09T23:48:13Z"
				} `json:"snippet"`
			} `json:"topLevelComment"`
		} `json:"snippet"`
	} `json:"items"`
}

func getComments(videoID string, maxResults int) (models.Reviews, error) {
	// https://developers.google.com/youtube/v3/docs/commentThreads/list
	params := url.Values{}
	params.Add("key", os.Getenv("YOUTUBE_API_KEY"))
	params.Add("textFormat", "plainText")
	params.Add("part", "snippet")
	params.Add("videoId", videoID)
	params.Add("maxResults", strconv.Itoa(maxResults))
	params.Add("order", "relevance")
	fullURL := fmt.Sprintf("%s?%s", YoutubeCommentsListUrl, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("Status code of getting reviews " + strconv.Itoa(res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var jsonResponse YoutubeCommentsResponse
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, err
	}

	reviews := models.Reviews{}
	for _, responseItem := range jsonResponse.Items {
		date := strings.Split(responseItem.Snippet.TopLevelComment.Snippet.PublishedAt, "T")[0]
		reviews = append(reviews, &models.Review{
			Name:    responseItem.Snippet.TopLevelComment.Snippet.AuthorDisplayName,
			Message: responseItem.Snippet.TopLevelComment.Snippet.TextDisplay,
			Date:    date,
		})
	}

	return reviews, nil
}

func ScrapYoutube(videoUrl string) (models.Reviews, error) {
	const maxResults = 50

	if strings.HasPrefix(videoUrl, "https://youtu.be/") {
		videoId := strings.Split(strings.Split(videoUrl, "/")[3], "?")[0]
		return getComments(videoId, maxResults)
	} else if strings.HasPrefix(videoUrl, "https://www.youtube.com/watch") {
		u, err := url.Parse(videoUrl)
		if err != nil {
			return nil, err
		}
		videoID := u.Query().Get("v")

		return getComments(videoID, maxResults)
	} else if strings.HasPrefix(videoUrl, "https://www.youtube.com/shorts/") {
		videoId := strings.Split(strings.Split(videoUrl, "/")[4], "?")[0]
		return getComments(videoId, maxResults)
	}
	return nil, errors.New("Invalid Youtube URL")
}

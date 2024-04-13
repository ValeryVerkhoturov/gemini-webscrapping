package common

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"gemini-webscrapping/bbb_org"
	"gemini-webscrapping/models"
	"gemini-webscrapping/youtube_com"
	"github.com/google/generative-ai-go/genai"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/transform"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	FindReviewPrompt = `
You will get html document with multiple reviews of business.
You need to extract reviewer name, his message and his mark and date of review.
If message is too big, it should be summarized.
Output should be only json array with reviews with properties "name" as string, "message" as string, "mark" as float, "date" as string with format "YYYY-MM-DD".
If you did not find any reviews, return empty array. Do not write anything else.
`
	FindNewsPrompt = `
You will get html document with news about business.
You need to extract author name, text of news and date of review.
If text is too big, it should be summarized. 
Output should be only json array with single news object with properties "name" as string, "message" as string, "date" as string with format "YYYY-MM-DD".
If you did not find news, return empty array. Do not write anything else.`
)

func getHTML(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Status Code: " + resp.Status)
	}
	defer resp.Body.Close()

	var reader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = resp.Body
	}

	contentType := resp.Header.Get("Content-Type")
	charset := "utf-8"
	if contentType != "" {
		parts := strings.Split(contentType, "charset=")
		if len(parts) > 1 {
			charset = strings.TrimSpace(parts[1])
		}
	}

	encoder := getEncoder(charset)
	if encoder == nil {
		return nil, fmt.Errorf("Unsupported charset: %s", charset)
	}
	reader = transform.NewReader(reader, encoder.NewDecoder())

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func getEncoder(charset string) encoding.Encoding {
	encoder, err := htmlindex.Get(charset)
	if err != nil {
		return nil
	}
	return encoder
}

func body(doc *html.Node) (*html.Node, error) {
	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if body != nil {
		return body, nil
	}
	return nil, errors.New("Missing <body> in the node tree")
}

func removeTagsAndAttributes(n *html.Node) {
	if n.Type == html.ElementNode {
		if n.Data == "svg" || n.Data == "img" || n.Data == "picture" ||
			n.Data == "style" || n.Data == "canvas" || n.Data == "figure" {
			n.Parent.RemoveChild(n)
			return
		}
		for i := len(n.Attr) - 1; i >= 0; i-- {
			n.Attr[i].Key = ""
			n.Attr[i].Val = ""
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		removeTagsAndAttributes(c)
	}
}

func renderNode(node *html.Node) (string, error) {
	var b bytes.Buffer
	err := html.Render(&b, node)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func extractReviews(client *genai.Client, strHTML string, prompt string) (models.Reviews, error) {
	ctx := context.Background()
	model := client.GenerativeModel("models/gemini-1.5-pro-latest")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt), genai.Text(strHTML))
	if err != nil {
		return models.Reviews{}, err
	}
	respString := strings.Replace(fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), "\n", "", -1)
	re := regexp.MustCompile(`\[.*?\]`)
	result := re.FindStringSubmatch(respString)

	var reviews models.Reviews = models.Reviews{}
	if len(result) > 0 {
		err := reviews.UnmarshalLLMText(result[0])
		if err != nil {
			return models.Reviews{}, err
		}
		return reviews, nil
	}
	return models.Reviews{}, errors.New("Wrong LLM output")
}

func ScrapReviews(url string) (models.Reviews, error) {
	if strings.HasPrefix(url, "https://www.bbb.org/") && len(strings.Split(url, "/")) >= 9 {
		return bbb_org.ScrapBBB(url)
	} else if strings.HasPrefix(url, "https://www.youtube.com/") || strings.HasPrefix(url, "https://youtu.be/") {
		return youtube_com.ScrapYoutube(url)
	} else {
		fullPage, err := getHTML(url)
		if err != nil {
			return nil, err
		}
		doc, err := html.Parse(bytes.NewReader(fullPage))
		if err != nil {
			return nil, err
		}
		bodyNode, err := body(doc)
		if err != nil {
			return nil, err
		}
		removeTagsAndAttributes(bodyNode)
		body, err := renderNode(bodyNode)
		if err != nil {
			return nil, err
		}

		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
		if err != nil {
			return nil, err
		}
		defer client.Close()

		var reviews = models.Reviews{}
		for _, prompt := range []string{FindReviewPrompt, FindNewsPrompt} {
			extractedReviews, err := extractReviews(client, body, prompt)
			if err != nil {
				return nil, err
			}
			reviews = append(reviews, extractedReviews...)
		}
		return reviews, nil
	}
}

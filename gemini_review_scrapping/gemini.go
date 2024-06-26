package gemini_review_scrapping

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
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
	"strings"
	"time"
)

const (
	FindReviewFromCommentsPrompt = `
You will get multiple reviews of business.
You need to summarize all opinions and write review of business.
If you did not write complex review, write nothing. 
Write plain text without markdown.
`
	FindReviewFromNewsPrompt = `
You will get html document with news about business.
You need to summarize news and write review of business.
If you did not write complex review, write nothing.
Write plain text without markdown.
`
	SummarizeReviewFromSomeReviewsPrompt = `
You will get some reviews about business.
You need to summarize them and write single review of business.
If you did not write complex review, write nothing.
Write plain text without markdown.
`
)

var (
	rateLimitedExtractor RateLimitedExtractor
)

func init() {
	rateLimitedExtractor = *newRateLimitedExtractor(10*time.Second, 2)
}

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

func getHTMLBody(doc *html.Node) (*html.Node, error) {
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

func extractReview(client *genai.Client, strHTML string, prompt string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	model := client.GenerativeModel("models/gemini-1.5-flash") // models/gemini-1.5-pro-latest

	resp, err := model.GenerateContent(ctx, genai.Text(prompt), genai.Text(strHTML))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}

func ScrapNewsSiteReview(url string) (string, error) {
	fullPage, err := getHTML(url)
	if err != nil {
		return "", err
	}
	doc, err := html.Parse(bytes.NewReader(fullPage))
	if err != nil {
		return "", err
	}
	bodyNode, err := getHTMLBody(doc)
	if err != nil {
		return "", err
	}
	removeTagsAndAttributes(bodyNode)
	body, err := renderNode(bodyNode)
	if err != nil {
		return "", err
	}

	return scrapGeminiInternal(body, FindReviewFromNewsPrompt)
}

func ScrapGeminiFromMessages(messages []string) (string, error) {
	return scrapGeminiInternal(strings.Join(messages, "\n*****\n"), FindReviewFromCommentsPrompt)
}

func SummarizeReviews(reviews ...string) (string, error) {
	if len(reviews) == 0 {
		return "", errors.New("No input reviews to summarize")
	}
	return scrapGeminiInternal(strings.Join(reviews, "\n*****\n"), SummarizeReviewFromSomeReviewsPrompt)
}

func scrapGeminiInternal(strHTML string, prompt string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return "", err
	}
	defer client.Close()

	review, err := rateLimitedExtractor.extractReviewWithRateLimit(client, strHTML, prompt)
	if err != nil {
		return "", err
	}
	if review == "" {
		return "", errors.New("No Gemini output")
	}

	return review, nil
}

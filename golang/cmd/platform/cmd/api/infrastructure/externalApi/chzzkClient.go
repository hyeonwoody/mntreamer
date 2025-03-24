package externalApi

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"golang.org/x/time/rate"
)

type ChzzkClient struct {
	clnt *http.Client
	lmtr *rate.Limiter
}

func NewChzzkClient() *ChzzkClient {

	err := playwright.Install()
	if err != nil {
		log.Fatalf("Could not install Playwright: %v", err)
	}

	return &ChzzkClient{
		clnt: &http.Client{},
		lmtr: rate.NewLimiter(rate.Limit(14), 1)}
}

func (chzz *ChzzkClient) getHtml(url string) (string, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("🛑error creating request: %w", err)
	}
	return chzz.readHtml(req)
}

func (chzz *ChzzkClient) readHtml(req *http.Request) (string, error) {
	for {
		if !chzz.lmtr.Allow() {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		resp, err := chzz.clnt.Do(req)
		if err != nil {
			return "", fmt.Errorf("🛑error sending request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("🛑error reading response body: %w", err)
		}
		bodyString := string(body)
		return bodyString, nil
	}
}

func (chzz *ChzzkClient) Do(request *http.Request) (*http.Response, error) {
	return chzz.clnt.Do(request)
}

func (chzz *ChzzkClient) GetDocument(html string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (chzz *ChzzkClient) Get(url string) ([]byte, error) {
	panic("not implemented")
}

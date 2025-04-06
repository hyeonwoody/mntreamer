package business

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mntreamer/platform/cmd/api/infrastructure/externalApi"
	mntreamerModel "mntreamer/shared/model"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

type ChzzkBusiness struct {
	clnt externalApi.IClient
}

func NewChzzkBusiness(clnt externalApi.IClient) *ChzzkBusiness {
	return &ChzzkBusiness{clnt: clnt}
}

func (chzz *ChzzkBusiness) GetPlatformName() string {
	return "chzzk"
}

func (chzz *ChzzkBusiness) GetChannelId(nickname string) (string, error) {
	url := chzz.makeQueryUrl(nickname)
	html := renderQueryHtml(url)
	doc, err := chzz.clnt.GetDocument(html)
	if err != nil {
		return "", err
	}
	return chzz.extractChannelId(doc)
}

func (chzz *ChzzkBusiness) GetChannelName(nickname string) (string, error) {
	url := chzz.makeQueryUrl(nickname)
	html := renderQueryHtml(url)
	doc, err := chzz.clnt.GetDocument(html)
	if err != nil {
		return "", err
	}
	return chzz.extractChannelName(doc)
}

func renderQueryHtml(url string) string {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true), // Run headlessly
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	// Create a new browser page
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	_, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		log.Fatalf("could not visit page: %v", err)
	}

	locator := page.Locator("[class^='channel_item_channel']").First()
	err = locator.WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(10000),
	})
	if err != nil {
		log.Fatalf("could not wait for content: %v", err)
	}
	htmlContent, err := page.Content()
	if err != nil {
		return ""
	}
	return htmlContent
}

func (chzz *ChzzkBusiness) extractChannelName(doc *goquery.Document) (string, error) {
	channelElement := doc.Find("strong[class^='channel_item_channel']").First()
	if channelElement.Length() == 0 {
		return "", fmt.Errorf("channel element not found")
	}
	name := channelElement.Find("span[class^='name_text']").Text()
	if name == "" {
		return "", fmt.Errorf("channel name not found")
	}
	return name, nil
}
func (chzz *ChzzkBusiness) extractChannelId(doc *goquery.Document) (string, error) {
	selection := doc.Find("a[class^='channel_item_wrapper']")
	href, exist := selection.Attr("href")
	if exist {
		return strings.TrimPrefix(href, "/"), nil
	}
	return "", fmt.Errorf("channel Id not found")
}

func (chzz *ChzzkBusiness) makeQueryUrl(nickname string) string {
	return "https://chzzk.naver.com/search?query=" + nickname
}

func (chzz *ChzzkBusiness) GetMediaDetail(streamer *mntreamerModel.Streamer) (*mntreamerModel.Media, error) {
	apiUrl := fmt.Sprintf("http://api.chzzk.naver.com/service/v3/channels/%s/live-detail", streamer.ChannelId)
	liveDetail, err := chzz.extractLiveDetail(apiUrl)
	if err != nil {
		return nil, err
	}
	livePlayback, err := chzz.parseLivePlaybackJson(liveDetail.LivePlaybackJson)
	if err != nil {
		return nil, err
	}

	m3u8Url, err := chzz.getM3u8PlaylistUrl(livePlayback, err)
	if err != nil {
		return nil, err
	}

	mediaUrl, err := chzz.parseM3u8ChunklistUrl(m3u8Url)
	if err != nil {
		return nil, err
	}

	return mntreamerModel.NewMedia(liveDetail.LiveTitle, mediaUrl, liveDetail.LiveImageUrl, chzz.GetPlatformName()), nil
}

func (*ChzzkBusiness) getM3u8PlaylistUrl(livePlayback map[string]interface{}, err error) (string, error) {
	mediaList, ok := livePlayback["media"].([]interface{})
	if !ok || len(mediaList) == 0 {
		return "", fmt.Errorf("🛑media list is missing or empty: %w", err)
	}

	mediaItem, ok := mediaList[0].(map[string]interface{})
	if !ok {
		fmt.Println("Error: media item is not a valid map")
		return "", fmt.Errorf("🛑media item is not a valid map: %w", err)
	}

	m3u8Url, ok := mediaItem["path"].(string)
	if !ok {
		fmt.Println("Error: path key not found or not a string")
		return "", fmt.Errorf("🛑path key not found or not a string: %w", err)
	}
	return m3u8Url, nil
}

func (chzz *ChzzkBusiness) parseM3u8ChunklistUrl(m3u8Url string) (string, error) {
	resp, err := http.Get(m3u8Url)
	if err != nil {
		return "", fmt.Errorf("🛑failed to fetch playlist.m3u8 %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("🛑failed to read response body: %w", err)
	}

	playlistContent := string(body)

	re := regexp.MustCompile(`(?m)#EXT-X-STREAM-INF:.*RESOLUTION=852x480.*\n(.*chunklist\.m3u8)`)
	matches := re.FindStringSubmatch(playlistContent)

	if len(matches) < 2 {
		return "", fmt.Errorf("🛑480p stream URL not found in playlist.m3u8: %w", err)
	}
	relativeUrl := matches[1]
	index := strings.Index(m3u8Url, "hls_playlist.m3u8")
	baseUrl := m3u8Url[:index]
	lastSlashIndex := strings.LastIndex(baseUrl, "/")
	baseUrl = baseUrl[:lastSlashIndex+1]
	fullChunklistUrl := baseUrl + relativeUrl
	return fullChunklistUrl, nil
}

func (chzz *ChzzkBusiness) extractLiveDetail(apiUrl string) (*mntreamerModel.LiveDetail, error) {
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Println(string(reqDump))
	res, err := chzz.clnt.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code error: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var response struct {
		Content mntreamerModel.LiveDetail `json:"content"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	response.Content.LiveImageUrl = strings.ReplaceAll(response.Content.LiveImageUrl, "{type}", "480")

	return &response.Content, nil
}

func (chzz *ChzzkBusiness) parseLivePlaybackJson(livePlaybackJson string) (map[string]interface{}, error) {
	var livePlayback map[string]interface{}
	err := json.Unmarshal([]byte(livePlaybackJson), &livePlayback)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil, err
	}
	return livePlayback, nil
}

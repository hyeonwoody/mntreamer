package externalApi

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type IClient interface {
	getHtml(url string) (string, error)
	GetDocument(html string) (*goquery.Document, error)
	readHtml(req *http.Request) (string, error)
	Do(request *http.Request) (*http.Response, error)
}

type ClientStrategy struct {
	clientMap map[uint16]IClient
}

func NewClientStrategy(clientMap map[uint16]IClient) *ClientStrategy {
	return &ClientStrategy{clientMap: clientMap}
}

func (strat *ClientStrategy) GetClient(platformId uint16) IClient {
	return strat.clientMap[platformId]
}

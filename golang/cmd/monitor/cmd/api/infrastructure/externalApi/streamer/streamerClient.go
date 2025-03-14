package client

type StreamerClient struct {
}

func NewStreamerClient() IStreamerClient {
	return &StreamerClient{}
}

func (s *StreamerClient) 
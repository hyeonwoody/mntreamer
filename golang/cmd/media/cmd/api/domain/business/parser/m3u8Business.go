package business

import (
	"fmt"
	"io"
	"mntreamer/media/cmd/model"
	"strconv"
	"strings"
)

type M3u8Business struct {
}

func NewM3u8Business() *M3u8Business {
	return &M3u8Business{}
}

func (m *M3u8Business) Encode(mediaPlayList interface{}) model.IBuffer {
	mpl, ok := mediaPlayList.(*model.MediaPlaylist)
	if !ok {
		return nil
	}
	return mpl.Encode()
}

func (m *M3u8Business) Decode(reader io.Reader) (interface{}, error) {
	buf := model.NewM3u8Buffer(20*1024, 4*1024)
	buf.ReadFrom(reader)
	return m.decode(buf)
}

func (m *M3u8Business) decode(buf *model.M3u8Buffer) (*model.MediaPlaylist, error) {
	mediaPlaylist := model.NewMediaPlaylist(0, 30)
	decoding := &model.Decoding{}
	eof := false
	for !eof {
		line, err := buf.ReadLine()
		if err == io.EOF {
			eof = true
		} else if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if len(line) < 1 || line == "\r" {
			continue
		}
		err = m.decodeLine(mediaPlaylist, line, decoding)
		if err != nil {
			return nil, err
		}
	}
	return mediaPlaylist, nil
}

func (m *M3u8Business) decodeLine(mediaPlaylist *model.MediaPlaylist, line string, decoding *model.Decoding) error {
	var err error
	switch {

	case !decoding.INF && strings.HasPrefix(line, "#EXTINF:"):
		decoding.INF = true
		separatorIdx := strings.Index(line, ",")
		duration := line[8:separatorIdx]
		if len(duration) > 0 {
			decoding.Duration, err = strconv.ParseFloat(duration, 64)
			if err != nil {
				return fmt.Errorf("🛑failed to parse duration: %w", err)
			}
		}
	case !strings.HasPrefix(line, "#"):
		if decoding.INF {
			err := mediaPlaylist.AppendChunk(line, decoding.Duration)
			if err == model.PLAYLISTFULL {
				mediaPlaylist.Expand()
				err = mediaPlaylist.AppendChunk(line, decoding.Duration)
			}
			if err != nil {
				return err
			}
			if decoding.DISCONTINUITY {
				decoding.DISCONTINUITY = false
				mediaPlaylist.SetDiscontinuity(true)
			}
			decoding.INF = false
		}
	case !decoding.DISCONTINUITY && strings.HasPrefix(line, "#EXT-X-DISCONTINUITY"):
		decoding.DISCONTINUITY = true
	case line == "#EXTM3U":
		decoding.M3U = true
	case strings.HasPrefix(line, "#EXT-X-VERSION:"):
		if _, err = fmt.Sscanf(line, "#EXT-X-VERSION:%d", &mediaPlaylist.Version); err != nil {
			return fmt.Errorf("🛑failed to parse version: %w", err)
		}
	case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"):
		if _, err = fmt.Sscanf(line, "#EXT-X-MEDIA-SEQUENCE:%d", &mediaPlaylist.SeqNo); err != nil {
			return fmt.Errorf("🛑failed to parse media sequence number: %w", err)
		}
	case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
		if _, err = fmt.Sscanf(line, "#EXT-X-TARGETDURATION:%f", &mediaPlaylist.TargetDuration); err != nil {
			return fmt.Errorf("🛑failed to parse target duration: %w", err)
		}
	case line == "#EXT-X-ENDLIST":
		mediaPlaylist.EndList = true
		//ignore #EXT-X-ALLOW-CACHE:YES
	}

	return err
}

package model

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// #EXTM3U is the format identifier tag
// #EXT-X-TARGETDURATION is a decimal-integer, target duration in seconds
type MediaPlaylist struct {
	TargetDuration   float64
	SeqNo            uint64 // EXT-X-MEDIA-SEQUENCE is the Media Sequence Number of the first Media Segment that appears in a Playlist file
	segments         []*MediaSegment
	Iframe           bool   // EXT-X-I-FRAMES-ONLY
	DiscontinuitySeq uint64 // EXT-X-DISCONTINUITY-SEQUENCE
	StartTime        float64
	StartTimePrecise bool
	EndList          bool
	keyformat        int
	winSize          uint // max number of segments displayed in an encoded playlist; need set to zero for VOD playlists
	capacity         uint // total Capacity of slice used for the playlist
	head             uint // head of FIFO, we add segments to head
	tail             uint // Tail of FIFO, we remove segments from Tail
	count            uint // number of segments added to the playlist
	buf              IBuffer
	Version          uint8 // EXT-X-VERSION is the compatibility version of the Playlist file
}

func (mpl *MediaPlaylist) lastSegmentIndex() uint {
	if mpl.tail == 0 {
		return mpl.capacity - 1
	}
	return mpl.tail - 1
}

func (mpl *MediaPlaylist) SetDiscontinuity(discontinuity bool) {
	if mpl.count == 0 {
	}
	mpl.segments[mpl.lastSegmentIndex()].Discontinuity = discontinuity
}

func (mpl *MediaPlaylist) SetDiscontinuityWithIndex(idx uint, discontinuity bool) {
	if idx == 0 {
		return
	}
	seg, _ := mpl.GetSegment(idx)
	if seg != nil {
		seg.Discontinuity = true
	}

}

func NewMediaPlaylist(winSize uint, capacity uint) *MediaPlaylist {
	return &MediaPlaylist{
		Version:  3,
		capacity: capacity,
		segments: make([]*MediaSegment, capacity),
	}
}

func (mpl *MediaPlaylist) Expand() {
	mpl.segments = append(mpl.segments, make([]*MediaSegment, mpl.count)...)
	mpl.capacity = uint(len(mpl.segments))
	mpl.tail = mpl.count
}

func (mpl *MediaPlaylist) Count() uint {
	return mpl.count
}

func (mpl *MediaPlaylist) GetSegment(index uint) (*MediaSegment, error) {
	actualIdx := (index) % mpl.capacity
	segment := mpl.segments[actualIdx]
	return segment, nil
}

func (mpl *MediaPlaylist) PullSegment(index uint) (*MediaSegment, error) {
	actualIdx := (mpl.head + index) % mpl.capacity
	segment := mpl.segments[actualIdx]
	for i := actualIdx; i != mpl.tail; i = (i + 1) % mpl.capacity {
		nextIndex := (i + 1) % mpl.capacity
		mpl.segments[i] = mpl.segments[nextIndex]
	}
	mpl.tail = (mpl.tail - 1) % mpl.capacity
	mpl.count--
	return segment, nil
}

func (mpl *MediaPlaylist) AppendChunk(uri string, duration float64) error {
	segment := &MediaSegment{}
	segment.Uri = uri
	segment.Duration = duration
	return mpl.appendSegment(segment)
}

var PLAYLISTFULL = errors.New("🛑media playlist is full")

func (mpl *MediaPlaylist) appendSegment(segment *MediaSegment) error {
	if mpl.head == mpl.tail && mpl.count > 0 {
		return PLAYLISTFULL
	}
	segment.SeqId = mpl.SeqNo
	if mpl.count > 0 {
		segment.SeqId = mpl.segments[(mpl.capacity+mpl.tail-1)%mpl.capacity].SeqId + 1
	}
	mpl.segments[mpl.tail] = segment
	mpl.tail = (mpl.tail + 1) % mpl.capacity
	mpl.count++
	return nil
}

func (mpl *MediaPlaylist) Encode() IBuffer {
	sb := &strings.Builder{}
	sb.WriteString("#EXTM3U\n")
	sb.WriteString("#EXT-X-VERSION:")
	sb.WriteString(strconv.FormatUint(uint64(mpl.Version), 10))
	sb.WriteString("\n")
	sb.WriteString("#EXT-X-MEDIA-SEQUENCE:")
	sb.WriteString(strconv.FormatUint(mpl.SeqNo, 10))
	sb.WriteString("\n")
	sb.WriteString("#EXT-X-TARGETDURATION:")
	sb.WriteString(strconv.FormatInt(int64(math.Ceil(mpl.TargetDuration)), 10))
	sb.WriteString("\n")

	head := mpl.head
	count := mpl.count
	durationDp := make((map[float64]string))
	for i := uint(0); (i < mpl.winSize || mpl.winSize == 0) && count > 0; count-- {
		segment := mpl.segments[head]
		head = (head + 1) % mpl.capacity
		if segment == nil {
			continue
		}
		if mpl.winSize > 0 {
			i++
		}
		duration, ok := durationDp[segment.Duration]
		if segment.Discontinuity {
			sb.WriteString("#EXT-X-DISCONTINUITY\n")
		}
		if !ok {
			duration = strconv.FormatFloat(segment.Duration, 'f', 6, 32)
			durationDp[segment.Duration] = duration
		}
		sb.WriteString("#EXTINF:")
		sb.WriteString(duration)
		sb.WriteString(",\n")
		sb.WriteString(segment.Uri)
		sb.WriteString("\n")
	}
	if mpl.EndList {
		sb.WriteString("#EXT-X-ENDLIST\n")
	}
	mpl.buf = NewM3u8Buffer(sb.Len(), 0)
	mpl.buf.Append(sb)
	return mpl.buf
}

type Decoding struct {
	M3U           bool
	INF           bool
	DISCONTINUITY bool
	SCTE          bool
	Duration      float64
}

type MediaSegment struct {
	SeqId         uint64
	Uri           string
	Duration      float64 // the first parameter of EXTINF; floating-point duration values
	Discontinuity bool    // EXT-X-DISCONTINUITY indicates an encoding discontinuity (file format, tracks, timestamp, encoding parameters, encoding sequence)
	SCTE          *SCTE   // SCTE-35 used for Ad signaling in HLS
}

func (ms *MediaSegment) SetDiscontinuity(parm bool) {
	ms.Discontinuity = parm
}

type MediaType uint

const (
	// use 0 for not defined type
	EVENT MediaType = iota + 1
	VOD
)

type SCTE struct {
	CueType SCTE35CueType
	Cue     string
	ID      string
	Time    float64
	Elapsed float64
}

type SCTE35CueType uint

const (
	SCTE35Cue_OUT SCTE35CueType = iota // a splice out point (start of break)
	SCTE35Cue_MID                      // a segment between start and end cue points
	SCTE35Cue_IN                       // a splice in point (end of break)
)

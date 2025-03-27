package model

import (
	"bytes"
	"time"
)

// #EXTM3U is the format identifier tag
// #EXT-X-TARGETDURATION is a decimal-integer, target duration in seconds
type MediaPlaylist struct {
	TargetDuration   float64
	SeqNo            uint64 // EXT-X-MEDIA-SEQUENCE is the Media Sequence Number of the first Media Segment that appears in a Playlist file
	Segments         []*MediaSegment
	Iframe           bool   // EXT-X-I-FRAMES-ONLY
	DiscontinuitySeq uint64 // EXT-X-DISCONTINUITY-SEQUENCE
	StartTime        float64
	StartTimePrecise bool
	keyformat        int
	winsize          uint // max number of segments displayed in an encoded playlist; need set to zero for VOD playlists
	capacity         uint // total capacity of slice used for the playlist
	head             uint // head of FIFO, we add segments to head
	tail             uint // tail of FIFO, we remove segments from tail
	count            uint // number of segments added to the playlist
	buf              bytes.Buffer
	ver              uint8 // EXT-X-VERSION is the compatibility version of the Playlist file
	Custom           map[string]CustomTag
	customDecoders   []CustomDecoder
}

type MediaSegment struct {
	Uri      string
	Duration float64 // the first parameter of EXTINF; floating-point duration values
	Discontinuity   bool      // EXT-X-DISCONTINUITY indicates an encoding discontinuity (file format, tracks, timestamp, encoding parameters, encoding sequence)
	SCTE            *SCTE     // SCTE-35 used for Ad signaling in HLS
	ProgramDateTime time.Time // EXT-X-PROGRAM-DATE-TIME tag associates the first sample of a media segment with an absolute date and/or time
	Custom          map[string]CustomTag
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

type CustomTag interface {
	// TagName should return the full indentifier including the leading '#' as well as the
	// trailing ':' if the tag also contains a value or attribute list
	TagName() string
	// Encode should return the complete tag string as a *bytes.Buffer.
	// Return nil to not write anything to the m3u8.
	Encode() *bytes.Buffer
	// String should return the encoded tag as a string.
	String() string
}

type CustomDecoder interface {
	// TagName should return the full indentifier including the leading '#' as well as the
	// trailing ':' if the tag also contains a value or attribute list
	TagName() string
	// Decode parses a line from the playlist and returns the CustomTag representation
	Decode(line string) (CustomTag, error)
	// SegmentTag should return true if this CustomDecoder should apply per segment.
	// Should returns false if it a MediaPlaylist header tag.
	// This value is ignored for MasterPlaylists.
	SegmentTag() bool
}

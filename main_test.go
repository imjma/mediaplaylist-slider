package main

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/grafov/m3u8"
)

func TestMediaHandler(t *testing.T) {
	f, _ := os.Open("fixtures/media-vod.m3u8")
	p, _ := m3u8.NewMediaPlaylist(4, 4)
	_ = p.DecodeFrom(bufio.NewReader(f), true)

	server := httptest.NewServer(mediaHandler(p))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}

func TestNewSegmentCache(t *testing.T) {
	p := testNewMediaPlaylistStream()
	newSegmentCache(p.Segments)

	for i, seg := range p.Segments {
		if *seg != *segmentsCache[i] {
			t.Errorf("exp: %+v\ngot: %+v", segmentsCache[i], seg)
		}
	}
}

func testslide(t *testing.T) {
	p := testNewMediaPlaylistStream()
	newSegmentCache(p.Segments)
	slide(p, segmentsCache[0])
	if *p.Segments[len(p.Segments)-1] != *segmentsCache[0] {
		t.Fatal("slide: append failed")
	}
}

func testNewMediaPlaylistStream() *m3u8.MediaPlaylist {
	f, _ := os.Open("fixtures/media-stream.m3u8")
	p, _ := m3u8.NewMediaPlaylist(4, 4)
	_ = p.DecodeFrom(bufio.NewReader(f), true)
	return p
}

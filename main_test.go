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

	if resp.Header.Get("Content-Type") != "application/vnd.apple.mpegurl" {
		t.Fatalf("Received wrong Content-Type: %s\n", resp.Header.Get("Content-Type"))
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

func TestSlideOK(t *testing.T) {
	p := testNewMediaPlaylistStream()
	err := slide(p, p.Segments[0])
	if err != nil {
		t.Fatal(err)
	}
}

func testNewMediaPlaylistStream() *m3u8.MediaPlaylist {
	f, _ := os.Open("fixtures/media-stream.m3u8")
	p, _ := m3u8.NewMediaPlaylist(4, 4)
	_ = p.DecodeFrom(bufio.NewReader(f), true)
	return p
}

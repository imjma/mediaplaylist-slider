package main

import (
	"bufio"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/grafov/m3u8"
)

func TestMediaHandler(t *testing.T) {
	p, err := decodeMediaPlaylist("fixtures/media-vod.m3u8")
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(mediaHandler(p))
	defer server.Close()

	resp, err := http.NewRecorder("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatal("Received non-200 response: %d\n", resp.StatusCode)
	}
}

func decodeMediaPlaylist(file string) (*m3u8.MediaPlaylist, error) {
	mediapl := new(m3u8.MediaPlaylist)
	f, err := os.Open(file)
	if err != nil {
		return mediapl, err
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		return mediapl, err
	}
	if listType == m3u8.MEDIA {
		return p.(*m3u8.MediaPlaylist), nil
	}
	return mediapl, errors.New("File is not media playlist")
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/grafov/m3u8"
)

var segmentsCache []*m3u8.MediaSegment

func main() {
	m3u8File := "media.m3u8"
	f, err := os.Open(m3u8File)
	if err != nil {
		panic(err)
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		http.HandleFunc("/", mediaHandler(mediapl))
		http.ListenAndServe(":9080", nil)
	case m3u8.MASTER:
		masterpl := p.(*m3u8.MasterPlaylist)
		fmt.Printf("%+v\n", masterpl)
	}
}

func mediaHandler(p *m3u8.MediaPlaylist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !p.Closed && segmentsCache == nil {
			for _, seg := range p.Segments {
				if seg == nil {
					continue
				}
				segmentsCache = append(segmentsCache, seg)
			}
			segmentsCache[0].Discontinuity = true
			go func(mediapl *m3u8.MediaPlaylist) {
				c := time.Tick(3 * time.Second)
				i := 0
				for _ = range c {
					seg := segmentsCache[i]
					p.Remove()
					p.AppendSegment(seg)
					p.ResetCache()
					i++
					if i >= len(segmentsCache) {
						i = 0
					}
				}
			}(p)
		}

		w.Header()["Content-Type"] = []string{"application/vnd.apple.mpegurl"}
		w.Header()["Access-Control-Allow-Origin"] = []string{"*"}

		io.WriteString(w, fmt.Sprintf("%+v", p))
	}
}

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/dustin/go-humanize"
)

type DownloadTorrentResult struct {
	path  string
	Error error
}

func downloadTorrent(pathName string) <-chan DownloadTorrentResult {
	ch := make(chan DownloadTorrentResult)
	go func() {
		close(ch)
	}()
	return ch
}

func test_download() {
	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = "./data"
	clientConfig.DisableWebtorrent = true
	clientConfig.NoDefaultPortForwarding = true
	// clientConfig.NoDHT = true
	// // ClientTrackerConfig: torrent.ClientTrackerConfig{DisableTrackers: true},
	// // DefaultStorage:          storage.NewSqlitePieceCompletion(),
	// Seed: true,
	// // NoDefaultPortForwarding: true,  // TODO wtf?
	// DisableUTP: false, // going to choose UTP because it should have a lower priority than tcp connections, need some speed tests
	// // DisableTCP:              true,
	// // DisableIPv6:             true,
	// // HTTPUserAgent:     "reconbot/swarm-download",
	// // 		HttpRequestDirector func(*http.Request) error, // custom for signed in stuff

	c, err := torrent.NewClient(clientConfig)
	if err != nil {
		log.Fatalf("Failed to create torrent client: %v", err)
	}
	defer c.Close()
	t, _ := c.AddTorrentFromFile("./download/ubuntu-24.10-live-server-amd64.iso.torrent")
	if t.Info() != nil {
		log.Print("info!")
	} else {
		log.Print("getting ", t.InfoHash())
		<-t.GotInfo()
		log.Print("got info!")
	}
	t.VerifyData()
	t.DownloadAll() // mark all torrents for download
	torrentStats(t, false)
	c.WaitAll()
	log.Print("ermahgerd, torrent downloaded")
}

// todo stop the loop when the torrent finishes
func torrentStats(t *torrent.Torrent, pieceStates bool) {
	go func() {
		start := time.Now()
		if t.Info() == nil {
			fmt.Printf("%v: getting torrent info for %q\n", time.Since(start), t.Name())
			<-t.GotInfo()
		}
		lastStats := t.Stats()
		var lastLine string
		interval := 3 * time.Second
		for range time.Tick(interval) {
			var completedPieces, partialPieces int
			psrs := t.PieceStateRuns()
			for _, r := range psrs {
				if r.Complete {
					completedPieces += r.Length
				}
				if r.Partial {
					partialPieces += r.Length
				}
			}
			stats := t.Stats()
			byteRate := int64(time.Second)
			byteRate *= stats.BytesReadUsefulData.Int64() - lastStats.BytesReadUsefulData.Int64()
			byteRate /= int64(interval)
			line := fmt.Sprintf(
				"%v: downloading %q: %s/%s, %d/%d pieces completed (%d partial): %v/s from %d of %d peers \"%s\"\n",
				time.Since(start),
				t.Name(),
				humanize.Bytes(uint64(t.BytesCompleted())),
				humanize.Bytes(uint64(t.Length())),
				completedPieces,
				t.NumPieces(),
				partialPieces,
				humanize.Bytes(uint64(byteRate)),
				t.Stats().ActivePeers,
				t.Stats().TotalPeers,
				t.Info().Source,
			)
			if line != lastLine {
				lastLine = line
				os.Stdout.WriteString(line)
			}
			if pieceStates {
				fmt.Println(psrs)
			}
			lastStats = stats
		}
	}()
}

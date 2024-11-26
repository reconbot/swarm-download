package main

import (
	"errors"
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
		// download the torrent file and stash it in
		client, err := createClient()
		defer client.Close()
		if err != nil {
			wrappedError := errors.New(fmt.Sprintf("Failed to create torrent client: %v", err))
			ch <- DownloadTorrentResult{Error: wrappedError}
			return
		}
	}()
	return ch
}

func foo() {
	client, err := createClient()
	if err != nil {
		log.Fatalf("Failed to create torrent client: %v", err)
	}
	defer client.Close()
	t, _ := client.AddTorrentFromFile("./download/ubuntu-24.10-live-server-amd64.iso.torrent")
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
	client.WaitAll()
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

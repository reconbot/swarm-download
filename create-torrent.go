package main

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

func fileReadable(pathName string) (bool, error) {
	stat, err := os.Stat(pathName)
	if err == nil {
		if stat.IsDir() {
			return false, nil
		}
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type CreateTorrentResult struct {
	path  string
	Error error
}

func createTorrent(pathName string) <-chan CreateTorrentResult {
	ch := make(chan CreateTorrentResult)
	go func() {
		defer close(ch)
		readable, err := fileReadable(pathName)
		if err != nil {
			println(err)
			ch <- CreateTorrentResult{Error: err}
			return
		}
		if readable == false {
			err := errors.New("Cannot read: " + pathName)
			println(err)
			ch <- CreateTorrentResult{Error: err}
			return
		}
		println("file readable ", pathName)
		torrentPath := pathName + ".torrent"

		// https://web.archive.org/web/20220926063641/https://wiki.theory.org/index.php/BitTorrentSpecification
		metadata := metainfo.MetaInfo{
			AnnounceList: make([][]string, 0),
			CreatedBy:    "github.com/reconbot/swarm-download",
			CreationDate: time.Now().Unix(),
			Comment:      "Lets go fast together!",
		}
		isPrivate := false
		fileInfo := metainfo.Info{
			PieceLength: 100 * 1024 * 1024, // 100mb is a wild guess this needs thought
			Private:     &isPrivate,        // if true we can only use peers from a tracker and not peer exchange
		}
		println("building file from path", pathName)
		err = fileInfo.BuildFromFilePath(pathName)
		if err != nil {
			err := errors.New("Cannot read: " + pathName)
			println(err)
			ch <- CreateTorrentResult{Error: err}
			return
		}

		fileInfo.Name = filepath.Base(pathName)
		metadata.InfoBytes, err = bencode.Marshal(fileInfo)
		if err != nil {
			return
		}
		torrentFile, err := os.Create(torrentPath)
		if err != nil {
			println(err)
			ch <- CreateTorrentResult{Error: err}
			return
		}
		err = metadata.Write(torrentFile)
		if err := torrentFile.Close(); err != nil {
			println(err)
			ch <- CreateTorrentResult{Error: err}
			return
		}
		println("path! ", torrentPath)
		ch <- CreateTorrentResult{path: torrentPath}
		return
	}()

	return ch
}

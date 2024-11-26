package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

type AppConfig struct {
	DataDir     string // where to store internal data
	DownloadDir string // where to save downloaded files
}

func main() {

	var createTorrentCmd = &cobra.Command{
		Use:   "create-torrent [file]",
		Short: "Create a .torrent file out of a file",
		Long:  `Creates a file.torrent from a file with no extra information`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("create-torrent: " + args[0])
			result := <-createTorrent(args[0])
			if result.Error != nil {
				log.Panic("Error creating torrent: ", result)
			}
			fmt.Println("torrent created!: ", result)
		},
	}

	var downloadCmd = &cobra.Command{
		Use:   "download [URI]",
		Short: "Download the uri with local peers",
		Long:  `Download the URI and seed it while you're downloading it`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("download: " + args[0])
			result := <-downloadTorrent(args[0])
			if result.Error != nil {
				log.Panic("Error downloaded torrent: ", result)
			}
			fmt.Println("torrent downloaded!: ", result)
		},
	}

	var infoCmd = &cobra.Command{
		Use:   "info [URI]",
		Short: "Print info about a particular download",
		Long:  `doesn't do anything right now`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("info: " + args[0])
		},
	}

	var app = &cobra.Command{Use: "app"}
	app.AddCommand(createTorrentCmd, downloadCmd, infoCmd)
	app.Execute()
}

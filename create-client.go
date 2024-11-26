package main

import "github.com/anacrolix/torrent"

// Create a torrent client. You will need to close it when you're done
func createClient() (*torrent.Client, error) {
	// full list of options https://github.com/anacrolix/torrent/blob/4d8437a0562147d265591337e04e464dfbf7c348/config.go#L211
	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = "./data"
	clientConfig.DisableWebtorrent = true
	clientConfig.NoDHT = true
	clientConfig.Seed = true

	// disable upnp https://github.com/anacrolix/torrent/blob/4d8437a0562147d265591337e04e464dfbf7c348/portfwd.go#L37-L55
	clientConfig.NoDefaultPortForwarding = true
	clientConfig.DisableIPv6 = true
	clientConfig.ClientTrackerConfig.DisableTrackers = true
	clientConfig.HTTPUserAgent = "reconbot/swarm-download"

	// // DefaultStorage:          storage.NewSqlitePieceCompletion(),
	// // going to choose UTP because it should have a lower priority than tcp connections, need some speed tests
	// // DisableUTP = false,
	// // DisableTCP = true,
	// // HttpRequestDirector func(*http.Request) error, // custom for signed in stuff

	return torrent.NewClient(clientConfig)
}

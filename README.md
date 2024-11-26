# Swarm Downloader

An http downloader that leverages local peers to increase speed. Designed to transfer a changing set large files to a large number of computers quickly making best use a limited upstream network connection. Can be run as a daemon or just a downloader that seeds while downloading.

## Usage 

```bash
swarm create-torrent ./local-file # creates ./local-file.torrent
swarm daemon # an optional long lived seeder that keeps seeding files as long as they exist
swarm download https://foo/file # downloads the file to the storage directory and exits seeding during download
swarm download -b https://foo/file # if the daemon is running tell it to download the file and exit
swarm info https://foo/file # if the daemon is running return swarm and local status of the file
swarm files # returns status of all local files
```

### create-torrent


### download

```bash
swarm download URL
```

This looks for both the URL and `{$URL}.torrent` if both exist a download starts.

### Config file

```json
{
  "storage-directory": "/data/",
  "bind-address": null, // all interfaces
  "preferred-chunk-size": "100mb" // this should be auto somehow
  "bind-address": "0.0.0.0",
}
```

## Decisions

The idea is to leverage local peers if available as they would be high bandwidth but otherwise download via http.

- Bittorrent is great so lets use what we can
- Local network is fast and preferred
- HTTP Source always available and downloads with parallelism.
- No public or extra infrastructure. No Trackers, No DHT
- No UPNP
- Handle web seeds at the application level. In different environments it might be advantageous to use the same data and torrent file in different http locations.

## Features

1. local peer discovery ([BEP-14](https://en.wikipedia.org/wiki/Local_Peer_Discovery) doesn't exist in our library so I made a simple one)
1. HTTP Seeding ([BEP-19](https://www.bittorrent.org/beps/bep_0019.html))
1. Just download or use a local daemon.
1. API is add via cli or delete a file
1. Create a torrent from a file

## Bep 19 Exploration

[Bep 19](https://www.bittorrent.org/beps/bep_0019.html) is design document for "Using HTTP or FTP servers as seeds for BitTorrent downloads.". In it it describes different approaches to using an aternative protocol in addition to bittorrent.

> In GetRight's implementation, I made a couple changes to the usual "rarest first" piece-selection method to better allow "gaps" to develop between pieces. That way there are longer spaces in the file for HTTP and FTP threads to fill. They can start at the beginning of a gap and download until they get to the end.

Pice selection is the primary change.
> If everything else is more-or-less equivalent, it is better to pick a piece to do from the gap of 2 when requesting a piece from a BitTorrent peer.
> In any gap, it is best to fill in from the end (ie, the highest piece number first).

They do this to allow for larger contiguous blocks when downloading from http or ftp

Change the "rarest first" piece selection to a "pretty rare with biggest distance from another completed piece".

> If the client knows the HTTP/FTP download is part of a BitTorrent download, when the very first connection is made it is better to start the HTTP/FTP download somewhere randomly in the file. This way it is more likely the first HTTP pieces it gets will be useful for sharing to the BitTorrent peers.
> If a BitTorrent download is already progressing when starting a HTTP/FTP connection, the HTTP/FTP should start at the beginning of the biggest gap. Given a bitfield "YYnnnnYnnY" it should start at #: "YY#nnnYnnY"

## Development

Large, Fast and Slow are relative to your situation.

- It's helpful to have a web server with a slow connection to your local network.
- It's helpful to have a few computers networked locally with a fast connection.

It's helpful have a large random file to play with. This project hopes to be performant at speeds of 20Gbit/s and files as large as 500Gb, currently untested or tuned.

```bash
# 2^31 == 2GB
openssl rand -out large.file $(( 2**31 ))
```

### Development Notes

- This [discussion suggests](https://github.com/anacrolix/torrent/discussions/953) high throughput is hard to do with this library

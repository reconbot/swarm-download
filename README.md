# Swarm Cache Manager

A network local p2p file cache with remote http (and azure blob store) support. Designed to transfer a changing set large files to a large number of hosts quickly. Can be run as a daemon or just a download client that seeds while downloading.

```bash
swarm-server&
swarm download https://foo/bar?v4 -b # queues the server to fetch the file while backgrounding the request
swarm download https://foo/bar?v4 # waits for the file to become available
swarm info https://foo/bar?v4 # returns swarm and local status of the file
swarm peers # returns peer information 
swarm files # returns status of all available files
```

```json
{
  "storage-directory": "/data/",
  "max-storage-size": "400Gb",
  "speed-limit-up":  "10Gbit",
  "speed-limit-down": null,
  "bind-address": null, // all interfaces
  "preferred-chunk-size": "100mb" // this should be auto somehow
}
```

```bash
swarm-download https://foo/bar?version3 --bind-address 0.0.0.0 --speed-limit-up 10Gbit --speed-limit-down 10Gbit --seed-until-idle=true
```

## Design

- Bittorrent is great so lets steal as much as we can from it and BEP19
- We trust our peers
- We trust the HTTP source to provide an md5 hash of the file in a header
- We check the file in the end with the md5 hash.
- Distribute fetching from the origin as much as possible.

### Advantages

Will we regret the trust later on? Maybe? But without generating a torrent like file this works nicely.

### Pitfalls

Azure blob store, and AWS s3 both provide MD5 hashes of their files, in head and get responses. This will be the swarm key of the file in question. If the file changes, the key changes and blocks from an http source with an incorrect hash will be rejected. A completed file is checked with this checksum at the end. It's recommended you check your file with something more secure later on. Other HTTP servers will need to provide an MD5 header of some sort. Sorry those are the rules.


Goals
- Able to run with minima coordination.

1. local peer discovery
1. remote download start random chunk
1. local chunk discovery
1. prioritize

  - local available chunks we don't have
  - remote chunks that aren't local that we don't have
  - 


## Bep 19
[Bep 19](https://www.bittorrent.org/beps/bep_0019.html) is design document for "Using HTTP or FTP servers as seeds for BitTorrent downloads.". In it it describes different approaches to using an aternative protocol in addition to bittorrent.

> In GetRight's implementation, I made a couple changes to the usual "rarest first" piece-selection method to better allow "gaps" to develop between pieces. That way there are longer spaces in the file for HTTP and FTP threads to fill. They can start at the beginning of a gap and download until they get to the end.

Pice selection is the primary change.
> If everything else is more-or-less equivalent, it is better to pick a piece to do from the gap of 2 when requesting a piece from a BitTorrent peer.
> In any gap, it is best to fill in from the end (ie, the highest piece number first).

They do this to allow for larger contiguous blocks when downloading from http or ftp

Change the "rarest first" piece selection to a "pretty rare with biggest distance from another completed piece".

> If the client knows the HTTP/FTP download is part of a BitTorrent download, when the very first connection is made it is better to start the HTTP/FTP download somewhere randomly in the file. This way it is more likely the first HTTP pieces it gets will be useful for sharing to the BitTorrent peers.
> If a BitTorrent download is already progressing when starting a HTTP/FTP connection, the HTTP/FTP should start at the beginning of the biggest gap. Given a bitfield "YYnnnnYnnY" it should start at #: "YY#nnnYnnY"

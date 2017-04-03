# udp-mirror [![Build Status](https://travis-ci.org/czerwonk/udp-mirror.svg)][travis]
Listens for UDP packets an sends copies to multiple receivers

# Install
```
go get -u github.com/czerwonk/udp-mirror
```
# Application
This tool is helpful if you want to use more than one netflow analasys tool at the same time.

# Use
In this example we want to listen for packets on port 4560. Each packet received should be mirrored and sent to 192.168.1.1:1234 and 192.168.1.2:3456 
```
udp-mirror -listen-address ":4560" -receivers "192.168.1.1:1234,192.168.1.2:3456"
```
[travis]: https://travis-ci.org/czerwonk/udp-mirror

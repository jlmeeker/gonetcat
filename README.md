gonetcat
========

A blend of netcat and dd combined into one executable. Written in Go (golang).

Bulding:
```
go build gonetcat.go
```

Executing:
```Bash
# Server (listen)
gonetcat -l

# Client
gonetcat -client
```

Other options:
```Bash
user@hostname $ gonetcat -h
Usage of gonetcat:
  -B=false: Show results in bytes instead of bits
  -U=false: Use UDP instead of TCP
  -bc="1000": Number of blocks to send (default is 1 thousand) optional suffixes are: k, m, g, t, p, e
  -bs="1000000": Block size (in bytes) for client send (default is 1 megabyte) optional suffixes are: k, m, g, t, p, e
  -c=1: How many consecutive times to run the client transfer test (0 is indefinitely)
  -client=false: Send to remote host
  -d="": Overrides -bc. Total data size to send (in bytes) optional suffixes are: k, m, g, t, p, e
  -l=false: Listen for incoming connections
  -p="2000": Port to listen on
  -repeat=false: Enable echo of received data (reply to sender with received data)
  -s="0.0.0.0": Interface address (or name) to listen on
  -unit="bps": Desired units in which to display results (bps, kbps, mbps, gbps, tbps, pbps, ebps)
```

Sample output:
```Bash
# run as a server and client on the same computer (testing a combination of cpu/memory speeds)
user@hostname $ gonetcat -l -client -bs 100000000 -bc 100
2013/12/31 20:26:20 23057959355.344212 bps (80000000000 bits sent in 3.469518 seconds)

user@hostname $ gonetcat -l -client -bs 100000000 -bc 100 -unit mbps
2013/12/31 20:26:53 23507.409472 Mbps (80000000000 bits sent in 3.403182 seconds)

user@hostname $ gonetcat -client -bs 10000 -bc 100000 -B -unit gbps -c 5
2014/01/01 10:20:08 1.494416 GBps (1000000000 bytes sent in 0.669158 seconds)
2014/01/01 10:20:09 1.537276 GBps (1000000000 bytes sent in 0.650501 seconds)
2014/01/01 10:20:10 1.544670 GBps (1000000000 bytes sent in 0.647388 seconds)
2014/01/01 10:20:10 1.513604 GBps (1000000000 bytes sent in 0.660675 seconds)
2014/01/01 10:20:11 1.537356 GBps (1000000000 bytes sent in 0.650467 seconds)
2014/01/01 10:20:11 Average over all runs:
2014/01/01 10:20:11 1.525233 GBps (1000000000 bytes sent in 0.655638 seconds)

user@hostname $ gonetcat -client -d 200g -B -unit mbps
2014/01/01 20:45:21 3204.516948 MBps (200000000000 bytes sent in 62.411903 seconds)

user@hostname $ gonetcat -client -bs 3m -bc 2k -unit mbps
2014/01/01 22:04:12 25601.544586 Mbps (48000000000 bits sent in 1.874887 seconds)
```

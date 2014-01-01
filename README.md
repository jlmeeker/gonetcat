gonetcat
========

A blend of netcat and dd combined into one executable. Written in Go (golang).

Bulding:
```
go build gonetcat.go
```

Executing:
```Bash
# Server
gonetcat -server

# Client
gonetcat -client
```

Other options:
```Bash
user@hostname $ gonetcat -h
Usage of gonetcat:
  -bcount=1000: Number of blocks to send (default is 1 thousand)
  -bsize=1000000: Block size (in bytes) for client send (default is 1 megabyte)
  -bytes=false: Show results in bytes instead of bits
  -client=false: Send to remote host
  -host="localhost": Interface address (or name) to listen on
  -port="2000": Port to listen on
  -proto="tcp": Protocol to listen on: tcp, udp
  -repeat=false: Enable echo of received data (reply to sender with received data)
  -runs=1: How many consecutive times to run the client transfer test (0 is indefinitely)
  -server=false: Listen for incoming connections
  -unit="bps": Desired units in which to display results (bps, kbps, mbps, gbps, tbps, pbps, ebps, zbps, ybps)
```

Sample output:
```Bash
# run as a server and client on the same computer (testing a combination of cpu/memory speeds)
user@hostname $ gonetcat -server -client -bsize 100000000 -bcount 100
2013/12/31 20:26:20 23057959355.344212 Bps (80000000000 bits sent in 3.469518 seconds)

user@hostname $ gonetcat -server -client -bsize 100000000 -bcount 100 -unit mbps
2013/12/31 20:26:53 23507.409472 Mbps (80000000000 bits sent in 3.403182 seconds)

user@hostname $ gonetcat -client -bsize 10000 -bcount 100000 -bytes -unit gbps -runs 5
2014/01/01 10:20:08 1.494416 GBps (1000000000 bytes sent in 0.669158 seconds)
2014/01/01 10:20:09 1.537276 GBps (1000000000 bytes sent in 0.650501 seconds)
2014/01/01 10:20:10 1.544670 GBps (1000000000 bytes sent in 0.647388 seconds)
2014/01/01 10:20:10 1.513604 GBps (1000000000 bytes sent in 0.660675 seconds)
2014/01/01 10:20:11 1.537356 GBps (1000000000 bytes sent in 0.650467 seconds)
2014/01/01 10:20:11 Average over all runs:
2014/01/01 10:20:11 1.525233 GBps (1000000000 bytes sent in 0.655638 seconds)
```

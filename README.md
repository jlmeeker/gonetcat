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
2013/12/31 21:25:44 1.513957 GBps (1000000000 bytes sent in 0.660521 seconds)
2013/12/31 21:25:44 1.522700 GBps (1000000000 bytes sent in 0.656728 seconds)
2013/12/31 21:25:45 1.562462 GBps (1000000000 bytes sent in 0.640015 seconds)
2013/12/31 21:25:45 1.533771 GBps (1000000000 bytes sent in 0.651988 seconds)
2013/12/31 21:25:46 1.550016 GBps (1000000000 bytes sent in 0.645155 seconds)
```

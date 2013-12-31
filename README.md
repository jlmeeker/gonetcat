gonetcat
========

A blend of netcat and dd combined into one executable.

Bulding:
```
go build gonetcat.go
```

Executing:
```
# Server
gonetcat -server

# Client
gonetcat -client

# Other options
gonetcat -h
Usage of gonetcat:
  -base=1000: Base divisor for doing conversions
  -bcount=1000000: Number of blocks to send
  -bsize=1: Block size for client send (in bytes)
  -client=false: Send to remote host
  -host="localhost": Interface address (or name) to listen on
  -port="2000": Port to listen on
  -proto="tcp": Protocol to listen on: tcp, udp
  -repeat=false: Enable echo of received data (reply to sender with received data)
  -server=false: Listen for incoming connections


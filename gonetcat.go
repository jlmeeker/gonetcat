package main

import (
	//"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

type xferResult struct {
	Bytes int64
	Seconds float64
}

func (r *xferResult) Bits() int64 {
	return r.Bytes*8
}

func (r *xferResult) Bips() float64 {
	return float64(r.Bits()) / r.Seconds
}

func (r *xferResult) KBips() float64 {
	return float64(r.Bips())/base
}

func (r *xferResult) MBips() float64 {
	return float64(r.KBips())/base
}

func (r *xferResult) GBips() float64 {
	return float64(r.MBips())/base
}

func (r *xferResult) Byps() float64 {
	return float64(r.Bytes)/r.Seconds
}

func (r *xferResult) KByps() float64 {
	return float64(r.Byps() / base)
}

func (r *xferResult) MByps() float64 {
	return r.KByps() / base
}

func (r *xferResult) GByps() float64 {
	return r.MByps() / base
}


func serverHandler() {
	// Listen
	l, err := net.Listen(proto, addr+":"+port)

	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			startTime := time.Now()
			var result xferResult			

			if repeat == true {
				// Echo all incoming data.
				result.Bytes,_ = io.Copy(c, c)
			} else {
				// Discard incoming data
				result.Bytes, _ = io.Copy(ioutil.Discard, c)
			}

			// Shut down the connection.
			endTime := time.Now()
			c.Close()
			
			result.Seconds = endTime.Sub(startTime).Seconds()
			log.Printf("%f Mb/s (%d bytes received in %f seconds)", result.MBips(), result.Bytes, result.Seconds)
		}(conn)
	}
}

// Generate and send data
func clientHandler() {
	zero := make([]byte, blocksz, blocksz)
	var result xferResult
	var i int64

	d, err := net.Dial(proto, addr+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()
	
	startTime := time.Now()
	for i = 0; i < blockcount; i++ {
		newBytes,_ := d.Write(zero)
		result.Bytes += int64(newBytes)
	}

	endTime := time.Now()
	result.Seconds = endTime.Sub(startTime).Seconds()
	log.Printf("%f Mb/s (%d bytes sent in %f seconds)", result.MBips(), result.Bytes, result.Seconds)
}

var base float64
var addr string
var port string
var proto string
var repeat bool
var server bool
var client bool
var blocksz int64
var blockcount int64

func init() {
	const (
		defaultAddress = "localhost"
		addrDescr = "Interface address (or name) to listen on"
		defaultPort = "2000"
		portDescr = "Port to listen on"
		defaultProto = "tcp"
		protoDescr = "Protocol to listen on: tcp, udp"
		defaultRepeat = false
		repeatDescr = "Enable echo of received data (reply to sender with received data)"
		defaultServer = false
		serverDescr = "Listen for incoming connections"
		defaultClient = false
		clientDescr = "Send to remote host"
		defaultBlockSize = 1
		blockSizeDescr = "Block size for client send (in bytes)"
		defaultBlockCount = 1000000
		blockCountDescr = "Number of blocks to send"
		defaultBase = 1000
		baseDescr = "Base divisor for doing conversions"
	)
	flag.StringVar(&addr, "host", defaultAddress, addrDescr)
	flag.StringVar(&port, "port", defaultPort, portDescr)
	flag.StringVar(&proto, "proto", defaultProto, protoDescr)
	flag.BoolVar(&repeat, "repeat", defaultRepeat, repeatDescr)
	flag.BoolVar(&server, "server", defaultServer, serverDescr)
	flag.BoolVar(&client, "client", defaultClient, clientDescr)
	flag.Int64Var(&blocksz, "bsize", defaultBlockSize, blockSizeDescr)
	flag.Int64Var(&blockcount, "bcount", defaultBlockCount, blockCountDescr)
	flag.Float64Var(&base, "base", defaultBase, baseDescr)

	flag.Parse()
}

func main() {

	if server {
		go serverHandler()
	} 

	if client {
		time.Sleep(2000)
		clientHandler()
		return
	} 

	if server == false && client == false {
		fmt.Println("You must specify either -server or -client.")
		return
	}

	var input string
    fmt.Scanln(&input)
}
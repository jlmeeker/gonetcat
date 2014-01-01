package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"strings"
	"time"
)

type xferResult struct {
	Bytes float64
	Seconds float64
}

func (r *xferResult) Bits() float64 {
	return r.Bytes * 8
}

func prepUnits() map[string]float64 {
    var units = make(map[string]float64)

    unitNames := [...]string{"bps", "kbps", "mbps", "gbps", "tbps", "pbps", "ebps", "zbps", "ybps"}
    exponent := -3  // set this to -3 so it is zero for the first run (bps)
    for i := range unitNames {
    	exponent += 3 // all powers increase by a factor or 1000 
		units[unitNames[i]] = math.Pow10(exponent)
    }
    return units
}


func formatValue(rawval float64, format string) float64 {
    units := prepUnits()
    return rawval / units[strings.ToLower(format)]
}

func processResult(bytes int64, seconds float64) {
	var result xferResult
	var rate float64
	var bitbyte string
	var totalxferBb int64

	result.Bytes = float64(bytes)
	result.Seconds = seconds

	if usebytes == true {
		rate = formatValue(result.Bytes/result.Seconds, unit)
		bitbyte = "bytes"
		unit = strings.Replace(unit, "b", "B", 1)
		totalxferBb = int64(result.Bytes)
	} else {
		rate = formatValue(result.Bits()/result.Seconds, unit)
		bitbyte = "bits"
		totalxferBb = int64(result.Bits())
	}

	unit = strings.Title(unit)
	log.Printf("%f %s (%d %s sent in %f seconds)", rate, unit, totalxferBb, bitbyte, result.Seconds)
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
			var bytesXferred int64	

			startTime := time.Now()
			if repeat == true {
				// Echo all incoming data.
				bytesXferred,_ = io.Copy(c, c)
			} else {
				// Discard incoming data
				bytesXferred, _ = io.Copy(ioutil.Discard, c)
			}
			endTime := time.Now()

			// Shut down the connection
			c.Close()
			processResult(bytesXferred, endTime.Sub(startTime).Seconds())
		}(conn)
	}
}

// Generate and send data
func clientHandler() {
	zero := make([]byte, blocksz, blocksz)
	var i int64
	var bytesXferred int64

	d, err := net.Dial(proto, addr+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()
	
	startTime := time.Now()
	for i = 0; i < blockcount; i++ {
		newBytes,_ := d.Write(zero)
		bytesXferred += int64(newBytes)
	}
	endTime := time.Now()
	processResult(bytesXferred, endTime.Sub(startTime).Seconds())
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
var unit string
var usebytes bool

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
		defaultBlockSize = 1000000
		blockSizeDescr = "Block size (in bytes) for client send (default is 1 megabyte)"
		defaultBlockCount = 1000
		blockCountDescr = "Number of blocks to send (default is 1 thousand)"
		defaultBase = 1000
		baseDescr = "Base divisor for doing conversions"
		defaultUnit = "bps"
		unitDescr = "Desired units in which to display results (bps, kbps, mbps, gbps, tbps, pbps, ebps, zbps, ybps)"
		defaultBytes = false
		bytesDescr = "Show results in bytes instead of bits"
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
	flag.StringVar(&unit, "unit", defaultUnit, unitDescr)
	flag.BoolVar(&usebytes, "bytes", defaultBytes, bytesDescr)

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
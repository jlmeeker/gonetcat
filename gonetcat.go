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
	"sync"
	"time"
)

type xferResult struct {
	Bytes float64
	Seconds float64
}

// Get our result bytes into bits
func (r *xferResult) Bits() float64 {
	return r.Bytes * 8
}


// Build units array (used in formatValue for division operations)
func prepUnits() map[string]float64 {
    var units = make(map[string]float64)

    unitNames := [...]string{"bps", "kbps", "mbps", "gbps", "tbps", "pbps", "ebps", "zbps", "ybps"}
    exponent := -3  // set this to -3 so it is zero for the first run (bps)
    for i := range unitNames {
    	exponent += 3 // all powers increase by a factor of 1000 
		units[unitNames[i]] = math.Pow10(exponent)
    }
    return units
}


// Get our rate in the specified unit
func formatValue(rawval float64, format string) float64 {
    units := prepUnits()
    return rawval / units[strings.ToLower(format)]
}


// Calculate stats and (optionally) print them to the screen
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


// Run as a server
func serverHandler(showOutput bool) {
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

			seconds := endTime.Sub(startTime).Seconds()
			// Shut down the connection
			c.Close()

			if showOutput {
				processResult(bytesXferred, seconds)
			}
		}(conn)
	}
}


// Run as a client
func clientHandler(showOutput bool) {
	defer wg.Done()

	zero := make([]byte, blocksz, blocksz)
	var i int64
	var bytesXferred int64
	var runsBytesXferred int64
	var runsSeconds float64

	for runNumber := 0; runNumber < runs; runNumber++ {
		bytesXferred = 0
		d, err := net.Dial(proto, addr+":"+port)
		if err != nil {
			log.Fatal(err)
		}
		
		startTime := time.Now()
		for i = 0; i < blockcount; i++ {
			newBytes,_ := d.Write(zero)
			bytesXferred += int64(newBytes)
		}
		endTime := time.Now()

		seconds := endTime.Sub(startTime).Seconds()
		d.Close()
		
		// Crunch numbers and display output (if enabled)
		if showOutput {
			processResult(bytesXferred, seconds)
		}
		
		// Keep track of stats over all runs
		runsBytesXferred += bytesXferred
		runsSeconds += seconds
	}

	if runs > 1 && showOutput {
		log.Println("Average over all runs:")
		avgBytesPerRun := runsBytesXferred/int64(runs)
		avgSeconds := runsSeconds/float64(runs)
		processResult(avgBytesPerRun, avgSeconds)
	}
}


// Global variables
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
var runs int


// Initialize the app
func init() {

	// Define default values and description strings for all flags
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
		defaultUnit = "bps"
		unitDescr = "Desired units in which to display results (bps, kbps, mbps, gbps, tbps, pbps, ebps, zbps, ybps)"
		defaultBytes = false
		bytesDescr = "Show results in bytes instead of bits"
		defaultRuns = 1
		runsDescr = "How many consecutive times to run the client transfer test (0 is indefinitely)"
	)
	flag.StringVar(&addr, "host", defaultAddress, addrDescr)
	flag.StringVar(&port, "port", defaultPort, portDescr)
	flag.StringVar(&proto, "proto", defaultProto, protoDescr)
	flag.BoolVar(&repeat, "repeat", defaultRepeat, repeatDescr)
	flag.BoolVar(&server, "server", defaultServer, serverDescr)
	flag.BoolVar(&client, "client", defaultClient, clientDescr)
	flag.Int64Var(&blocksz, "bsize", defaultBlockSize, blockSizeDescr)
	flag.Int64Var(&blockcount, "bcount", defaultBlockCount, blockCountDescr)
	flag.StringVar(&unit, "unit", defaultUnit, unitDescr)
	flag.BoolVar(&usebytes, "bytes", defaultBytes, bytesDescr)
	flag.IntVar(&runs, "runs", defaultRuns, runsDescr)

	// Validate flags from CLI
	flag.Parse()
}

var wg sync.WaitGroup

func main() {
	
	if server && client {
		go serverHandler(false) // don't log since we're waiting only on the client, so it will log
		time.Sleep(2*time.Second) // Wait for server to listen before we start client tests

		wg.Add(1)  // Only wait for client to complete, server will be killed when client runs are complete
		go clientHandler(true)
	} else if server {
		serverHandler(true) // do not background, run indefinitely
	} else if client {
		wg.Add(1)
		go clientHandler(true)  // It doesn't matter if we background or not here, so why not?
	} 

	if server == false && client == false {
		fmt.Println("You must specify either -server or -client.")
		return
	}

	wg.Wait() // Don't exit until all registered jobs are finished
}
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"strconv"
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

    unitNames := [...]string{"bps", "kbps", "mbps", "gbps", "tbps", "pbps", "ebps"}
    exponent := -3  // set this to -3 so it is zero for the first run (bps)
    for i := range unitNames {
    	exponent += 3 // all powers increase by a factor of 1000 
		units[unitNames[i]] = math.Pow10(exponent)
    }
    return units
}


// Get our rate in the specified unit
func formatValue(method string, rawval float64, format string) (value float64) {
    units := prepUnits()

    if method == "reduce" {
    	value = rawval / units[strings.ToLower(format)]
    } else {
    	value = rawval * units[strings.ToLower(format)]
    }

    return 
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
		rate = formatValue("reduce", result.Bytes/result.Seconds, unit)
		bitbyte = "bytes"
		unit = strings.Replace(unit, "b", "B", 1)
		totalxferBb = int64(result.Bytes)
	} else {
		rate = formatValue("reduce", result.Bits()/result.Seconds, unit)
		bitbyte = "bits"
		totalxferBb = int64(result.Bits())
	}

	unit = strings.Title(unit)
	log.Printf("%f %s (%d %s sent in %f seconds)", rate, unit, totalxferBb, bitbyte, result.Seconds)
}


// Convert string value into usable integer
func parseDataSize(incoming string) (size int64) {
	incomingBytes := []byte(incoming)

	if incoming == "" {
		return
	}

	sizeStr := string(incomingBytes[0:len(incomingBytes)-1])
	suffix := string(incomingBytes[len(incomingBytes)-1])

	if strings.Contains("kmgtpe", suffix) {
		suffix = suffix+"bps"  // lazy re-use of formatValue function
	} else {
		sizeStr = sizeStr + suffix
		suffix = "bps"
	}

	size,_ = strconv.ParseInt(sizeStr,0,0)
	size = int64(formatValue("increase", float64(size), suffix))
	return
}

// Run as a server
func serverHandler(proto string, showOutput bool) {
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
func clientHandler(proto string, showOutput bool) {
	defer wg.Done()

	zero := make([]byte, blockszInt, blockszInt)
	var i int64
	var bytesXferred int64
	var runsBytesXferred int64
	var runsSeconds float64
	var loopForever bool

	if runs == 0 {
		loopForever = true
		runs = 1
	}

	for runNumber := 0; runNumber < runs; runNumber++ {
		if stopExecution {
			break
		}

		bytesXferred = 0
		d, err := net.Dial(proto, addr+":"+port)
		if err != nil {
			log.Fatal(err)
		}
		
		startTime := time.Now()
		for i = 0; i < blockcountInt; i++ {
			if stopExecution {
				break
			}
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

		if loopForever {
			runs += 1
		}
	}

	if runs > 1 && showOutput {
		log.Println("Average over all runs:")
		avgBytesPerRun := runsBytesXferred/int64(runs)
		avgSeconds := runsSeconds/float64(runs)
		processResult(avgBytesPerRun, avgSeconds)
	}
}


// Global variables
var stopExecution bool
var addr string
var port string
var udp bool
var repeat bool
var server bool
var client bool
var blocksz string
var blockcount string
var unit string
var usebytes bool
var runs int
var dataSize string

var blockszInt int64
var blockcountInt int64


// Initialize the app
func init() {

	// Define default values and description strings for all flags
	const (
		defaultAddress = "localhost"
		addrDescr = "Interface address (or name) to listen on"
		defaultPort = "2000"
		portDescr = "Port to listen on"
		defaultUdp = false
		udpDescr = "Use UDP instead of TCP"
		defaultRepeat = false
		repeatDescr = "Enable echo of received data (reply to sender with received data)"
		defaultListen = false
		listenDescr = "Listen for incoming connections"
		defaultClient = false
		clientDescr = "Send to remote host"
		defaultBlockSize = "1000000"
		blockSizeDescr = "Block size (in bytes) for client send (default is 1 megabyte) optional suffixes are: k, m, g, t, p, e"
		defaultBlockCount = "1000"
		blockCountDescr = "Number of blocks to send (default is 1 thousand) optional suffixes are: k, m, g, t, p, e"
		defaultUnit = "bps"
		unitDescr = "Desired units in which to display results (bps, kbps, mbps, gbps, tbps, pbps, ebps)"
		defaultBytes = false
		bytesDescr = "Show results in bytes instead of bits"
		defaultRuns = 1
		runsDescr = "How many consecutive times to run the client transfer test (0 is indefinitely)"
		defaultDataSize = ""
		dataSizeDescr = "Overrides -bc. Total data size to send (in bytes) optional suffixes are: k, m, g, t, p, e"
	)
	flag.StringVar(&addr, "s", defaultAddress, addrDescr)
	flag.StringVar(&port, "p", defaultPort, portDescr)
	flag.BoolVar(&udp, "U", defaultUdp, udpDescr)
	flag.BoolVar(&repeat, "repeat", defaultRepeat, repeatDescr)
	flag.BoolVar(&server, "l", defaultListen, listenDescr)
	flag.BoolVar(&client, "client", defaultClient, clientDescr)
	flag.StringVar(&blocksz, "bs", defaultBlockSize, blockSizeDescr)
	flag.StringVar(&blockcount, "bc", defaultBlockCount, blockCountDescr)
	flag.StringVar(&unit, "unit", defaultUnit, unitDescr)
	flag.BoolVar(&usebytes, "B", defaultBytes, bytesDescr)
	flag.IntVar(&runs, "c", defaultRuns, runsDescr)
	flag.StringVar(&dataSize, "d", defaultDataSize, dataSizeDescr)

	// Validate flags from CLI
	flag.Parse()

	// Let things run
	stopExecution = false
}

var wg sync.WaitGroup

func main() {
	proto := "tcp"
	if udp {
		proto = "udp"
	}

	blockszInt = int64(parseDataSize(blocksz))
	blockcountInt = int64(parseDataSize(blockcount))

	// Adjust block count if dataSize was set
	parsedDataSize := parseDataSize(dataSize)
	if dataSize != "" && parsedDataSize != int64(0) {
		if (parsedDataSize > blockszInt) {
			blockcountInt = parsedDataSize/blockszInt
		} else {
			log.Println("Specified data size is smaller than block size, reducing block size.")
			blockszInt = parsedDataSize
			blockcountInt = parsedDataSize/blockszInt
		}
		
	} else if dataSize != "" && parsedDataSize == int64(0) {
		log.Println("Error parsing -d value, using -bc value instead.")
	}

	// Catch Ctrl-C so we can exit cleanly
	c := make(chan os.Signal, 1)                                       
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)                                     
	go func() {                                                        
		for _ = range c {                                             
			stopExecution = true
			if client == false {
				os.Exit(1)
			}
		}                                                            
	}()
	
	// Run server and/or client
	if server && client {
		go serverHandler(proto, false) // don't log since we're waiting only on the client, so it will log
		time.Sleep(2*time.Second) // Wait for server to listen before we start client tests

		wg.Add(1)  // Only wait for client to complete, server will be killed when client runs are complete
		go clientHandler(proto, true)
	} else if server {
		serverHandler(proto, true) // do not background, run indefinitely
	} else if client {
		wg.Add(1)
		go clientHandler(proto, true)  // It doesn't matter if we background or not here, so why not?
	} 

	if server == false && client == false {
		fmt.Println("You must specify at least one of -l, -client.")
		return
	}

	wg.Wait() // Don't exit until all registered jobs are finished
}
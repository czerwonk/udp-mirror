package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const version string = "0.1"

var (
	listenAddress     = flag.String("listen-address", ":9999", "UDP port to listen for incoming packets")
	receiverAddresses = flag.String("receivers", "", "comma seperated list of copy receivers")
	bufferSize        = flag.Int("buffer-size", 1024, "size of read buffer")
	debug             = flag.Bool("debug", false, "debug mode")
	showVersion       = flag.Bool("version", false, "show version info")
)

type receiver struct {
	address string
	channel chan []byte
}

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: udp-mirror [ ... ]\n\nParameters:\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	if len(*receiverAddresses) == 0 {
		fmt.Println("No receivers defined!")
		os.Exit(1)
	}

	receivers := getReceivers()

	for _, r := range receivers {
		go tee(r)
	}

	startServer(receivers)
}

func printVersion() {
	fmt.Println("udp-mirror")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author: Daniel Czerwonk")
	fmt.Println("Source code: https://github.com/czerwonk/udp-mirror")
}

func startServer(receivers []*receiver) {
	log.Println("Starting listening")
	conn, err := net.ListenPacket("udp", *listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("Listening on %s. Waiting for packets", *listenAddress)

	for {
		buf := make([]byte, *bufferSize)
		len, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println(err)
		}

		if *debug {
			log.Printf("Received (%d)", len)
		}

		go func() {
			for _, r := range receivers {
				r.channel <- buf[:len]
			}
		}()
	}
}

func getReceivers() []*receiver {
	receivers := make([]*receiver, 0)

	for _, x := range strings.Split(*receiverAddresses, ",") {
		r := &receiver{address: strings.TrimSpace(x), channel: make(chan []byte)}
		receivers = append(receivers, r)
	}

	return receivers
}

func tee(r *receiver) {
	conn, err := net.Dial("udp", r.address)
	if err != nil {
		log.Println("Could not add receiver %s: %s", r.address, err)
		return
	}

	log.Printf("Adding receiver: %s\n", r.address)

	for {
		d := <-r.channel
		_, err := conn.Write(d)
		if err != nil {
			log.Printf("%s: %s", r.address, err)
		}

		if *debug {
			log.Printf("Packet sent to %s (%d)", r.address, len(d))
		}
	}
}

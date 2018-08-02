package main

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/alandtsang/zmux"
)

var buffer []byte
var pkts int32
var recvBytes int32

func main() {
	go displayTimerProc()

	listener, err := net.Listen("tcp", ":9999")
	//listener, err := net.ListenPacket("udp", ":9999")
	if err != nil {
		log.Fatal(err)
	}
	// Accept a TCP connection
	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", conn)

	/*conf := &zmux.Config{
		AcceptBacklog:          256,
		EnableKeepAlive:        true,
		KeepAliveInterval:      30 * time.Second,
		ConnectionWriteTimeout: 10 * time.Second,
		MaxStreamWindowSize:    1024 * 1024,
		LogOutput:              os.Stderr,
	}*/

	// Setup server side of zmux
	session, err := zmux.Server(conn, nil)
	//session, err := zmux.Server(conn, conf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", session)

	for {
		// Accept a stream
		stream, err := session.Accept()
		if err != nil {
			panic(err)
		}

		go handleStream(&stream)

		/*buffer = make([]byte, 1600)

		for {
			// Listen for a message
			n, err := stream.Read(buffer)
			if err != nil {

			}
			pkts++
			recvBytes += n
		}*/
	}
}

func handleStream(stream *net.Conn) {
	//buffer = make([]byte, 1600)
	buffer = make([]byte, 1024*64)

	for {
		// Listen for a message
		n, err := (*stream).Read(buffer)
		if err != nil {

		}
		atomic.AddInt32(&pkts, 1)
		atomic.AddInt32(&recvBytes, int32(n))
		//pkts++
		//recvBytes += n
	}
}

func displayTimerProc() {
	for {
		time.Sleep(time.Second)
		//fmt.Printf("Pkts %d, Bytes %d, Rate %d Mbps\n",
		//  pkts, bytes, bytes*8/(1000*1000))
		fmt.Printf("Pkts %d, Bytes %d, Rate %d M/s\n",
			pkts, recvBytes, recvBytes/(1000*1000))
		pkts = 0
		recvBytes = 0
	}
}

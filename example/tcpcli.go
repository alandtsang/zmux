package main

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/alandtsang/zmux"
)

var buffer []byte
var pkts int32
var sendBytes int32

func main() {
	go displayTimerProc()

	// Get a TCP connection
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	//conn, err := net.Dial("tcp", "192.168.110.8:9999")
	if err != nil {
		panic(err)
	}

	/*conf := &zmux.Config{
		AcceptBacklog:          256,
		EnableKeepAlive:        true,
		KeepAliveInterval:      30 * time.Second,
		ConnectionWriteTimeout: 10 * time.Second,
		MaxStreamWindowSize:    1024 * 1024,
		LogOutput:              os.Stderr,
	}*/
	// Setup client side of zmux
	session, err := zmux.Client(conn, nil)
	//session, err := zmux.Client(conn, conf)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		//for i := 0; i < 1; i++ {
		// Open a new stream
		stream, err := session.Open()
		if err != nil {
			panic(err)
		}

		go Write(&stream)
	}

	time.Sleep(60 * time.Second)
}

func Write(stream *net.Conn) {
	buffer = make([]byte, 1024*64)

	for {
		//for i := 0; i < 10; i++ {
		//n, err := (*stream).Write(buffer[:1472])
		n, err := (*stream).Write(buffer)
		if err != nil {
			fmt.Println(err)
		}

		atomic.AddInt32(&pkts, 1)
		atomic.AddInt32(&sendBytes, int32(n))
	}
}

func displayTimerProc() {
	for {
		time.Sleep(time.Second)
		//fmt.Printf("Pkts %d, Bytes %d, Rate %d Mbps\n",
		//  pkts, bytes, bytes*8/(1000*1000))
		fmt.Printf("Pkts %d, Bytes %d, Rate %d M/s\n",
			pkts, sendBytes, sendBytes/(1000*1000))
		pkts = 0
		sendBytes = 0
	}
}

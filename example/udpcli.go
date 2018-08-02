package main

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	"./util"
	"github.com/alandtsang/zmux"
)

var buffer []byte
var pkts int32
var sendBytes int32

type Conn = util.Conn

type client struct {
	conn       net.UDPConn
	remoteAddr net.Addr
	hostname   string
}

func DialAddr(addr string) (*client, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}

	hostname, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	return &client{
		conn:       *udpConn,
		hostname:   hostname,
		remoteAddr: udpAddr,
	}, nil
}

func (c *client) NewConn() *Conn {
	return util.NewConn(c.conn, c.remoteAddr)
}

func main() {
	go displayTimerProc()

	client, err := DialAddr("127.0.0.1:9999")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("client=%+v\n", client)

	// Get a TCP connection
	conn := client.NewConn()
	fmt.Printf("conn=%+v\n", conn)

	// Setup client side of zmux
	session, err := zmux.Client(conn, nil)
	//session, err := zmux.Client(conn, conf)
	if err != nil {
		panic(err)
	}

	//for i := 0; i < 10; i++ {
	for i := 0; i < 1; i++ {
		// Open a new stream
		stream, err := session.Open()
		if err != nil {
			panic(err)
		}
		//fmt.Printf("stream=%+v\n", stream)
		go Write(&stream)
	}

	time.Sleep(60 * time.Second)
}

func Write(stream *net.Conn) {
	buffer = make([]byte, 1024*2)

	for {
		//for i := 0; i < 10; i++ {
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

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

type Conn = util.Conn

var buffer []byte
var pkts int32
var recvBytes int32

type Listener interface {
	Accept() (*Conn, error)
	Addr() net.Addr
	Close() error
}

type server struct {
	conn       net.UDPConn
	remoteAddr net.Addr
	connQueue  chan Conn
}

var _ Listener = &server{}

func (s *server) Accept() (*Conn, error) {
	var conn Conn
	select {
	case conn = <-s.connQueue:
		fmt.Printf("========= conn=%+v\n", conn)
		return &conn, nil
	}
}

func (s *server) Close() error {
	err := s.conn.Close()
	return err
}

func (s *server) Addr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *server) serve() {
	for {
		/*var buffer [1024]byte
		n, remoteAddr, err := s.conn.ReadFrom(buffer[:])
		if err != nil {
			fmt.Println(err)
			s.Close()
			return
		}
		s.remoteAddr = remoteAddr
		fmt.Println("+++++ remoteAddr=", remoteAddr)
		fmt.Println("+++++ ", n, string(buffer[:n]))*/

		//atomic.AddInt32(&pkts, 1)
		//atomic.AddInt32(&recvBytes, int32(n))

		/*data := buffer[:n]
		if err := s.handlePacket(remoteAddr, data); err != nil {
			fmt.Println("error handling packet: ", err)
		}*/

		//conn := util.NewConn(s.conn, remoteAddr)
		conn := util.NewConn(s.conn, nil)
		s.connQueue <- *conn
	}
}

//func (s *UDPServer) handlePacket(remoteAddr net.Addr, packet []byte) error {
//}

func ListenUDP(network, address string) (Listener, error) {
	la, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP(network, la)
	if err != nil {
		log.Fatal(err)
	}

	s := &server{
		conn:      *conn,
		connQueue: make(chan Conn, 5),
	}

	go s.serve()
	return s, nil
}

func main() {
	go displayTimerProc()

	//listener, err := net.Listen("tcp", ":9999")
	//conn, err := net.ListenPacket("udp", ":9999")
	listener, err := ListenUDP("udp", ":9999")
	if err != nil {
		log.Fatal(err)
	}

	// Accept a TCP connection
	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	// Setup server side of zmux
	session, err := zmux.Server(conn, nil)
	//session, err := zmux.Server(conn, conf)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Accept a stream
		stream, err := session.Accept()
		if err != nil {
			panic(err)
		}

		go handleStream(&stream)
	}
}

func handleStream(stream *net.Conn) {
	buffer = make([]byte, 1600)

	fmt.Println("-> handleStream")
	for {
		// Listen for a message
		n, err := (*stream).Read(buffer)
		//fmt.Println("===", n)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(string(buffer[:n]))
		atomic.AddInt32(&pkts, 1)
		atomic.AddInt32(&recvBytes, int32(n))
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

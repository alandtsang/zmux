package util

import (
	"net"
)

type Conn struct {
	conn       net.UDPConn
	remoteAddr net.Addr
}

func NewConn(conn net.UDPConn, remoteAddr net.Addr) *Conn {
	return &Conn{
		conn:       conn,
		remoteAddr: remoteAddr,
	}
}

func (c *Conn) Read(b []byte) (int, error) {
	//fmt.Println("Read()")
	n, remoteAddr, err := c.conn.ReadFrom(b)
	c.remoteAddr = remoteAddr
	//fmt.Println("n=", n)
	return n, err
}

func (c *Conn) Write(b []byte) (int, error) {
	n, err := c.conn.WriteTo(b, c.remoteAddr)
	return n, err
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

/*func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}*/

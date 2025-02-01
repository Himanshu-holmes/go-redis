package main

import (
	"fmt"
	"net"
)

type Peer struct {
	conn  net.Conn
	msgCh chan []byte
}

func NewPeer(conn net.Conn, msgCh chan []byte) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
	}
}

func (p *Peer) readLoop() error {

	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		fmt.Println("n", n)
		if err != nil {
			fmt.Println("error reading from conn", err)
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
		p.msgCh <- msgBuf
	}
}

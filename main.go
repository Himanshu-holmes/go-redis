package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
)

const defaultListenAddress = ":3000"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan []byte
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddress
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan []byte),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	fmt.Println("Listening", s.ListenAddr)
	go s.loop()
	slog.Info("server started", "addr", s.ListenAddr)

	return s.acceptLoop()
}
func (s *Server) handleRawMessage(rawMsg []byte) error {

	return nil
}
func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgCh:
			if err := s.handleRawMessage(rawMsg); err != nil {
				slog.Error("raw message error", "err", err)
			}
			
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true

		}

	}
}

func (s *Server) acceptLoop() error {

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}
		go s.handleConn(conn)

	}
}

func (s *Server) handleConn(conn net.Conn) {

	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	slog.Info("new peer connected", "RemoteAddr", conn.RemoteAddr())
	slog.Info("server peers", "peer", s.peers)
	err := peer.readLoop()
	fmt.Println("error reading from conn", err)
	if err != nil {
		slog.Info("peer read error", "err", err, "RemoteAddr", conn.RemoteAddr())
	}

}

func main() {
	server := NewServer(Config{})

	log.Fatal(server.Start())
}

package gsservice

import (
	"fmt"
	"github.com/jfy0o0/goStealer/net/gsservice/iface"
	"github.com/jfy0o0/goStealer/net/gstcp"
	"net"
	"time"
)

type Server struct {
	port                int
	MsgHandler          iface.IMsgHandler
	ConnMgr             iface.IConnectionManager
	idProducer          iface.IConnectionIDProducer
	listener            *net.TCPListener
	FreshExpireInterval int64
}

func NewServer(port int, ConnMgr iface.IConnectionManager, MsgHandler iface.IMsgHandler) *Server {
	return &Server{
		port:                port,
		MsgHandler:          MsgHandler,
		ConnMgr:             ConnMgr,
		idProducer:          nil,
		FreshExpireInterval: 120,
	}
}

func (s *Server) SetFreshExpireInterval(expire int64) {
	s.FreshExpireInterval = expire
}

func (s *Server) SetIDProducer(idProducer iface.IConnectionIDProducer) {
	s.idProducer = idProducer
}

func (s *Server) Start() (err error) {
	s.listener, err = net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: s.port,
	})
	if err != nil {
		return err
	}
	defer s.listener.Close()

	fmt.Println("[server TCP listener start SUCCESS]:", s.port)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}
		c := gstcp.NewConnByNetConn(conn)

		connID := conn.RemoteAddr().String()
		if s.idProducer != nil {
			connID = s.idProducer.ProduceID()
		}
		dealConn := newConnectionAsServer(c, connID, s.MsgHandler)
		go func(iConn iface.IConnection) {
			s.ConnMgr.Add(iConn)
			defer s.ConnMgr.Del(iConn)
			iConn.Start()
		}(dealConn)
	}

	s.ConnMgr.Clear()

	return nil
}

func (s *Server) Stop() {
	s.listener.SetDeadline(time.Now())
}

func (s *Server) timer() {
	s.ConnMgr.Walk(func(m map[string]iface.IConnection) {
		now := time.Now().Unix()
		for k, v := range m {
			if !v.IsCmdChan() {
				continue
			}
			if v.GetFresh()+s.FreshExpireInterval < now {
				delete(m, k)
				v.Stop()
			}
		}
	})
}

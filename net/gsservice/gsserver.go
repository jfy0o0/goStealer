package gsservice

import (
	"context"
	"fmt"
	"github.com/jfy0o0/goStealer/net/gsservice/iface"
	"github.com/jfy0o0/goStealer/net/gstcp"
	"net"
	"time"
)

type Server struct {
	ip         string
	port       int
	MsgHandler iface.IMsgHandler
	ConnMgr    iface.IConnectionManager
	//DoExitChan chan os.Signal
	idProducer iface.IConnectionIDProducer
	listener   net.Listener
}

func NewServer(ip string, port int, ConnMgr iface.IConnectionManager, MsgHandler iface.IMsgHandler) *Server {
	return &Server{
		ip:         ip,
		port:       port,
		MsgHandler: MsgHandler,
		ConnMgr:    ConnMgr,
		//DoExitChan: make(chan os.Signal, 1),
		idProducer: nil,
	}
}

func (s *Server) SetIDProducer(idProducer iface.IConnectionIDProducer) {
	s.idProducer = idProducer
}

func (s *Server) Start() (err error) {
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
	if err != nil {
		return err
	}
	defer s.listener.Close()

	fmt.Println("[server TCP listener start SUCCESS]:", s.ip, s.port)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}
		c := gstcp.NewConnByNetConn(conn)

		ctx, cancel := context.WithCancel(context.Background())
		connID := conn.RemoteAddr().String()
		if s.idProducer != nil {
			connID = s.idProducer.ProduceID()
		}
		dealConn := newConnectionAsServer(ctx, cancel, c, connID, s.MsgHandler)
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
	s.listener.Close()
}

func (s *Server) timer() {
	s.ConnMgr.Walk(func(m map[string]iface.IConnection) {
		now := time.Now().Unix()
		for k, v := range m {
			if !v.IsCmdChan() {
				continue
			}
			if v.GetFresh()+180 < now {
				delete(m, k)
			}
		}
	})
}

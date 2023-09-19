package gsnet

import (
	"github.com/jfy0o0/goStealer/container/gsmap"
	"github.com/jfy0o0/goStealer/net/gstcp"
	"github.com/jfy0o0/goStealer/os/gstimer"
	"net"
)

type Server[T any] struct {
	Connections *gsmap.AnyAnyMap[string, *Session[T]]
	Config      *ServerConfig[T]
	timer       *gstimer.Timer
	listener    net.Listener
}

func NewServer[T any](configs ...*ServerConfig[T]) *Server[T] {
	config := GetDefaultServerConfig[T]()
	if len(configs) > 0 {
		config = configs[0]
	}
	server := &Server[T]{
		Connections: gsmap.NewAnyAnyMap[string, *Session[T]](true),
		Config:      config,
		timer:       gstimer.New(),
	}
	server.timer.Stop()
	if config.CheckBeatHeartInterval != 0 {
		server.timer.AddSingleton(config.CheckBeatHeartInterval, server.timerCheckHeartBeat)
	}
	return server
}
func (s *Server[T]) timerCheckHeartBeat() {
	s.Connections.LockFunc(func(m map[string]*Session[T]) {
		for k, session := range m {
			if err := session.CheckHeartBeat(); err != nil {
				s.Config.OnHeartBeatFailed(k)
				continue
			}
			s.Config.OnHeartBeatSuccessful(k)
		}
	})
}
func (s *Server[T]) Run() (err error) {
	s.timer.Start()
	s.listener, err = net.Listen("tcp", s.Config.ListenAddr)
	if err != nil {
		return err
	}
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			break
		}
		go s.process(gstcp.NewConnByNetConn(conn))
	}
	return err
}
func (s *Server[T]) process(conn *gstcp.Conn) {
	defer conn.Close()
	if err := s.Config.OnConnectedStart(conn); err != nil {
		return
	}
	clientHello, err := s.Config.Hello.HandShakeAsServer(conn)
	if err != nil {
		return
	}
	if clientHello.Data.Key == "" {
		clientHello.Data.Key = conn.RemoteAddr().String()
	}
	newConn := gstcp.UpgradeConnAsServer(conn)
	if err = s.Config.OnConnectedHandClientHello(clientHello); err != nil {
		return
	}
	session := newServerSession[T](s)
	session.Hello = clientHello

	if err = session.InitSelf(true, newConn); err != nil {
		return
	}
	//session, err := newServerSession[T](s, newConn, clientHello)
	//if err != nil {
	//	return
	//}

	s.Connections.Set(clientHello.Data.Key, session)
	session.Run()
	s.Connections.Remove(clientHello.Data.Key)
}
func (s *Server[T]) Close() {
	s.Connections.LockFunc(func(m map[string]*Session[T]) {
		for k, v := range m {
			v.Stop()
			delete(m, k)
		}
	})

	s.listener.Close()
}

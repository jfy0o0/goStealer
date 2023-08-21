package gsnet

import (
	"github.com/jfy0o0/goStealer/container/gsmap"
	"github.com/jfy0o0/goStealer/net/gstcp"
	"github.com/jfy0o0/goStealer/os/gstimer"
	"net"
)

type Server[T any] struct {
	Connections *gsmap.AnyAnyMap[string, *WorkerSession[T]]
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
		Connections: gsmap.NewAnyAnyMap[string, *WorkerSession[T]](true),
		Config:      config,
		timer:       gstimer.New(),
	}
	server.timer.Stop()
	server.timer.AddSingleton(config.CheckBeatHeartInterval, server.timerCheckHeartBeat)
	return server
}
func (s *Server[T]) timerCheckHeartBeat() {
	s.Connections.LockFunc(func(m map[string]*WorkerSession[T]) {

		for k, workSession := range m {
			if workSession.YamuxSession == nil {
				continue
			}
			if _, err := workSession.YamuxSession.Ping(); err != nil {
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
	newConn := gstcp.UpgradeConnAsServer(conn)
	if err = s.Config.OnConnectedHandClientHello(clientHello); err != nil {
		return
	}

	session, err := newServerSession[T](s, newConn, clientHello)
	if err != nil {
		return
	}
	s.Connections.Set(clientHello.Data.Key, session)
	session.Run()
	s.Connections.Remove(clientHello.Data.Key)
}

package gsnet

import (
	"context"
	"fmt"
	"github.com/hashicorp/yamux"
	"github.com/jfy0o0/goStealer/errors/gserror"
	"github.com/jfy0o0/goStealer/net/gstcp"
	"net"
)

type CommunicationYamuxMuti[T any] struct {
	// Channel to notify that the connection has exited/stopped
	// (告知该链接已经退出/停止的channel)
	ctx           context.Context
	cancel        context.CancelFunc
	Tx            chan interface{}
	Conn          *gstcp.Conn
	ParentSession *Session[T]
	YamuxSession  *yamux.Session
	isServer      bool
}

func (c *CommunicationYamuxMuti[T]) InitSelf(isServer bool, conn *gstcp.Conn) (err error) {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//c.Tx = make(chan interface{}, 1024)
	c.isServer = isServer
	c.Conn = conn
	if c.isServer {
		c.YamuxSession, err = yamux.Server(c.Conn, nil)
	} else {
		c.YamuxSession, err = yamux.Server(c.Conn, nil)
	}
	return err
}

func (c *CommunicationYamuxMuti[T]) Run() {
	go c.runTx()
	for {
		stream, err := c.YamuxSession.Accept()
		if err != nil {
			break
		}
		go c.handleStream(gstcp.NewConnByNetConn(stream))
	}
	//<-c.ctx.Done()
	fmt.Println("exit run")
}

func (c *CommunicationYamuxMuti[T]) runTx() {
	for {
		select {
		case <-c.ctx.Done():
			c.YamuxSession.Close()
			return
		case msg, ok := <-c.Tx:
			if !ok {
				continue
			}
			var err error
			var conn net.Conn
			for {
				conn, err = c.YamuxSession.Open()
				if err == nil {
					break
				}
				if c.ctx.Err() != nil {
					return
				}
			}
			c.ParentSession.Adapter.OnSendMsg(gstcp.NewConnByNetConn(conn), msg)
			conn.Close()
		}
	}
}

func (s *CommunicationYamuxMuti[T]) handleStream(conn *gstcp.Conn) {
	defer conn.Close()
	s.ParentSession.Adapter.OnMsg(conn)
}

func (c *CommunicationYamuxMuti[T]) Stop() {
	c.YamuxSession.Close()
	c.cancel()
}

func (c *CommunicationYamuxMuti[T]) Push(msg interface{}) {
	c.Tx <- msg
}
func (c *CommunicationYamuxMuti[T]) CheckHeartBeat() error {
	if c.YamuxSession == nil {
		return gserror.Newf("yamux session is nil")
	}
	_, err := c.YamuxSession.Ping()

	return err
}

package gsnet

import (
	"context"
	"fmt"
	"github.com/jfy0o0/goStealer/net/gstcp"
)

type CommunicationUserDefined[T any] struct {
	// Channel to notify that the connection has exited/stopped
	// (告知该链接已经退出/停止的channel)
	ctx           context.Context
	cancel        context.CancelFunc
	Tx            chan interface{}
	Conn          *gstcp.Conn
	ParentSession *Session[T]
	isServer      bool
}

func (c *CommunicationUserDefined[T]) InitSelf(isServer bool, conn *gstcp.Conn) error {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//c.Tx = make(chan interface{}, 1024)
	c.isServer = isServer
	c.Conn = conn
	return nil
}
func (c *CommunicationUserDefined[T]) Run() {
	go c.runTx()
	go c.runRx()
	<-c.ctx.Done()
	fmt.Println("CommunicationUserDefined run exit")
}

func (c *CommunicationUserDefined[T]) runRx() {
	c.ParentSession.Adapter.OnMsg(c.Conn)
	c.Stop()
}

func (c *CommunicationUserDefined[T]) runTx() {
	for {
		select {
		case <-c.ctx.Done():
			c.Conn.Close()
			fmt.Println("close conn ")
			return
		case msg, ok := <-c.Tx:
			if !ok {
				continue
			}
			c.ParentSession.Adapter.OnSendMsg(c.Conn, msg)
		}
	}
}

func (c *CommunicationUserDefined[T]) Stop() {
	c.cancel()
}

func (c *CommunicationUserDefined[T]) Push(msg interface{}) {
	c.Tx <- msg
}
func (c *CommunicationUserDefined[T]) CheckHeartBeat() error {
	return nil
}

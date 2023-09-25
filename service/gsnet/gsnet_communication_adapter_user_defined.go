package gsnet

import (
	"context"
	"github.com/jfy0o0/goStealer/container/gstype"
	"github.com/jfy0o0/goStealer/net/gstcp"
	"log"
	"time"
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
	IsRun         *gstype.Bool
}

func NewCommunicationUserDefinedFromConfig[T any](config SessionConfig, session *Session[T]) *CommunicationUserDefined[T] {
	c := &CommunicationUserDefined[T]{
		ParentSession: session,
		Tx:            make(chan interface{}, config.TxCap),
		IsRun:         gstype.NewBool(false),
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go c.runTx()
	return c
}

func (c *CommunicationUserDefined[T]) InitSelf(isServer bool, conn *gstcp.Conn) error {
	c.isServer = isServer
	c.Conn = conn
	c.IsRun.Set(true)
	log.Println("re init self=  ", conn.LocalAddr().String())
	return nil
}
func (c *CommunicationUserDefined[T]) Run() {
	c.runRx()
	c.IsRun.Set(false)
	log.Println("CommunicationUserDefined run exit")
}

func (c *CommunicationUserDefined[T]) runRx() {
	c.ParentSession.Adapter.OnMsg(c.Conn)
}

func (c *CommunicationUserDefined[T]) runTx() {
	//defer log.Println("exit tx")
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg, ok := <-c.Tx:
			if !ok {
				return
			}
			for {
				if !c.IsRun.Val() {
					if c.ctx.Err() != nil {
						return
					}
					time.Sleep(time.Second)
					continue
				}
				if err := c.ParentSession.Adapter.OnSendMsg(c.Conn, msg); err == nil {
					break
				}
				time.Sleep(time.Second)
			}
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

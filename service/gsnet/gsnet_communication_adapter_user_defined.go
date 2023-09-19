package gsnet

import (
	"context"
	"fmt"
	"github.com/jfy0o0/goStealer/container/gstype"
	"github.com/jfy0o0/goStealer/net/gstcp"
	"log"
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

	return c
}

func (c *CommunicationUserDefined[T]) InitSelf(isServer bool, conn *gstcp.Conn) error {
	//c.Tx = make(chan interface{}, 1024)
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.isServer = isServer
	c.Conn = conn
	log.Println("re init self=  ", conn.LocalAddr().String())
	return nil
}
func (c *CommunicationUserDefined[T]) Run() {
	c.IsRun.Set(true)
	defer c.IsRun.Set(false)
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
	//defer log.Println("tx chan exit ...")
	for {
		select {
		case <-c.ctx.Done():
			c.Conn.Close()
			//fmt.Println("close conn ")
			return
		case msg, ok := <-c.Tx:
			if !ok {
				continue
			}
			for {
				if !c.IsRun.Val() {
					if c.ctx.Err() != nil {
						return
					}
					continue
				}
				c.ParentSession.Adapter.OnSendMsg(c.Conn, msg)
				break
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

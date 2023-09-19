package gsnet

import (
	"github.com/jfy0o0/goStealer/container/gstype"
	"github.com/jfy0o0/goStealer/net/gstcp"
	"log"
	"net"
	"time"
)

type Client[T any] struct {
	Config  *ClientConfig[T]
	isRun   *gstype.Bool
	Session *Session[T]
}

func NewClient[T any](configs ...*ClientConfig[T]) *Client[T] {
	config := GetDefaultClientConfig[T]()
	if len(configs) > 0 {
		config = configs[0]
	}
	client := &Client[T]{
		Config: config,
		isRun:  gstype.NewBool(false),
	}
	client.Session = newClientSession[T](client)
	return client
}

func (c *Client[T]) Run() {
	//go c.runTx()
	var sleepTime = 5
	for {
		conn, err := net.Dial("tcp", c.Config.ConnAddr)
		if err != nil {
			log.Println(err)
			sleepTime *= 2
			if sleepTime > 180 {
				sleepTime = 180
			}
		} else {
			log.Printf("conn to %v success ", conn.RemoteAddr().String())
			c.do(gstcp.NewConnByNetConn(conn))
			sleepTime = 5
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}
func (c *Client[T]) Stop() {
	c.isRun.Set(false)
	c.Session.Stop()
}

//func (c *Client[T]) runTx() {
//	for v := range c.Tx {
//		c.Session.Push(msg)
//
//	}
//}

func (c *Client[T]) do(conn *gstcp.Conn) {
	defer conn.Close()
	defer c.isRun.Set(false)
	var err error
	if err = c.Config.OnConnectedStart(conn); err != nil {
		return
	}

	serverHello, err := c.Config.Hello.HandShakeAsClient(conn)
	if err != nil {
		return
	}
	newConn := gstcp.UpgradeConnAsClient(conn)
	if err = c.Config.OnConnectedHandServerHello(serverHello); err != nil {
		return
	}

	//if c.Session != nil {
	//	c.Session.Stop()
	//	c.Session = nil
	//}
	//c.Session, err = newClientSession[T](c, newConn, serverHello)
	//if err != nil {
	//	return
	//}
	c.Session.Hello = serverHello
	if err = c.Session.CommunicationAdapter.InitSelf(false, newConn); err != nil {
		return
	}

	c.isRun.Set(true)

	c.Session.Run()
}

func (c *Client[T]) IsConnected() bool {
	return c.isRun.Val()
}

package gsservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jfy0o0/goStealer/errors/gserror"
	"github.com/jfy0o0/goStealer/net/gsservice/iface"
	"github.com/jfy0o0/goStealer/net/gsservice/internal"
	"github.com/jfy0o0/goStealer/net/gstcp"
	"sync"
	"sync/atomic"
	"time"
)

type Connection struct {
	conn         *gstcp.Conn
	connID       string
	property     map[string]interface{}
	propertyLock *sync.RWMutex
	fresh        int64
	msgHandler   iface.IMsgHandler
	txChan       chan *Msg
	ctx          context.Context
	cxtFunc      func()
	isClient     bool
	timer        *time.Timer
	isCmd        bool
}

func newConnectionAsServer(conn *gstcp.Conn, id string, handler iface.IMsgHandler) iface.IConnection {
	c := &Connection{
		conn:         conn,
		connID:       id,
		fresh:        0,
		msgHandler:   handler,
		propertyLock: &sync.RWMutex{},
		property:     make(map[string]interface{}),
		txChan:       make(chan *Msg, 100),
		isClient:     false,
		isCmd:        false,
	}
	c.ctx, c.cxtFunc = context.WithCancel(context.TODO())
	return c
}

func NewConnectionAsClient(conn *gstcp.Conn, id string, handler iface.IMsgHandler, isCmd bool) iface.IConnection {
	c := &Connection{
		conn:         conn,
		connID:       id,
		fresh:        0,
		msgHandler:   handler,
		propertyLock: &sync.RWMutex{},
		property:     make(map[string]interface{}),
		txChan:       make(chan *Msg, 100),
		isClient:     true,
		isCmd:        isCmd,
		timer:        time.NewTimer(time.Minute),
	}
	c.ctx, c.cxtFunc = context.WithCancel(context.TODO())
	return c
}

func (c *Connection) GetConnectionID() string {
	return c.connID
}

func (c *Connection) SetFresh(i int64) {
	atomic.StoreInt64(&c.fresh, i)
}

func (c *Connection) GetFresh() int64 {
	return atomic.LoadInt64(&c.fresh)
}

func (c *Connection) SetProperty(name string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[name] = value
}

func (c *Connection) GetProperty(name string) (v interface{}, ok bool) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	v, ok = c.property[name]
	return

}

func (c *Connection) DelProperty(name string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if _, ok := c.property[name]; ok {
		delete(c.property, name)
	}
}

func (c *Connection) startReader() error {
	for {
		if err := c.read(); err != nil {
			fmt.Println(err)
			break
		}
	}
	fmt.Printf("conn key : [%v] reader exit \n", c.connID)
	c.Stop()
	return nil
}

func (c *Connection) read() error {
	tp, err := c.conn.Recv(1)
	if err != nil {
		return err
	}
	switch tp[0] {
	case Fresh:
		c.fresh = time.Now().Unix()
	case Lv:
		binMsg, err := c.conn.RecvPkg(gstcp.PkgOption{
			HeaderSize: 4,
		})
		if err != nil {
			return err
		}
		msg := NewMsg(tp[0], binMsg)
		req := NewRequest(c, msg)
		if err := c.msgHandler.HandleCmdChan(req); err != nil {
			return err
		}
	default:
		return gserror.Newf("undefined tp %v \n", tp[0])
	}

	return nil
}

func (c *Connection) startWriter() (err error) {
	for {
		if err = c.write(); err != nil {
			break
		}
	}
	fmt.Printf("conn key : [%v] writer exit \n", c.connID)
	c.Stop()
	return err
}

func (c *Connection) write() error {
	select {
	// read txChan
	case msg := <-c.txChan:
		if msg.Tp == Fresh {
			if _, err := c.conn.Write([]byte{Fresh}); err != nil {
				return err
			}
		} else if msg.Tp == Lv {
			if _, err := c.conn.Write([]byte{Lv}); err != nil {
				return err
			}
			if err := c.conn.SendPkg(msg.BinData, gstcp.PkgOption{
				HeaderSize: 4,
			}); err != nil {
				return err
			}
		}

	case <-c.ctx.Done():
		return errors.New("cxt exit")
	}

	return nil
}

func (c *Connection) startWriteFresh() (err error) {
	for {
		if err = c.writeFresh(); err != nil {
			break
		}
	}
	c.Stop()
	return
}
func (c *Connection) writeFresh() error {
	select {
	case <-c.timer.C:
		c.SendFreshMsg()
	case <-c.ctx.Done():
		return errors.New("timer exit")
	}

	return nil
}

func (c *Connection) Start() {

	fmt.Printf("[new conn ,id [%v] , local [%v] , remote [%v] ] \n", c.connID, c.conn.LocalAddr().String(), c.conn.RemoteAddr().String())
	c.msgHandler.HandleOnConnect(c)
	defer c.msgHandler.HandleOffConnect(c)
	if c.isClient {
		if err := c.conn.Send(internal.ProtocolHeader); err != nil {
			return
		}
		c.conn = gstcp.UpgradeConnAsClient(c.conn)
		chanTp := CmdChan
		if !c.isCmd {
			chanTp = DataChan
		}

		if _, err := c.conn.Write([]byte{chanTp}); err != nil {
			return
		}
		switch chanTp {
		case CmdChan:
			c.runAsCmdChan()
		case DataChan:
			c.runAsDataChan()
		default:
			fmt.Println("error chan type")
			return
		}
	} else {
		header, err := c.conn.Recv(len(internal.ProtocolHeader))
		if err != nil {
			return
		}

		if !bytes.Equal(header, internal.ProtocolHeader) {
			return
		}
		c.conn = gstcp.UpgradeConnAsServer(c.conn)
		chanTp, err := c.conn.Recv(1)
		if err != nil {
			return
		}

		switch chanTp[0] {
		case CmdChan:
			c.isCmd = true
			c.SetFresh(time.Now().Unix())
			c.runAsCmdChan()
		case DataChan:
			c.isCmd = false
			c.runAsDataChan()
		default:
			fmt.Println("error chan type")
			return
		}
	}

}

func (c *Connection) runAsCmdChan() {
	go func() {
		if err := c.startWriter(); err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		if err := c.startReader(); err != nil {
			fmt.Println(err)
		}
	}()

	if c.isClient && c.isCmd {
		go func() {
			if err := c.startWriteFresh(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	<-c.ctx.Done()
}

func (c *Connection) runAsDataChan() {
	if err := c.msgHandler.HandleDataChan(c); err != nil {
		fmt.Println(err)
	}
}

func (c *Connection) Stop() {
	c.cxtFunc()
	if c.isClient {
		c.timer.Stop()
	}
	c.conn.Close()
}

func (c *Connection) SendLvMsg(data []byte) {
	c.txChan <- NewMsg(Lv, data)
}

func (c *Connection) SendFreshMsg() {
	c.txChan <- NewMsg(Fresh, nil)
}

func (c *Connection) IsCmdChan() bool {
	return c.isCmd
}

func (c *Connection) GetRawConn() *gstcp.Conn {
	return c.conn
}

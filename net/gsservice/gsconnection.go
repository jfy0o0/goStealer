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
	cxt          context.Context
	cxtFunc      func()
	wg           sync.WaitGroup
	isClient     bool
	timer        *time.Timer
}

func NewConnection(cxt context.Context, cxtFunc func(), conn *gstcp.Conn, id string, handler iface.IMsgHandler, isClient bool) iface.IConnection {
	c := &Connection{
		conn:         conn,
		connID:       id,
		fresh:        0,
		msgHandler:   handler,
		propertyLock: &sync.RWMutex{},
		property:     make(map[string]interface{}),
		txChan:       make(chan *Msg, 100),
		cxt:          cxt,
		cxtFunc:      cxtFunc,
		isClient:     isClient,
	}

	if c.isClient {
		c.timer = time.NewTimer(time.Minute)
	}

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

func (c *Connection) StartReader() error {
	defer c.wg.Done()

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
	case Json:
		binMsg, err := c.conn.RecvPkg(gstcp.PkgOption{
			HeaderSize: 4,
		})
		if err != nil {
			return err
		}
		msg := NewMsg(tp[0], binMsg)
		req := NewRequest(c, msg)
		if err := c.msgHandler.Handle(req); err != nil {
			return err
		}
	default:
		return gserror.Newf("undefined tp %v \n", tp[0])
	}

	return nil
}

func (c *Connection) StartWriter() (err error) {
	defer c.wg.Done()
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
		} else if msg.Tp == Json {
			if _, err := c.conn.Write([]byte{Json}); err != nil {
				return err
			}
			if err := c.conn.SendPkg(msg.BinData, gstcp.PkgOption{
				HeaderSize: 4,
			}); err != nil {
				return err
			}
		}

	case <-c.cxt.Done():
		return errors.New("cxt exit")
	}

	return nil
}

func (c *Connection) StartWriteFresh() (err error) {
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
		if err := c.SendMsg(Fresh, nil); err != nil {
			return err
		}

	case <-c.cxt.Done():
		return errors.New("timer exit")
	}

	return nil
}
func (c *Connection) Start() {

	fmt.Printf("[new conn ,id [%v] , local [%v] , remote [%v] ] \n", c.connID, c.conn.LocalAddr().String(), c.conn.RemoteAddr().String())

	if c.isClient {
		if err := c.conn.Send(internal.ProtocolHeader); err != nil {
			return
		}
		c.conn = gstcp.UpgradeConnAsClient(c.conn)
	} else {
		header, err := c.conn.Recv(len(internal.ProtocolHeader))
		if err != nil {
			return
		}

		if !bytes.Equal(header, internal.ProtocolHeader) {
			return
		}
		c.conn = gstcp.UpgradeConnAsServer(c.conn)
	}
	c.wg.Add(2)

	go func() {
		if err := c.StartWriter(); err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		if err := c.StartReader(); err != nil {
			fmt.Println(err)
		}
	}()
	if c.isClient {
		go func() {
			if err := c.StartWriteFresh(); err != nil {
				fmt.Println(err)
			}
		}()
	}

	c.wg.Wait()

}

func (c *Connection) Stop() {
	c.cxtFunc()
	if c.isClient {
		c.timer.Stop()
	}
	c.conn.Close()
}

func (c *Connection) SendMsg(tp byte, data []byte) error {
	c.txChan <- NewMsg(tp, data)
	return nil
}

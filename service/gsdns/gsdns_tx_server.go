package gsdns

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/jfy0o0/goStealer/container/gsmap"
	"github.com/jfy0o0/goStealer/net/gsudp"
	"github.com/jfy0o0/goStealer/os/gstimer"
	"github.com/jfy0o0/goStealer/util/gscall"
	"github.com/jfy0o0/goStealer/util/gsid_pool"
	"golang.org/x/net/dns/dnsmessage"
	"time"
)

type TxServer[T any] struct {
	requestMap *gsmap.AnyAnyMap[uint16, *gscall.Command[*UserDnsRequest[T], *UserDnsResponse[T]]]
	txChan     chan *UserDnsRequest[T]
	rxChan     chan<- *UserDnsResponse[T]
	realConn   *gsudp.Conn
	idPool     *gsid_pool.IDPool
	timer      *gstimer.Timer
	cancel     context.CancelFunc
	*TxServerConfig[T]
}

func NewTxServer[T any](rxChan chan<- *UserDnsResponse[T], config ...*TxServerConfig[T]) *TxServer[T] {

	s := &TxServer[T]{
		requestMap: gsmap.NewAnyAnyMap[uint16, *gscall.Command[*UserDnsRequest[T], *UserDnsResponse[T]]](true),
		txChan:     make(chan *UserDnsRequest[T], 10*10000),
		rxChan:     rxChan,
		idPool:     gsid_pool.New(1 << 16),
		timer:      gstimer.New(),
	}

	if len(config) > 0 {
		s.TxServerConfig = config[0]
	} else {
		s.TxServerConfig = GetDefaultTxServerConfig[T]()
	}

	s.timer.Stop()

	s.timer.AddSingleton(time.Duration(s.ScanInterval)*time.Second, s.onTimer)

	fmt.Printf(" tx server | id : [%v] ,to : [%v] , cap : [%v] ,scan [%v] ,appendHeader [%v] \n",
		s.Id, s.ToDnsIP, s.idPool.Cap(), s.ScanInterval, s.AppendHeader)
	return s
}

func (s *TxServer[T]) Run() error {
	s.timer.Start()
	realConn, err := gsudp.NewConn(s.ToDnsIP)
	if err != nil {
		return err
	}
	defer realConn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.realConn = realConn

	s.cancel = cancel

	go s.runTx(ctx)

	go s.runRx()

	<-ctx.Done()

	return nil
}

func (s *TxServer[T]) onTimer() {
	now := time.Now().Unix()
	s.requestMap.LockFunc(func(m map[uint16]*gscall.Command[*UserDnsRequest[T], *UserDnsResponse[T]]) {
		for k, v := range m {
			if v.GetValue().endTime < now {
				if v.GetValue().count.Val() == 0 {
					if s.OnNoResponse != nil {
						s.OnNoResponse(v.GetValue())
					}
				}
				delete(m, k)
				s.idPool.DeleteID(uint64(k))
			}
		}
	})
}

func (s *TxServer[T]) runTx(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.realConn.SetDeadline(time.Now())
			return
		case v := <-s.txChan:

			newID := s.idPool.NewID()

			cmd := newCommand(newID, v)

			cmd.GetValue().msg.ID = uint16(newID)

			data, err := s.pack(cmd)
			if err != nil {
				s.idPool.DeleteID(newID)
				continue
			}

			s.realConn.Write(data)

			s.requestMap.Set(uint16(newID), cmd)
		}
	}
}

func (s *TxServer[T]) runRx() {
	buf := make([]byte, 2048)

	for {
		n, srcAddr, err := s.realConn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		if srcAddr.String() != s.ToDnsIP {
			continue
		}

		msg := &dnsmessage.Message{}
		if err := msg.Unpack(buf[:n]); err != nil {
			continue
		}

		cmd := s.requestMap.Get(msg.ID)
		if cmd == nil {
			continue
		}
		cmd.GetValue().count.Add(1)
		msg.ID = cmd.GetValue().oldDnsID
		s.rxChan <- &UserDnsResponse[T]{
			msg:   msg,
			addr:  cmd.GetValue().addr,
			value: cmd.GetValue().value,
		}

	}
}

func (s *TxServer[T]) SendRequest(req *UserDnsRequest[T]) {
	s.txChan <- req
}

func (s *TxServer[T]) Stop() {
	s.cancel()
}

func (s *TxServer[T]) pack(cmd *gscall.Command[*UserDnsRequest[T], *UserDnsResponse[T]]) (d []byte, err error) {
	data, err := cmd.GetValue().msg.Pack()
	if err != nil {
		return d, err
	}

	if !s.AppendHeader {
		return data, nil
	}

	buffer := bytes.NewBuffer(userProtocolHeader)
	if cmd.GetValue().srcIP == 0 && cmd.GetValue().springBoardIP == 0 {
		buffer.Write(cmd.GetValue().addr.IP.To4())
		binary.Write(buffer, binary.BigEndian, &s.LocalIP)
	} else {
		binary.Write(buffer, binary.BigEndian, &cmd.GetValue().srcIP)
		binary.Write(buffer, binary.BigEndian, &cmd.GetValue().springBoardIP)
	}
	buffer.Write(data)

	return buffer.Bytes(), nil

}

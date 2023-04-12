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
	id           int
	requestMap   *gsmap.AnyAnyMap[uint16, *gscall.Command[*UserDnsRequest[T], *UserDnsResponse[T]]]
	txChan       chan *UserDnsRequest[T]
	rxChan       chan<- *UserDnsResponse[T]
	toDnsIP      string
	realConn     *gsudp.Conn
	idPool       *gsid_pool.IDPool
	timer        *gstimer.Timer
	cancel       context.CancelFunc
	appendHeader bool
	localIP      uint32
	adapter      DnsAdapter[T]
}

func NewTxServer[T any](id int, toDnsIP string, rxChan chan<- *UserDnsResponse[T], scanInterval int, appendHeader bool, localIP uint32, adapter DnsAdapter[T]) *TxServer[T] {

	s := &TxServer[T]{
		id:           id,
		requestMap:   gsmap.NewAnyAnyMap[uint16, *gscall.Command[*UserDnsRequest[T], *UserDnsResponse[T]]](true),
		txChan:       make(chan *UserDnsRequest[T], 10*10000),
		rxChan:       rxChan,
		toDnsIP:      toDnsIP + ":53",
		idPool:       gsid_pool.New(1 << 16),
		timer:        gstimer.New(),
		appendHeader: appendHeader,
		localIP:      localIP,
		adapter:      adapter,
	}
	s.timer.Stop()

	s.timer.AddSingleton(time.Duration(scanInterval)*time.Second, s.onTimer)

	fmt.Sprintf(" tx server | id : [%v] , cap : [%v] ,scan [%v] ,appendHeader [%v] ", s.id, s.idPool.Cap(), scanInterval, s.appendHeader)
	return s
}

func (s *TxServer[T]) Run() error {
	s.timer.Start()
	realConn, err := gsudp.NewConn(s.toDnsIP)
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
					s.adapter.OnNoResponse(v.GetValue())
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

		if srcAddr.String() != s.toDnsIP {
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

	if !s.appendHeader {
		return data, nil
	}

	buffer := bytes.NewBuffer(userProtocolHeader)
	if cmd.GetValue().srcIP == 0 && cmd.GetValue().springBoardIP == 0 {
		buffer.Write(cmd.GetValue().addr.IP.To4())
		binary.Write(buffer, binary.BigEndian, &s.localIP)
	} else {
		binary.Write(buffer, binary.BigEndian, &cmd.GetValue().srcIP)
		binary.Write(buffer, binary.BigEndian, &cmd.GetValue().springBoardIP)
	}
	buffer.Write(data)

	return buffer.Bytes(), nil

}

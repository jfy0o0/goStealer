package main

import (
	"context"
	"fmt"
	"github.com/jfy0o0/goStealer/net/gsservice"
	"github.com/jfy0o0/goStealer/net/gsservice/iface"
	"github.com/jfy0o0/goStealer/net/gstcp"
	"log"
	"runtime"
	"time"
)

func test(skip int) {
	call(skip)
}

func call(skip int) {
	pc, file, line, ok := runtime.Caller(skip)
	pcName := runtime.FuncForPC(pc).Name()
	fmt.Println(fmt.Sprintf("%v   %s   %d   %t   %s", pc, file, line, ok, pcName))
}

type MsgClientHandler struct {
}

func (mh MsgClientHandler) HandleCmdChan(message iface.IRequest) error {
	fmt.Println("client get message :", string(message.GetMsg().GetMsgData()))
	return nil
}
func (mh MsgClientHandler) HandleDataChan() {
	fmt.Println("client handle data message :")

}

type MsgServerHandler struct {
}

func (mh MsgServerHandler) HandleCmdChan(message iface.IRequest) error {
	fmt.Println("server get message :", string(message.GetMsg().GetMsgData()))
	return nil
}
func (mh MsgServerHandler) HandleDataChan() {
	fmt.Println("server handle data message :")
}
func main() {
	handler := MsgServerHandler{}
	s := gsservice.NewServer("0.0.0.0", 1367, gsservice.NewConnMgr(), handler)
	go s.Start()

	time.Sleep(time.Second * 2)

	c, err := gstcp.NewConn("127.0.0.1:1367")
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	h := MsgClientHandler{}
	dealConn := gsservice.NewConnectionAsClient(ctx, cancel, c, "1", h, true)
	go dealConn.Start()
	time.Sleep(time.Second * 2)
	dealConn.SendMsg(gsservice.Json, []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))

	go func() {
		for i := 0; ; i++ {
			s.ConnMgr.Walk(func(m map[string]iface.IConnection) {
				fmt.Println("map size :", len(m))
				for _, v := range m {
					v.SendMsg(gsservice.Json, []byte(fmt.Sprintf("hello , i am server %v ", i)))
				}
			})
			time.Sleep(time.Second * 2)
		}

	}()

	//gsservice.NewConnection()

	//c.Write([]byte{0x01, 0x02, 0x04, 0x08, 0x08, 0x04, 0x02, 0x01})
	//c2 := gstcp.UpgradeConnAsClient(c)
	//a := &AAAAA{
	//	Name: "aaaaaaa",
	//	Age:  123,
	//}
	//d, err := json.Marshal(a)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//m := gsservice.NewMsg(Json, d)
	//
	//d2, err := json.Marshal(m)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(string(d2))
	//c2.Write([]byte{Json})
	//err = c2.SendPkg(d2, gstcp.PkgOption{
	//	HeaderSize: 4,
	//})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//c.Close()

	//select {}
	for {
		time.Sleep(time.Minute)
	}
}

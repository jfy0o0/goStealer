package gstcp

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

func process(conn *Conn) {
	defer conn.Close() // 关闭连接

	hello := NewGsHello[int](0, 1, 1)
	clientHello, err := hello.HandShakeAsServer(conn)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(clientHello)
}
func runServer() {
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(NewConnByNetConn(conn)) // 启动一个goroutine处理连接
	}
}
func runClient() {
	c, err := NewConn("127.0.0.1:20000")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	hello := NewGsHello[int](1, 2, 2)
	serverHello, err := hello.HandShakeAsClient(c)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println(serverHello)
}
func TestGsHello_HandShake(t *testing.T) {
	go runServer()
	time.Sleep(time.Second)

	runClient()

	time.Sleep(time.Second)
}

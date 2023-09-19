package gsnet

import (
	"github.com/jfy0o0/goStealer/net/gstcp"
	"log"
	"strconv"
	"testing"
	"time"
)

type TestServerAdapter struct {
}

func (s *TestServerAdapter) OnHeartBeatSuccessful(string) {
	log.Println("OnHeartBeatSuccessful")
}
func (s *TestServerAdapter) OnHeartBeatFailed(string) {
	log.Println("OnHeartBeatFailed")
}
func (s *TestServerAdapter) OnConnectedStart(conn *gstcp.Conn) error {
	log.Println("OnConnectedStart")
	return nil
}
func (s *TestServerAdapter) OnConnectedHandClientHello(hello *gstcp.GsHello[HelloExtend[string]]) error {
	log.Println("OnConnectedHandClientHello")
	return nil
}

type TestClientAdapter struct {
}

func (c *TestClientAdapter) OnConnectedStart(conn *gstcp.Conn) error {
	log.Println("client OnConnectedStart")
	return nil
}
func (c *TestClientAdapter) OnConnectedHandServerHello(hello *gstcp.GsHello[HelloExtend[string]]) error {
	log.Println("client OnConnectedHandServerHello")
	return nil
}

type TestSessionAdapter struct {
}

func (s *TestSessionAdapter) OnMsg(conn *gstcp.Conn) {
	for {
		x, err := conn.RecvPkg(gstcp.PkgOption{
			HeaderSize: 4,
		})
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("OnMsg", string(x))
	}

}
func (s *TestSessionAdapter) OnSendMsg(conn *gstcp.Conn, msg interface{}) {
	log.Println("OnSendMsg", msg)
	conn.SendPkg([]byte(msg.(string)), gstcp.PkgOption{
		HeaderSize:  4,
		MaxDataSize: 0,
		Retry:       gstcp.Retry{},
	})
}

func TestGsNetServer(t *testing.T) {
	serverAdapter := &TestServerAdapter{}
	sessionAdapter := &TestSessionAdapter{}
	s := NewServer[string](&ServerConfig[string]{
		ListenAddr:             "0.0.0.0:19999",
		CheckBeatHeartInterval: time.Minute * 5,
		Hello:                  gstcp.NewGsHello[HelloExtend[string]](1, 1, HelloExtend[string]{Key: "asd", V: "abc"}),
		ServerAdapter:          serverAdapter,
		SessionAdapter:         sessionAdapter,
		SessionConf:            SessionConfig{CommunicationType: CommunicationTypeUserDefined},
	})

	s.Run()

}

func TestGsNetClient(t *testing.T) {
	clientAdapter := &TestClientAdapter{}
	sessionAdapter := &TestSessionAdapter{}

	s := NewClient[string](&ClientConfig[string]{
		ConnAddr:       "0.0.0.0:19999",
		Hello:          gstcp.NewGsHello[HelloExtend[string]](1, 1, HelloExtend[string]{Key: "asd", V: "abc"}),
		ClientAdapter:  clientAdapter,
		SessionAdapter: sessionAdapter,
		SessionConf:    SessionConfig{CommunicationType: CommunicationTypeUserDefined},
	})

	go func() {
		for i := 0; ; i++ {
			time.Sleep(time.Second)
			s.Session.Push("hi" + strconv.FormatInt(int64(i), 10))
		}
	}()

	s.Run()

}

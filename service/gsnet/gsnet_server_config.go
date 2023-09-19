package gsnet

import (
	"github.com/jfy0o0/goStealer/net/gstcp"
	"strconv"
	"time"
)

type ServerAdapter[T any] interface {
	OnHeartBeatSuccessful(string)
	OnHeartBeatFailed(string)
	OnConnectedStart(conn *gstcp.Conn) error
	OnConnectedHandClientHello(hello *gstcp.GsHello[HelloExtend[T]]) error
}

type ServerConfig[T any] struct {
	ListenAddr             string
	CheckBeatHeartInterval time.Duration
	Hello                  *gstcp.GsHello[HelloExtend[T]]
	ServerAdapter[T]
	SessionAdapter[T]
	SessionConf SessionConfig
}

func GetDefaultServerConfig[T any]() *ServerConfig[T] {
	return &ServerConfig[T]{
		ListenAddr: "0.0.0.0:12321",
		Hello: gstcp.NewGsHello[HelloExtend[T]](1, 1, HelloExtend[T]{
			Key: strconv.FormatInt(time.Now().UnixNano(), 10),
		}),
		CheckBeatHeartInterval: 5 * time.Minute,
		SessionConf:            SessionConfig{CommunicationType: CommunicationTypeYamuxMuti},
	}
}

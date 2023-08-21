package gsnet

import (
	"github.com/jfy0o0/goStealer/net/gstcp"
	"strconv"
	"time"
)

type ClientAdapter[T any] interface {
	OnConnectedStart(conn *gstcp.Conn) error
	OnConnectedHandServerHello(hello *gstcp.GsHello[HelloExtend[T]]) error
}
type ClientConfig[T any] struct {
	ConnAddr string
	Hello    *gstcp.GsHello[HelloExtend[T]]
	ClientAdapter[T]
	SessionAdapter[T]
}

func GetDefaultClientConfig[T any]() *ClientConfig[T] {
	return &ClientConfig[T]{
		ConnAddr: "0.0.0.0:12321",
		Hello: gstcp.NewGsHello[HelloExtend[T]](1, 1, HelloExtend[T]{
			Key: strconv.FormatInt(time.Now().UnixNano(), 10),
		}),
	}
}

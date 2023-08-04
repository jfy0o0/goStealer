package gstcp

import (
	"bytes"
	"encoding/json"
	"github.com/jfy0o0/goStealer/errors/gscode"
	"github.com/jfy0o0/goStealer/errors/gserror"
	"time"
)

var helloMsg = []byte{
	0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
	0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA,
}

// MajorVersion: client == server
// StartTime : used time when process start
type GsHello[T any] struct {
	MajorVersion int
	MinorVersion int
	StartTime    int64
	Data         T
}

func NewGsHello[T any](majorVersion, minorVersion int, data T) *GsHello[T] {
	return &GsHello[T]{
		MajorVersion: majorVersion,
		MinorVersion: minorVersion,
		StartTime:    time.Now().Unix(),
		Data:         data,
	}
}

func (g *GsHello[T]) HandShakeAsClient(conn *Conn) (serverHello *GsHello[T], err error) {
	if _, err = conn.Write(helloMsg); err != nil {
		return
	}

	helloData, err := conn.Recv(len(helloMsg))
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(helloMsg, helloData) {
		return nil, gserror.NewCode(gscode.CodeNotAuthorized)
	}

	jsonData, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}

	if err = conn.SendPkg(jsonData, PkgOption{
		HeaderSize: 4,
	}); err != nil {
		return nil, err
	}

	result, err := conn.RecvPkg(PkgOption{
		HeaderSize: 4,
	})
	if err != nil {
		return nil, err
	}
	x := &GsHello[T]{}
	if err = json.Unmarshal(result, x); err != nil {
		return nil, err
	}

	if x.MajorVersion != g.MajorVersion {
		return nil, gserror.New("Version Incompatible")
	}
	return x, nil
}
func (g *GsHello[T]) HandShakeAsServer(conn *Conn) (clientHello *GsHello[T], err error) {
	helloData, err := conn.Recv(len(helloMsg))
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(helloMsg, helloData) {
		return nil, gserror.NewCode(gscode.CodeNotAuthorized)
	}
	if _, err = conn.Write(helloMsg); err != nil {
		return
	}
	result, err := conn.RecvPkg(PkgOption{
		HeaderSize: 4,
	})
	if err != nil {
		return nil, err
	}
	x := &GsHello[T]{}
	if err = json.Unmarshal(result, x); err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}

	if err = conn.SendPkg(jsonData, PkgOption{
		HeaderSize: 4,
	}); err != nil {
		return nil, err
	}

	if x.MajorVersion != g.MajorVersion {
		return nil, gserror.New("Version Incompatible")
	}

	return x, nil
}

package internal

import (
	"io"
	"net"
	"time"
)

func Relay(left, right net.Conn) {
	ch := make(chan bool)

	go func() {
		_, _ = io.Copy(right, left)
		_ = right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
		_ = left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
		ch <- true
	}()

	_, _ = io.Copy(left, right)
	_ = right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
	_ = left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
	<-ch

}

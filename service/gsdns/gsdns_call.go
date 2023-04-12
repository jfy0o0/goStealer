package gsdns

import "github.com/jfy0o0/goStealer/util/gscall"

func newCommand[T any](id uint64, request *UserDnsRequest[T]) *gscall.Command[*UserDnsRequest[T], *UserDnsResponse[T]] {
	return gscall.New[*UserDnsRequest[T], *UserDnsResponse[T]](id, request)
}

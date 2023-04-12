package gsdns

type DnsAdapter[T any] interface {
	OnNoResponse(*UserDnsRequest[T])
}

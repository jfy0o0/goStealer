package iface

type IConnectionManager interface {
	Add(connection IConnection)

	Del(connection IConnection) error

	Get(string) (IConnection, bool)

	Clear()

	Len() int

	Walk(func(map[string]IConnection))
}

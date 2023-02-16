package gsservice

import (
	"errors"
	"fmt"
	"github.com/jfy0o0/goStealer/net/gsservice/iface"
	"sync"
)

type ConnectionManager struct {
	connections map[string]iface.IConnection
	mu          sync.RWMutex
}

func NewConnMgr() iface.IConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]iface.IConnection),
		mu:          sync.RWMutex{},
	}
}

func (cm *ConnectionManager) Walk(f func(map[string]iface.IConnection)) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	f(cm.connections)
}

// add connection
func (cm *ConnectionManager) Add(connection iface.IConnection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.connections[connection.GetConnectionID()] = connection
}

// del connection
func (cm *ConnectionManager) Del(connection iface.IConnection) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	c, ok := cm.connections[connection.GetConnectionID()]
	if !ok {
		fmt.Println("connection", connection.GetConnectionID(), "do not exits!")
		return errors.New("remove nil pointer connection")
	}
	delete(cm.connections, connection.GetConnectionID())
	c.Stop()
	return nil
}

// get conn by id
func (cm *ConnectionManager) Get(id string) (conn iface.IConnection, ok bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	conn, ok = cm.connections[id]
	return
}

// clear connection
func (cm *ConnectionManager) Clear() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for id, conn := range cm.connections {
		// stop conn
		//conn.StopWithNotConnMgr()
		conn.Stop()

		// del conn
		delete(cm.connections, id)
	}
}

// get len
func (cm *ConnectionManager) Len() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return len(cm.connections)
}

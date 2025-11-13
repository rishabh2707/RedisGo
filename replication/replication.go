package replication

import (
	"net"
)

type SlaveConnections struct {
	connections []net.Conn
}

type ReplicaOffset struct {
	offset int64
}

var replicaOffset *ReplicaOffset

func InitReplicaOffset() {
	replicaOffset = &ReplicaOffset{offset: 0}
}

func SetReplicaOffset(offset int64) {
	replicaOffset.offset = offset
}

func GetReplicaOffset() int64 {
	return replicaOffset.offset
}

var slaveConnections *SlaveConnections

func InitSlaveConnections() {
	slaveConnections = &SlaveConnections{connections: make([]net.Conn, 0)}
}

func AddSlaveConnection(conn *net.Conn) {
	slaveConnections.connections = append(slaveConnections.connections, *conn)
}

func RemoveSlaveConnection(conn *net.Conn) {
	for i, c := range slaveConnections.connections {
		if c == *conn {
			slaveConnections.connections = append(slaveConnections.connections[:i], slaveConnections.connections[i+1:]...)
		}
	}
}

func GetSlaveConnections() []net.Conn {
	return slaveConnections.connections
}

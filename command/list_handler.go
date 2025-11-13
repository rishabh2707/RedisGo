package command

import (
	"net"
	"strconv"
	"time"

	"com.github.redisgo/config"
	"com.github.redisgo/database"
	"com.github.redisgo/replication"
	"com.github.redisgo/util"
)

type objectList struct {
	Value []string
}

func (cmd *Cmd) handleRPushCommand(conn *net.Conn) {
	if len(cmd.Args) < 3 {
		writeResponse(conn, "-ERR wrong number of arguments for 'rpush' command\r\n")
		return
	}
	if cmd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		return
	}
	key := cmd.Args[1]
	//value := cmd.Args[2]

	lock := database.GetKeyLock(key)

	lock.Lock()
	ObjectList, ok := database.Get(key)
	if !ok {
		ObjectList = &objectList{Value: cmd.Args[2:]}
	} else {
		ObjectList.(*objectList).Value = append(ObjectList.(*objectList).Value, cmd.Args[2:]...)
	}
	database.Store(key, ObjectList)
	lock.Unlock()
	if config.ServerConfig.Role == "master" {
		writeResponse(conn, util.ReturnIntegerResponse(len(ObjectList.(*objectList).Value)))
		go cmd.propagateToSlaves()
	}
	if config.ServerConfig.Role == "slave" {
		replication.SetReplicaOffset(replication.GetReplicaOffset() + int64(len(cmd.ToRespFormat())))
	}
}

func (cmd *Cmd) handleLPushCommand(conn *net.Conn) {
	if len(cmd.Args) < 3 {
		writeResponse(conn, "-ERR wrong number of arguments for 'lpush' command\r\n")
		return
	}

	if cmd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		return
	}

	key := cmd.Args[1]
	//value := cmd.Args[2]

	lock := database.GetKeyLock(key)
	defer lock.Unlock()

	lock.Lock()
	ObjectList, ok := database.Get(key)
	if !ok {
		ObjectList = &objectList{Value: []string{}}
	}
	for _, v := range cmd.Args[2:] {
		ObjectList.(*objectList).Value = append([]string{v}, ObjectList.(*objectList).Value...)
	}
	database.Store(key, ObjectList)
	if config.ServerConfig.Role == "master" {
		writeResponse(conn, util.ReturnIntegerResponse(len(ObjectList.(*objectList).Value)))
		go cmd.propagateToSlaves()
	}
	if config.ServerConfig.Role == "slave" {
		replication.SetReplicaOffset(replication.GetReplicaOffset() + int64(len(cmd.ToRespFormat())))
	}
}

func (cmd *Cmd) handleLLenCommand(conn *net.Conn) {
	if len(cmd.Args) < 2 {
		writeResponse(conn, "-ERR wrong number of arguments for 'llen' command\r\n")
		return
	}
	if cmd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		return
	}

	key := cmd.Args[1]

	lock := database.GetKeyLock(key)

	lock.Lock()
	ObjectList, ok := database.Get(key)
	if !ok {
		lock.Unlock()
		writeResponse(conn, util.ReturnIntegerResponse(0))
		return
	}
	lock.Unlock()
	writeResponse(conn, util.ReturnIntegerResponse(len(ObjectList.(*objectList).Value)))
}

func (cmd *Cmd) handleLPopCommand(conn *net.Conn) {
	if len(cmd.Args) < 2 {
		writeResponse(conn, "-ERR wrong number of arguments for 'lpop' command\r\n")
		return
	}
	if cmd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		return
	}

	key := cmd.Args[1]

	lock := database.GetKeyLock(key)

	lock.Lock()
	ObjectList, ok := database.Get(key)
	if !ok || len(ObjectList.(*objectList).Value) == 0 {
		lock.Unlock()
		writeResponse(conn, util.ReturnNullResponse())
		return
	}
	poppedValue := ObjectList.(*objectList).Value[0]
	ObjectList.(*objectList).Value = ObjectList.(*objectList).Value[1:]
	database.Store(key, ObjectList)
	lock.Unlock()
	if config.ServerConfig.Role == "master" {
		writeResponse(conn, util.ReturnBulkStringResponse(poppedValue))
		go cmd.propagateToSlaves()
	}
	if config.ServerConfig.Role == "slave" {
		replication.SetReplicaOffset(replication.GetReplicaOffset() + int64(len(cmd.ToRespFormat())))
	}
}

func (cmd *Cmd) handleBLPopCommand(conn *net.Conn) {
	if len(cmd.Args) < 3 {
		writeResponse(conn, "-ERR wrong number of arguments for 'blpop' command\r\n")
		return
	}
	if cmd.checkMultiExists() {
		writeResponse(conn, util.ReturnQueuedResponse())
		go cmd.propagateToSlaves()
		return
	}
	key := cmd.Args[1]
	timeout, err := strconv.ParseFloat(cmd.Args[2], 64)
	if err != nil {
		writeResponse(conn, "-ERR invalid timeout\r\n")
		return
	}
	timeoutInMilliseconds := float64(timeout * 1000)
	lock := database.GetKeyLock(key)

	if timeoutInMilliseconds != float64(0) {
		time.Sleep(time.Duration(timeoutInMilliseconds) * time.Millisecond)
		lock.Lock()
		ObjectList, ok := database.Get(key)
		if !ok || len(ObjectList.(*objectList).Value) == 0 {
			lock.Unlock()
			writeResponse(conn, util.ReturnNullResponse())
			return
		}
		poppedValue := ObjectList.(*objectList).Value[0]
		ObjectList.(*objectList).Value = ObjectList.(*objectList).Value[1:]
		database.Store(key, ObjectList)
		lock.Unlock()
		if config.ServerConfig.Role == "master" {
			writeResponse(conn, util.ReturnBulkStringResponse(poppedValue))
			go cmd.propagateToSlaves()
		}
		if config.ServerConfig.Role == "slave" {
			replication.SetReplicaOffset(replication.GetReplicaOffset() + int64(len(cmd.ToRespFormat())))
		}
		return
	}

	for {
		lock.Lock()
		ObjectList, ok := database.Get(key)
		if !ok {
			lock.Unlock()
			continue
		}

		if len(ObjectList.(*objectList).Value) > 0 {
			poppedValue := ObjectList.(*objectList).Value[0]
			ObjectList.(*objectList).Value = ObjectList.(*objectList).Value[1:]
			database.Store(key, ObjectList)
			lock.Unlock()
			if config.ServerConfig.Role == "master" {
				writeResponse(conn, util.ReturnBulkStringResponse(poppedValue))
				go cmd.propagateToSlaves()
			}
			if config.ServerConfig.Role == "slave" {
				replication.SetReplicaOffset(replication.GetReplicaOffset() + int64(len(cmd.ToRespFormat())))
			}
			return
		}
		lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
}

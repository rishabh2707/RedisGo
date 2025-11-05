package command

import (
	"strconv"
	"time"

	"com.github.redisgo/database"
	"com.github.redisgo/util"
)

type objectList struct {
	Value []string
}

func (cmd *Cmd) handleRPushCommand() string {
	if len(cmd.Args) < 3 {
		return "-ERR wrong number of arguments for 'rpush' command\r\n"
	}
	if cmd.checkMultiExists() {
		return util.ReturnQueuedResponse()
	}
	key := cmd.Args[1]
	//value := cmd.Args[2]

	lock := database.GetKeyLock(key)
	defer lock.Unlock()

	lock.Lock()
	ObjectList, ok := database.Get(key)
	if !ok {
		ObjectList = &objectList{Value: cmd.Args[2:]}
	} else {
		ObjectList.(*objectList).Value = append(ObjectList.(*objectList).Value, cmd.Args[2:]...)
	}
	database.Store(key, ObjectList)
	return util.ReturnIntegerResponse(len(ObjectList.(*objectList).Value))
}

func (cmd *Cmd) handleLPushCommand() string {
	if len(cmd.Args) < 3 {
		return "-ERR wrong number of arguments for 'lpush' command\r\n"
	}

	if cmd.checkMultiExists() {
		return util.ReturnQueuedResponse()
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
	return util.ReturnIntegerResponse(len(ObjectList.(*objectList).Value))
}

func (cmd *Cmd) handleLLenCommand() string {
	if len(cmd.Args) < 2 {
		return "-ERR wrong number of arguments for 'llen' command\r\n"
	}
	if cmd.checkMultiExists() {
		return util.ReturnQueuedResponse()
	}

	key := cmd.Args[1]

	lock := database.GetKeyLock(key)
	defer lock.Unlock()

	lock.Lock()
	ObjectList, ok := database.Get(key)
	if !ok {
		return util.ReturnIntegerResponse(0)
	}
	return util.ReturnIntegerResponse(len(ObjectList.(*objectList).Value))
}

func (cmd *Cmd) handleLPopCommand() string {
	if len(cmd.Args) < 2 {
		return "-ERR wrong number of arguments for 'lpop' command\r\n"
	}
	if cmd.checkMultiExists() {
		return util.ReturnQueuedResponse()
	}

	key := cmd.Args[1]

	lock := database.GetKeyLock(key)
	defer lock.Unlock()

	lock.Lock()
	ObjectList, ok := database.Get(key)
	if !ok {
		return util.ReturnNullResponse()
	}
	poppedValue := ObjectList.(*objectList).Value[0]
	ObjectList.(*objectList).Value = ObjectList.(*objectList).Value[1:]
	database.Store(key, ObjectList)
	return util.ParseNormalResponse(poppedValue)
}

func (cmd *Cmd) handleBLPopCommand() string {
	if len(cmd.Args) < 3 {
		return "-ERR wrong number of arguments for 'blpop' command\r\n"
	}
	if cmd.checkMultiExists() {
		return util.ReturnQueuedResponse()
	}
	key := cmd.Args[1]
	timeout, err := strconv.ParseFloat(cmd.Args[2], 64)
	if err != nil {
		return "-ERR invalid timeout\r\n"
	}
	timeoutInMilliseconds := float64(timeout * 1000)
	lock := database.GetKeyLock(key)
	defer lock.Unlock()

	if timeoutInMilliseconds != float64(0) {
		time.Sleep(time.Duration(timeoutInMilliseconds) * time.Millisecond)
		lock.Lock()
		ObjectList, ok := database.Get(key)
		if !ok {
			return util.ReturnNullResponse()
		}
		poppedValue := ObjectList.(*objectList).Value[0]
		ObjectList.(*objectList).Value = ObjectList.(*objectList).Value[1:]
		database.Store(key, ObjectList)
		return util.ParseNormalResponse(poppedValue)
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
			return util.ParseNormalResponse(poppedValue)
		}
		lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
}

package command

import (
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

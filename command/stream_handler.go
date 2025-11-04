package command

import (
	"com.github.redisgo/database"
	"com.github.redisgo/util"
)

type streamObject struct {
	Id    string
	Value map[string]string
}

type streamList struct {
	Value []*streamObject
}

func (cmnd *Cmd) handleXAddCommand() string {
	if len(cmnd.Args) < 5 {
		return "-ERR wrong number of arguments for 'xadd' command\r\n"
	}

	key := cmnd.Args[1]
	id := cmnd.Args[2]
	StreamObject := &streamObject{Id: id, Value: make(map[string]string)}
	for i := 3; i < len(cmnd.Args); i++ {
		StreamObject.Value[cmnd.Args[i]] = cmnd.Args[i+1]
		i++
	}

	lock := database.GetKeyLock(key)
	defer lock.Unlock()

	lock.Lock()
	StreamList, ok := database.Get(key)
	if !ok {
		StreamList = &streamList{Value: []*streamObject{StreamObject}}
	} else {
		StreamList.(*streamList).Value = append(StreamList.(*streamList).Value, StreamObject)
	}
	database.Store(key, StreamList)
	return util.ParseNormalResponse(StreamObject.Id)
}

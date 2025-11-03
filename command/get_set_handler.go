package command

import (
	"strconv"

	"com.github.redisgo/database"
)

type object struct {
	Value string
}

func (cmd *Cmd) handleSetCommand() string {
	if len(cmd.Args) < 3 {
		return "-ERR wrong number of arguments for 'set' command\r\n"
	}

	lock := database.GetKeyLock(cmd.Args[1])
	lock.Lock()
	defer lock.Unlock()

	key := cmd.Args[1]
	value := cmd.Args[2]

	Object := &object{Value: value}

	database.Store(key, Object)
	return "+OK\r\n"
}

func (cmd *Cmd) handleGetCommand() string {
	if len(cmd.Args) < 2 {
		return "-ERR wrong number of arguments for 'get' command\r\n"
	}

	lock := database.GetKeyLock(cmd.Args[1])
	lock.Lock()
	defer lock.Unlock()

	key := cmd.Args[1]
	value, ok := database.Get(key)
	if !ok {
		return "$-1\r\n"
	}

	Object := value.(*object)

	return "$" + strconv.Itoa(len(Object.Value)) + "\r\n" + Object.Value + "\r\n"
}

package command

import (
	"fmt"
	"strconv"
	"strings"

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
	previousId := ""
	if !ok {
		StreamList = &streamList{Value: []*streamObject{StreamObject}}
	} else {
		StreamList.(*streamList).Value = append(StreamList.(*streamList).Value, StreamObject)
		previousId = StreamList.(*streamList).Value[len(StreamList.(*streamList).Value)-1].Id
	}
	if strings.HasSuffix(id, "*") {
		firstPart := strings.Split(id, "-")[0]
		value, ok := database.Get(getSequenceNumberKey(key, firstPart))
		if !ok {
			database.Store(getSequenceNumberKey(key, firstPart), &object{Value: "0"})
			id = firstPart + "-0"
		} else {
			sequenceNumber, err := strconv.Atoi(value.(*object).Value)
			if err != nil {
				return util.ReturnErrorResponse("The ID specified in XADD must be a string")
			}
			sequenceNumber++
			id = firstPart + "-" + strconv.Itoa(sequenceNumber)
			database.Store(getSequenceNumberKey(key, firstPart), &object{Value: strconv.Itoa(sequenceNumber)})
		}
	}
	StreamObject.Id = id
	ok, err := validateStreamObjectId(previousId, StreamObject.Id)
	if !ok {
		return util.ReturnErrorResponse(err.Error())
	}
	database.Store(key, StreamList)
	return util.ParseNormalResponse(StreamObject.Id)
}

func validateStreamObjectId(previousId string, currentId string) (bool, error) {
	if strings.HasPrefix(currentId, "0-") {
		return false, fmt.Errorf("The ID specified in XADD must be greater than 0-0")
	}
	if !(strings.Compare(previousId, currentId) < 0) {
		return false, fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")
	}
	return true, nil
}

func getSequenceNumberKey(key string, firstPart string) string {
	return key + ":" + firstPart + ":sequenceNumber"
}

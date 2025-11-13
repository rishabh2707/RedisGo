package util

import (
	"strconv"

	"com.github.redisgo/config"
	"com.github.redisgo/replication"
)

func ReturnBulkStringResponse(response string) string {
	return "$" + strconv.Itoa(len(response)) + "\r\n" + response + "\r\n"
}

func ReturnNullResponse() string {
	return "$-1\r\n"
}

func ReturnOkResponse() string {
	return "+OK\r\n"
}

func ReturnQueuedResponse() string {
	return "+QUEUED\r\n"
}

func ReturnIntegerResponse(response int) string {
	return ":" + strconv.Itoa(response) + "\r\n"
}

func ReturnEmptyArrayResponse() string {
	return "*0\r\n"
}

func ReturnArrayResponse(response []string) string {
	result := "*" + strconv.Itoa(len(response)) + "\r\n"
	for _, v := range response {
		result += "$" + strconv.Itoa(len(v)) + "\r\n" + v + "\r\n"
	}
	return result
}

func ReturnErrorResponse(response string) string {
	return "-ERR " + response + "\r\n"
}

func ReturnReplicationInfoResponse() string {
	result := "#replication\n"
	result += "role:" + config.ServerConfig.Role + "\n"
	if config.ServerConfig.Role == "master" {
		result += "master_replid:" + config.ServerConfig.Master_replid + "\n"
		result += "master_repl_offset:" + strconv.FormatInt(config.ServerConfig.Master_repl_offset, 10) + "\n"
	}
	return ReturnBulkStringResponse(result)
}

func ReturnReplconfGetAckResponse() string {
	result := []string{"REPLCONF", "ACK", strconv.FormatInt(replication.GetReplicaOffset(), 10)}
	replication.SetReplicaOffset(replication.GetReplicaOffset())
	return ReturnArrayResponse(result)
}

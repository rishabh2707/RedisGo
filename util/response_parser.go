package util

import (
	"strconv"
)

func ParseNormalResponse(response string) string {
	return "$" + strconv.Itoa(len(response)) + "\r\n" + response + "\r\n"
}

func ReturnNullResponse() string {
	return "$-1\r\n"
}

func ReturnOkResponse() string {
	return "+OK\r\n"
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

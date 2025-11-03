package util

import "strconv"

func ParseNormalResponse(response string) string {
	return "$" + strconv.Itoa(len(response)) + "\r\n" + response + "\r\n"
}

func ReturnNullResponse() string {
	return "$-1\r\n"
}

func ReturnOkResponse() string {
	return "+OK\r\n"
}

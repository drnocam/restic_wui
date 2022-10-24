package main

import (
	"fmt"
)

type Message struct {
	Type    int
	Message string
}

type RJson struct {
	Message Message `json:"message"`
	Version int     `json:"vers"`
	Code    int     `json:"code"`
	Data    string  `json:"data"`
}

func JsonReturn(message Message, data string) string {
	ret := RJson{
		Message: message,
		Version: 1,
		Code:    200,
		Data:    data,
	}
	return fmt.Sprintf("{\"code\":%d, \"vers\":%d, \"message\":{\"Type\":%d,\"Message\":%q },\"data\":%s }", ret.Code, ret.Version, message.Type, message.Message, data)
}

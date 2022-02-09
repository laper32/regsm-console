package status

import (
	"encoding/json"
	"fmt"
)

type Code int

var (
	_message = map[int]string{}
)

func New(e int, msg string) Code {
	return add(e, msg)
}

func add(e int, msg string) Code {
	if _, ok := _message[e]; ok {
		panic(fmt.Sprintf("Code: %v already existd.", e))
	}
	_message[e] = msg

	return ToCode(e)
}

func ToCode(e int) Code { return Code(e) }

func (e Code) ToInt() int { return int(e) }

func (e Code) Message() string {
	if msg, ok := _message[e.ToInt()]; ok {
		return msg
	}
	return "Unknown Error code."
}

func (e Code) WriteDetail(detail interface{}) string {
	ret := make(map[string]interface{})
	ret["code"] = e
	ret["message"] = e.Message()
	if detail != nil {
		ret["detail"] = detail
	}
	output, _ := json.MarshalIndent(&ret, "", "    ") // we force 4 space rather than \t
	return string(output)
}

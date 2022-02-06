package status

import "fmt"

type Code int

var (
	_message = map[int]string{}
)

func New(e int, msg string) Code {
	if e <= 0 {
		panic("Reversed, and not allowed.")
	}
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

func Message(e int) string {
	if msg, ok := _message[e]; ok {
		return msg
	}
	return "Unknown Error"
}

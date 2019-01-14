package nerr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
)

type E struct {
	MessageLog []string
	Type       string
	Stack      []byte
}

func (e E) String() string {
	return e.Error()
}

func (e E) Error() string {
	b := strings.Builder{}
	for i := range e.MessageLog {
		b.WriteString(" - ")
		b.WriteString(e.MessageLog[i])
	}
	return b.String()
}

func (e *E) Addf(s string, v ...interface{}) *E {
	return e.Add(fmt.Sprintf(s, v...))
}

func (e *E) Add(s string) *E {
	e.MessageLog = append(e.MessageLog, strings.Trim(s, " \n\r"))
	return e
}

func (e *E) SetType(Type string) *E {
	e.Type = Type
	return e
}

func Translate(e error) *E {
	return Create(e.Error(), reflect.TypeOf(e).String())
}

func Create(msg string, Type string) *E {
	return &E{
		MessageLog: []string{msg},
		Type:       Type,
		Stack:      debug.Stack(),
	}
}

func Createf(Type string, format string, a ...interface{}) *E {
	return &E{
		MessageLog: []string{fmt.Sprintf(format, a...)},
		Type:       Type,
		Stack:      debug.Stack(),
	}
}

func (e *E) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		MessageLog []string `json:"message-log"`
		Type       string   `json:"type"`
		Stack      string   `json:"stack"`
	}{
		MessageLog: e.MessageLog,
		Type:       e.Type,
		Stack:      fmt.Sprintf("%s", e.Stack),
	})
}

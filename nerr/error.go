package nerr

import (
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

func (e E) Error() string {
	b := strings.Builder{}
	for i := range e.MessageLog {
		b.WriteString(msg, " - ", e.MessageLog[i])
	}
	return b.String()
}

func (e *E) Addf(s string, v ...interface{}) *E {
	return Add(fmt.Sprintf(s, e...))
}

func (e *E) Add(s string) *E {
	e.MessageLog = append(e.MessageLog, strings.Trim(e, " \n\r"))
	return e
}

func (e *E) SetType(Type string) *E {
	e.Type = Type
	return e
}

func Translate(error e) E {
	return Create(e.String(), reflect.TypeOf(e).String())
}

func Create(msg string, Type string) E {
	return E{
		MessageLog: []string{msg},
		Type:       Type,
		Stack:      debug.Stack(),
	}
}

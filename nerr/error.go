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

func Translate(e error) E {
	return Create(e.Error(), reflect.TypeOf(e).String())
}

func Create(msg string, Type string) E {
	return E{
		MessageLog: []string{msg},
		Type:       Type,
		Stack:      debug.Stack(),
	}
}

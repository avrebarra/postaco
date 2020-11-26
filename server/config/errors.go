package config

import "fmt"

type ErrServer struct {
	Code    ServerCode
	Message string
	Wrap    error
}

func (e ErrServer) Error() (msg string) {
	msg = fmt.Sprintf("%s (%s)", e.Message, e.Code)
	if e.Wrap != nil {
		msg += fmt.Sprintf(": %s", e.Wrap.Error())
	}
	return
}

func (e ErrServer) Unwrap() error {
	return e.Wrap
}

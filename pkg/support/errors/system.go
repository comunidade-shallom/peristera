package errors

import "strings"

type SystemError struct {
	BusinessError
	Reason error `json:"reason"`
}

func System(reason error, message, code string) SystemError {
	return SystemError{
		Reason:        reason,
		BusinessError: Business(message, code),
	}
}

func (e SystemError) WithErr(err error) SystemError {
	e.Reason = err

	return e
}

func (e SystemError) Error() string {
	var builder strings.Builder

	builder.WriteString(e.BusinessError.Error())
	builder.WriteString(" (")
	builder.WriteString(e.Reason.Error())
	builder.WriteString(")")

	return builder.String()
}

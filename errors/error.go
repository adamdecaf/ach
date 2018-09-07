// Copyright 2018 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package errors

import (
	"errors"
	"fmt"
	"strings"
)

type Error struct {
	underlying []error
}

func New(msg string) error {
	return WithError(errors.New(msg))
}

func WithError(err error) error {
	if err == nil {
		return nil
	}

	var es []error
	return &Error{
		underlying: append(es, err),
	}
}

func Contains(err error, slug string) bool {
	return err != nil && strings.Contains(err.Error(), slug)
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	var buf strings.Builder
	if e.underlying == nil {
		e.underlying = make([]error, 0)
	}
	for i := range e.underlying {
		if i > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(e.underlying[i].Error())
	}
	return buf.String()
}

// func (e *Error) Append(err error) *Error {
// 	if err == nil {
// 		return nil
// 	}
// 	e.underlying = append(e.underlying, err)
// 	return e
// }

// func (e *Error) Wrap(msg string) *Error {
// 	return e.Append(errors.New(msg))
// }

func File(field, value, msg string) error {
	if field == "" && value == "" {
		return New(fmt.Sprintf("file error message=%q", msg))
	}
	return New(fmt.Sprintf("file error field=%q value=%q message=%q", field, value, msg))
}

func Parse(line int, record string, err error) error {
	if record == "" {
		return New(fmt.Sprintf("line:%d %T %s", line, err, err))
	}
	return New(fmt.Sprintf("line:%d record:%s %T %s", line, record, err, err))
}

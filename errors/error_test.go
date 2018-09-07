// Copyright 2018 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package errors

import (
	"testing"
)

func TestError(t *testing.T) {
	err := New("foo")
	if err == nil {
		t.Fatal("should have error")
	}

	thing := func(_ error) {}
	thing(err) // compile test
}

// TODO(adam): test nil calls on methods

func TestContains(t *testing.T) {
	if Contains(nil, "foo") {
		t.Error("nil contains nothing")
	}
	if !Contains(New("foo"), "f") {
		t.Error("foo contains f")
	}
}

func TestErrorFormat(t *testing.T) {
	err := New("foobar")
	if err.Error() != "  foobar" {
		t.Errorf("got %q", err.Error())
	}
}

func TestFileError(t *testing.T) {
	err := File("mock", "", "test message")
	if err.Error() != `  file error field="mock" value="" message="test message"` {
		t.Errorf("got %q", err.Error())
	}
}

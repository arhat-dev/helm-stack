// Package pty provides functions for working with Unix terminals.
package pty

import (
	"errors"
	"io"
)

// ErrUnsupported is returned if a function is not
// available on the current platform.
var ErrUnsupported = errors.New("unsupported")

// Opens a pty and its corresponding tty.
func Open() (Pty, Tty, error) {
	return open()
}

type (
	FdHolder interface {
		Fd() uintptr
	}

	// Pty for terminal control in current process
	// for unix systems, the real type is *os.File
	// for windows, the real type is a *WindowsPty for ConPTY handle
	Pty interface {
		// Fd intended to resize Tty of child process in current process
		FdHolder

		// string writer is used to identify Pty and Tty
		io.StringWriter
		io.ReadWriteCloser
	}

	// Tty for data i/o in child process
	// for unix systems, the real type is *os.File
	// for windows, the real type is a *WindowsTty, which is a combination of two pipe file
	Tty interface {
		// Fd only intended for manual InheritSize from Pty
		FdHolder

		io.ReadWriteCloser
	}
)

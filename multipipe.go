// This code is a modified version of io.MultiWrite from Go standard library.
package main

import (
	"io"
	"log"
)

type multiWriteCloser struct {
	wcs []io.WriteCloser
}

func (t *multiWriteCloser) Write(p []byte) (n int, err error) {
	for _, w := range t.wcs {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

func (t *multiWriteCloser) Close() error {
	for _, w := range t.wcs {
		if err := w.Close(); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

// MultiWriteCloser is a WriteCloser compliant version of io.MultiWriter.
func MultiWriteCloser(writers ...io.WriteCloser) io.WriteCloser {
	allWriters := make([]io.WriteCloser, 0, len(writers))
	for _, w := range writers {
		if mw, ok := w.(*multiWriteCloser); ok {
			allWriters = append(allWriters, mw.wcs...)
		} else {
			allWriters = append(allWriters, w)
		}
	}
	return &multiWriteCloser{allWriters}
}

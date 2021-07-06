package main

import (
	"io"
	"reflect"
)

func (c *cmd) addout(o *cmd) {
	c.outlet[o] = struct{}{}
}

func iscmd(n node) bool {
	return reflect.TypeOf(n) == reflect.TypeOf(&cmd{})
}

// needsplay returns non-zero if node is runnable cmd
func needsplay(n node) int {
	if c, k := n.(*cmd); k {
		_, k := c.inlet.(*buf)
		if c.inlet == nil {
			return 1
		} else if k {
			return 2
		}
	}
	return 0
}

// stcap creates pipe connections for cmd with free inlet.
func stcap(c *cmd) {
	makepipes(nil, c)
}

// stretv creates pipe connections for cmd whose inlet is connected to a buf.
func stretv(c *cmd) {
	// connect buf via pipe
	r, w := io.Pipe()
	i := c.inlet.(*buf)
	if _, k := i.outlets[c]; k {
		i.out = w
	} else { // if link isn't broken being in this branch will mean that cmd is connected to sellet
		i.sel = w
	}
	makepipes(r, c)
}

func makepipes(r io.ReadCloser, s node) {
	// todo ensure that cmd loops are actually prohibited
	*s.Input() = r

	outsubs := make([]io.WriteCloser, 0)
	errsubs := make([]io.WriteCloser, 0)
	for o := range *s.Outlets() {
		// buf on outlet = end
		if _, da := o.(*buf); da {
			continue
		}
		ro, wo := io.Pipe()
		outsubs = append(outsubs, wo)
		makepipes(ro, o)
	}
	for o := range *s.Errlets() {
		if _, da := o.(*buf); da {
			continue
		}
		re, we := io.Pipe()
		errsubs = append(errsubs, we)
		makepipes(re, o)
	}
	*s.Primary() = MultiWriteCloser(outsubs...)
	*s.Secondary() = MultiWriteCloser(errsubs...)
}

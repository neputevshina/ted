package main

import "io"

func (c *cmd) addout(o *cmd) {
	c.outlet[o] = struct{}{}
}

// starters searches for all possible pipe starters
//
// cmd could be starter if:
//
// • inlet of cmd is not an another cmd
//
// • inlet of cmd is a buffer
//
// buffers aren't starable. so consider a special case:
// buffer is copied to another if they are connected
// inlet-to-outlet.
//
// this is a todo for you, loser
func starters(p []node) {
	cds := make([]node, 0, len(p))
	for _, c := range p {
		if *c.Inlet() == nil {
			cds = append(cds, c)
		}
		switch (*c.Inlet()).(type) {
		case *buf:
			// makepipes()
		case *cmd:
			continue
		}
	}
	// return cds
}

func makepipes(r io.ReadCloser, s node) {
	// todo ensure that cmd loops are actually prohibited
	*s.Input() = r

	outsubs := make([]io.WriteCloser, 0)
	for o := range *s.Outlets() {
		r, w := io.Pipe()
		outsubs = append(outsubs, w)
		makepipes(r, o)
	}
	*s.Primary() = MultiWriteCloser(outsubs...)

	errsubs := make([]io.WriteCloser, 0)
	for o := range *s.Outlets() {
		r, w := io.Pipe()
		errsubs = append(errsubs, w)
		makepipes(r, o)
	}
	*s.Secondary() = MultiWriteCloser(errsubs...)
}

package main

import "reflect"

func (c *cmd) addout(o *cmd) {
	c.outlet[o] = struct{}{}
}

func iscmd(n node) bool {
	return reflect.TypeOf(n) == reflect.TypeOf(&cmd{})
}

func findstarters(p []node) {
	// := make([]node, 0, len(p))
	for _, n := range p {
		if c, k := n.(*cmd); k {
			if _, k := c.inlet.(*buf); c.inlet == nil || k {

			}
		}
	}
}

// func starters(p []node) {
// 	cds := make([]node, 0, len(p))
// 	for _, c := range p {
// 		if *c.Inlet() == nil {
// 			cds = append(cds, c)
// 		}
// 		switch (*c.Inlet()).(type) {
// 		case *buf:
// 			// makepipes()
// 		case *cmd:
// 			continue
// 		}
// 	}
// 	// return cds
// }

// func makepipes(r io.ReadCloser, s node) {
// 	// todo ensure that cmd loops are actually prohibited
// 	*s.Input() = r

// 	outsubs := make([]io.WriteCloser, 0)
// 	for o := range *s.Outlets() {
// 		r, w := io.Pipe()
// 		outsubs = append(outsubs, w)
// 		makepipes(r, o)
// 	}
// 	*s.Primary() = MultiWriteCloser(outsubs...)

// 	errsubs := make([]io.WriteCloser, 0)
// 	for o := range *s.Outlets() {
// 		r, w := io.Pipe()
// 		errsubs = append(errsubs, w)
// 		makepipes(r, o)
// 	}
// 	*s.Secondary() = MultiWriteCloser(errsubs...)
// }

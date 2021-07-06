package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

var _ = node(&cmd{})

type cmd struct {
	Where XYWH
	Cmd   []rune

	inlet  node
	Entry  *TedText
	outlet map[node]struct{}
	errlet map[node]struct{}

	in  io.ReadCloser
	out io.WriteCloser
	err io.WriteCloser

	status int
	killme chan struct{}
}

func newcmd(where XYWH) *cmd {
	c := &cmd{
		Where:  where,
		outlet: make(map[node]struct{}),
		errlet: make(map[node]struct{}),
		Cmd:    make([]rune, 0),
	}
	c.Entry = NewTedText(&c.Cmd, G, gcache, true, false)
	lo, up := c.Limits()
	if lo.W > up.W {
		up = lo
	}
	c.Where.W, c.Where.H = up.Val()
	return c
}

func linemeasure(text []rune, font *FontSprites) (px int) {
	for _, r := range text {
		if r == '\n' {
			break
		}
		px += font.Cache[r].m.Advance
	}
	return
}

func (c *cmd) Limits() (WH, WH) {
	lim := Wt(linemeasure(c.Cmd, c.Entry.Sprites)+2*BoxKnobsSize, 3*FontSize/2)
	return Wt(72, 3*FontSize/2), lim
}

func (c *cmd) Draw() {
	boxdraw(c)
	c.Entry.Draw()
}

func (c *cmd) Mouse(at XY, buttons int, delta int) int {
	// todo: ugly
	x, y, w, h := c.Where.Val()
	c.Entry.Where = Rect(x+BoxKnobsSize, y, w-BoxKnobsSize*2, h)
	if knobpos(c.Where).Inside(at) {
		if buttons == MouseLeft {
			return MoveMe
		}
	}
	if killerpos(c.Where).Inside(at) {
		return OverKiller
	}
	if inletpos(c.Where).Inside(at) {
		return OverInlet
	}
	if outletpos(c.Where).Inside(at) {
		return OverOutlet
	}
	c.Entry.Mouse(at, buttons, delta)
	return 0
}

func (c *cmd) TextInput(r rune) {
	c.Entry.TextInput(r)
	lo, up := c.Limits()
	if lo.W > up.W {
		up = lo
	}
	c.Where.W, c.Where.H = up.Val()
}

func (c *cmd) Rect() *XYWH {
	return &c.Where
}

func (c *cmd) Play(finish chan struct{}) {
	cm, ar := c.parsecmd()
	ex := exec.Command(cm, ar...)
	ex.Stdin = c.in
	ex.Stdout = c.out
	ex.Stderr = c.err
	retur := make(chan struct{})

	er := ex.Start()
	if er != nil {
		log.Println(er)
		return
	}

	go func() {
		s, er := ex.Process.Wait()
		if er != nil {
			log.Println(er)
			return
		}
		c.status = s.ExitCode()
		// fixme
		// give reader time to read entirety of the pipe
		// WILL fail on large inputs
		time.Sleep(1 * time.Millisecond)
		retur <- struct{}{}
	}()

	select {
	case <-c.killme:
		log.Println(ex.Process.Kill())
		c.status = -127
	case <-retur:
	}
	switch o := ex.Stdout.(type) {
	case *os.File:
		break
	default:
		_ = o
		ex.Stdout.(io.WriteCloser).Close()
	}
	if finish != nil {
		finish <- struct{}{}
	}
}

func (c *cmd) parsecmd() (cm string, ar []string) {
	cmfull := false
	acc := ""
	p := 0
	for i, r := range c.Cmd {
		if r == ' ' {
			if !cmfull {
				cm = acc
				cmfull = true
				acc = ""
			} else {
				ar = append(ar, acc)
			}
			p++
			continue
		}
		acc += string(r)
		if i == len(c.Cmd)-1 {
			ar = append(ar, acc)
			return
		}
	}
	return
}

func (c *cmd) Input() *io.ReadCloser {
	return &c.in
}

func (c *cmd) Primary() *io.WriteCloser {
	return &c.out
}

func (c *cmd) Secondary() *io.WriteCloser {
	return &c.err
}

func (c *cmd) Inlet() *node {
	return &c.inlet
}

func (c *cmd) Outlets() *map[node]struct{} {
	return &c.outlet
}

func (c *cmd) Errlets() *map[node]struct{} {
	return &c.errlet
}

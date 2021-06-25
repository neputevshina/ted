package main

import (
	"io"
	"log"
	"os/exec"
)

type cmd struct {
	Where  XYWH
	inlet  node
	Entry  *TedText
	outlet map[node]struct{}
	errlet map[node]struct{}
	Cmd    []rune
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

func (c *cmd) Inlet() *node {
	return &c.inlet
}

func (c *cmd) Limits() (WH, WH) {
	lim := Wt(linemeasure(c.Cmd, c.Entry.Sprites)+2*BoxKnobsSize, 3*FontSize/2)
	return Wt(72, 3*FontSize/2), lim
}

func (c *cmd) Outlets() *map[node]struct{} {
	return &c.outlet
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

func (c *cmd) Play(in *io.PipeReader, out *io.PipeWriter, err *io.PipeWriter) (status chan int, killme chan struct{}) {
	cm, ar := c.parsecmd()
	ex := exec.Command(cm, ar...)
	ex.Stdin = in
	ex.Stdout = out
	ex.Stderr = err
	status = make(chan int, 1)
	killme = make(chan struct{}, 1)

	retur := make(chan struct{})
	go func() {
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
			retur <- struct{}{}
			status <- s.ExitCode()
		}()
		select {
		case <-killme:
			log.Println(ex.Process.Kill())
			retur <- struct{}{}
		}
		<-retur
		return
	}()
	return status, killme
}

func (c *cmd) parsecmd() (cm string, ar []string) {
	cmfull := false
	acc := ""
	p := 0
	for _, r := range c.Cmd {
		if r == ' ' {
			if !cmfull {
				cm = acc
				cmfull = true
			} else {
				ar = append(ar, acc)
			}
			p++
		}
		acc += string(r)
	}
	return
}

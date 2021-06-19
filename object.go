package main

type box struct {
	Where     [4]int
	Inlet     interface{}
	Text      []string
	Scrollpos uint
}

type node struct {
	Inlet interface{}
	Cmd   string
}

type tedstate struct {
	Objects []interface{}
	Pos     [2]int
}

var state = func() tedstate {
	return tedstate{}
}()

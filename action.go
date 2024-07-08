package main

type Action int

const (
	_ Action = iota
	USE
	CREATE
	EDIT
	DELETE
	LIST
)

var actionString = []string{
	"",
	"use",
	"create",
	"edit",
	"delete",
	"list",
}

var actionStringToAction = func() map[string]Action {
	m := make(map[string]Action)
	for i, s := range actionString {
		m[s] = Action(i)
	}
	return m
}()

func (a Action) IsValid() bool {
	return a > 0 && int(a) < len(actionString)
}

func (a Action) String() string {
	return actionString[a]
}

func getAction(s string) Action {
	return actionStringToAction[s]
}

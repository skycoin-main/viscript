package hypervisor

import (
	"github.com/corpusc/viscript/gfx"
	//"github.com/corpusc/viscript/msg"
)

type Task struct {
	In    [][]byte
	Out   [][]byte
	Panel *gfx.ScrollablePanel // if nil, no visual needed
}
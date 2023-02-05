package debugHandlers

import (
	"container/list"
	"fmt"
	"github.com/orchestd/debug"
)

type debugConfig struct {
	TraceRepo debug.TraceApi
}

type defaultDebugBuilder struct {
	ll *list.List
}

func DebugBuilder() debug.Builder {
	return &defaultDebugBuilder{ll: list.New()}
}

func (cr *defaultDebugBuilder) SetTraceRepo(resolver debug.TraceApi) debug.Builder {
	cr.ll.PushBack(func(cfg *debugConfig) {
		cfg.TraceRepo = resolver
	})
	return cr
}
func (cr *defaultDebugBuilder) Build() (debug.Debug, error) {
	debugCfg := &debugConfig{}
	for e := cr.ll.Front(); e != nil; e = e.Next() {
		f := e.Value.(func(cfg *debugConfig))
		f(debugCfg)
	}
	if debugCfg.TraceRepo == nil {
		return nil, fmt.Errorf("cannot initalize configurations without env settings")
	}
	return NewDebugInterface(debugCfg.TraceRepo), nil
}

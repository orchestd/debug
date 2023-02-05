package debugHandlers

import (
	"context"
	"github.com/orchestd/debug"
	"github.com/orchestd/servicereply"
)

type DebugInterface struct {
	TraceApi debug.TraceApi
}

func NewDebugInterface(traceApi debug.TraceApi) debug.Debug {
	return DebugInterface{TraceApi: traceApi}
}

func (i DebugInterface) GetErrorFromTrace(c context.Context, req debug.GetTraceErrorResponseReq) ([]debug.TraceError, servicereply.ServiceReply) {
	return i.TraceApi.GetErrorFromTrace(c, req)
}

package debug

import (
	"context"
	"github.com/orchestd/servicereply"
)

type TraceApi interface {
	GetErrorFromTrace(c context.Context, req GetTraceErrorResponseReq) ([]TraceError, servicereply.ServiceReply)
}
type Builder interface {
	SetTraceRepo(traceRepo TraceApi) Builder
	Build() (Debug, error)
}

type Debug interface {
	GetErrorFromTrace(c context.Context, req GetTraceErrorResponseReq) ([]TraceError, servicereply.ServiceReply)
}

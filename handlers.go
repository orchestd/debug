package debug

import (
	"github.com/orchestd/debug/dependencybundler/interfaces/transport"
)

func InitHandlers(router transport.IRouter, m Debug) {
	if m == nil || router == nil {
		return
	}
	router.POST(GetErrorFromTrace, transport.HandleFunc(m.GetErrorFromTrace))
}

const GetErrorFromTrace = "getErrorFromTrace"

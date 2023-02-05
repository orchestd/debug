package debug

import "encoding/json"

type GetTraceErrorResponseReq struct {
	TraceId string `json:"traceId"`

}

type TraceError struct {
	Status            string `json:"status"`
	UserMessageId     string `json:"userMessageId"`
	FullResponseTrace json.RawMessage `json:"fullResponseTrace"`
}

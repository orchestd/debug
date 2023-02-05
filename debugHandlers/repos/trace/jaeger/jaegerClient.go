package jaeger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/orchestd/debug"
	"github.com/orchestd/dependencybundler/interfaces/transport"
	"github.com/orchestd/servicereply"
)

type jaegerApiClient struct {
	client transport.HttpClient
}

func NewJaegerApiClient(client transport.HttpClient) debug.TraceApi {
	return jaegerApiClient{client: client}
}

type SpanNode struct {
	spanId    string
	tag       string
	MessageId string
	children  []*SpanNode
}
type Span struct {
	SpanNode Spans
	TagsMap  map[string]string
}

func (j jaegerApiClient) GetErrorFromTrace(c context.Context, req debug.GetTraceErrorResponseReq) ([]debug.TraceError, servicereply.ServiceReply) {
	var res Trace
	if reply := j.client.Get(c, "jaeger", fmt.Sprintf("api/traces/%s", req.TraceId), &res, nil); !reply.IsSuccess() {
		return nil, reply
	}

	spanMap := make(map[string][]Span)

	for _, data := range res.Data {
		for _, span := range data.Spans {
			tagsMap := make(map[string]string)

			for _, tag := range span.Tags {
				tagsMap[tag.Key] = fmt.Sprint(tag.Value)
			}
			if _, ok := tagsMap["span.kind"]; ok {
				if len(span.References) > 0 {
					if sps, ok := spanMap[span.References[0].SpanID]; !ok {
						spanMap[span.References[0].SpanID] = []Span{{SpanNode: span, TagsMap: tagsMap}}
					} else {
						spanMap[span.References[0].SpanID] = append(sps, Span{
							SpanNode: span,
							TagsMap:  tagsMap,
						})
					}
				} else {
					spanMap["head"] = []Span{{SpanNode: span, TagsMap: tagsMap}}
				}
			}
		}
	}
	var lastTraces []debug.TraceError
	getLastSpans(spanMap, "head", spanMap["head"][0], &lastTraces)
	return lastTraces, nil
}
func getLastSpans(spanMap map[string][]Span, spanID string, spanObj Span, lastTraces *[]debug.TraceError) {
	if s, ok := spanMap[spanID]; !ok {
		if val, ok := spanObj.TagsMap["span.kind"]; ok && val == "server" {
			if val, ok := spanObj.TagsMap["dependencyBundler.status"]; ok {
				var tr debug.TraceError
				tr.Status = val
				if val, ok := spanObj.TagsMap["dependencyBundler.id"]; ok {
					tr.UserMessageId = val
				}
				if len(spanObj.SpanNode.Logs) == 3 && len(tr.Status) > 0 && len(spanObj.SpanNode.Logs[2].Fields) > 0 {
					var jsonraw json.RawMessage
					if err := json.Unmarshal([]byte(spanObj.SpanNode.Logs[2].Fields[0].Value), &jsonraw); err == nil {
						tr.FullResponseTrace = jsonraw
					}
				}
				*lastTraces = append(*lastTraces, tr)
			}
		}
		return
	} else {
		for _, val := range s {
			getLastSpans(spanMap, val.SpanNode.SpanID, val, lastTraces)
		}
	}
}

type Trace struct {
	Data   []Data      `json:"data"`
	Total  int         `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Errors interface{} `json:"errors"`
}
type References struct {
	RefType string `json:"refType"`
	TraceID string `json:"traceID"`
	SpanID  string `json:"spanID"`
}
type Tags struct {
	Key   string      `json:"key"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
type Fields struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}
type Logs struct {
	Timestamp int64    `json:"timestamp"`
	Fields    []Fields `json:"fields"`
}
type Spans struct {
	TraceID       string       `json:"traceID"`
	SpanID        string       `json:"spanID"`
	Flags         int          `json:"flags"`
	OperationName string       `json:"operationName"`
	References    []References `json:"references"`
	StartTime     int64        `json:"startTime"`
	Duration      int          `json:"duration"`
	Tags          []Tags       `json:"tags"`
	Logs          []Logs       `json:"logs"`
	ProcessID     string       `json:"processID"`
	Warnings      interface{}  `json:"warnings"`
}
type P1 struct {
	ServiceName string `json:"serviceName"`
	Tags        []Tags `json:"tags"`
}
type P2 struct {
	ServiceName string `json:"serviceName"`
	Tags        []Tags `json:"tags"`
}
type Processes struct {
	P1 P1 `json:"p1"`
	P2 P2 `json:"p2"`
}
type Data struct {
	TraceID   string      `json:"traceID"`
	Spans     []Spans     `json:"spans"`
	Processes Processes   `json:"processes"`
	Warnings  interface{} `json:"warnings"`
}

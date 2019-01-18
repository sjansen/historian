package apigw

import (
	"bytes"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type responseWriter struct {
	buffer      bytes.Buffer
	header      http.Header
	status      int
	wroteHeader bool
}

func (w *responseWriter) Close() (resp events.APIGatewayProxyResponse, err error) {
	resp.Body = w.buffer.String()
	resp.Headers = make(map[string]string)
	resp.MultiValueHeaders = make(map[string][]string)
	for k, v := range w.Header() {
		if len(v) == 1 {
			resp.Headers[k] = v[len(v)-1]
		} else if len(v) > 1 {
			resp.MultiValueHeaders[k] = v
		}
	}
	resp.StatusCode = w.status
	return
}

func (w *responseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.buffer.Write(b)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.status = statusCode
	w.wroteHeader = true
}

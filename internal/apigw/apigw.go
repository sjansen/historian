package apigw

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ListenAndServe(addr string, h http.Handler) error {
	lambda.Start(func(
		ctx context.Context, req events.APIGatewayProxyRequest,
	) (
		res events.APIGatewayProxyResponse, err error,
	) {
		r, err := newRequest(req)
		if err != nil {
			return
		}

		w := &responseWriter{}
		h.ServeHTTP(w, r)

		return w.Close()
	})
	return nil
}

func newRequest(e events.APIGatewayProxyRequest) (*http.Request, error) {
	u, err := url.Parse(e.Path)
	if err != nil {
		return nil, err
	}

	q := url.Values(e.MultiValueQueryStringParameters)
	u.RawQuery = q.Encode()

	u.Host = e.Headers["Host"]
	u.Scheme = e.Headers["X-Forwarded-Proto"]
	if u.Host != "" && u.Scheme != "" {
		port, ok := e.Headers["X-Forwarded-Port"]
		switch {
		case ok && u.Scheme == "http" && port != "80":
			u.Host += ":" + port
		case ok && u.Scheme == "https" && port != "443":
			u.Host += ":" + port
		}
	}

	var r io.Reader
	if e.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(e.Body)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(b)
	} else {
		r = strings.NewReader(e.Body)
	}

	req, err := http.NewRequest(e.HTTPMethod, u.String(), r)
	if err != nil {
		return nil, err
	}

	if e.MultiValueHeaders != nil {
		req.Header = http.Header(e.MultiValueHeaders)
	} else {
		for k, v := range e.Headers {
			req.Header.Set(k, v)
		}
	}

	return req, nil
}

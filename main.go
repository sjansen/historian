package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/oklog/ulid"
)

var entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
var svc = dynamodb.New(session.Must(session.NewSession()))
var table = aws.String(os.Getenv("HISTORIAN_TABLE"))
var twoWeeks = time.Hour * 24 * 14

func handler(req events.ALBTargetGroupRequest) (resp events.ALBTargetGroupResponse, err error) {
	fmt.Printf("x-amzn-trace-id=%q.\n", req.Headers["x-amzn-trace-id"])

	msg := ""
	data := map[string]interface{}{
		"HTTPMethod": req.HTTPMethod,
		"Path":       req.Path,
		"Headers":    req.Headers,
		"Body":       req.Body,
		"ParsedBody": "",
	}
	if encoding, ok := req.Headers["content-type"]; ok && req.IsBase64Encoded {
		if encoding == "application/x-www-form-urlencoded" {
			if body, err := base64.StdEncoding.DecodeString(req.Body); err == nil {
				values, _ := url.ParseQuery(string(body))
				data["ParsedBody"] = values
				msg = values.Get("msg")
			} else {
				data["ParsedBody"] = err.Error()
			}
		}
	}

	now := time.Now().UTC()
	id := ulid.MustNew(ulid.Timestamp(now), entropy)
	timestamp := "0"
	expires := strconv.FormatInt(now.Add(twoWeeks).Unix(), 10)
	input := &dynamodb.PutItemInput{
		TableName: table,
		Item: map[string]*dynamodb.AttributeValue{
			"event-id":  {S: aws.String(id.String())},
			"timestamp": {N: aws.String(timestamp)},
			"expires":   {N: aws.String(expires)},
			"msg":       {S: aws.String(msg)},
		},
	}
	if _, err = svc.PutItem(input); err != nil {
		fmt.Printf("putitem=%q", err.Error())
	}

	var body bytes.Buffer
	if err = tmpl.Execute(&body, data); err != nil {
		return
	}

	resp.StatusCode = 200
	resp.Headers = map[string]string{"Content-Type": "text/html; charset=utf-8"}
	resp.Body = body.String()

	return
}

func main() {
	lambda.Start(handler)
}

var tmpl = template.Must(template.New("response").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>historian</title>
</head>
<body>
  <h1>{{ .HTTPMethod }} {{ .Path }}</h1>

  <h2>Headers</h2>
  <table>
    {{range $key, $val := .Headers}}<tr>
      <th><code>{{ $key }}</code></th>
      <td><code>{{ $val }}</code></td>
    </tr>
    {{end}}
  </table>

  <h2>Body</h2>
  <table>
    <tr>
      <th>Raw</th>
      <td><code>{{ .Body }}</code></td>
    </tr>
    <tr>
      <th>Parsed</th>
      <td><code>{{ .ParsedBody }}</code></td>
    </tr>
  </table>

  <form method="post">
    <div>
      <label for="msg">Message?</label>
      <input name="msg" id="msg" value="Spoon!">
    </div>
    <button>Send</button>
  </form>

</body>
</html>
`))

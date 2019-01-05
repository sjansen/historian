package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.ALBTargetGroupRequest) (resp events.ALBTargetGroupResponse, err error) {
	fmt.Printf("x-amzn-trace-id=%q.\n", req.Headers["x-amzn-trace-id"])

	var body bytes.Buffer
	if err = tmpl.Execute(&body, &req); err != nil {
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

  <form method="post">
    <div>
      <label for="msg">Message?</label>
      <input name="msg" id="msg" value="Spoon!">
    </div>
    <button>Send</button>
  </form>

  <h1>Metadata</h1>
  <table>
    <tr>
      <th>Method</th>
      <td>{{ .HTTPMethod }}</td>
    </tr>
    <tr>
      <th>Path</th>
      <td>{{ .Path }}</td>
    </tr>
  </table>

  <h1>Headers</h1>
  <table>
    {{range $key, $val := .Headers}}<tr>
      <th><code>{{ $key }}</code></th>
      <td><code>{{ $val }}</code></td>
    </tr>
    {{end}}
  </table>

  <h1>Body</h1>
  <pre>{{ .Body }}</pre>

</body>
</html>
`))

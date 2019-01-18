package api

import (
	"fmt"
	"net/http"
	"text/template"
)

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"HTTPMethod": r.Method,
		"Path":       r.URL.RawPath,
		"Headers":    r.Header,
	}

	headers := w.Header()
	headers.Add("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		fmt.Printf("error=%q\n", err.Error())
		w.WriteHeader(500)
	}
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
    {{range $key, $vals := .Headers}}{{range $idx, $val := $vals}}<tr>
      <th><code>{{ $key }}</code></th>
      <td><code>{{ $val }}</code></td>
    </tr>
    {{end}}{{end}}
  </table>

</body>
</html>
`))

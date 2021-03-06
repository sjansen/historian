package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/sjansen/historian/internal/dto"
	"github.com/sjansen/historian/internal/message"
)

type MessageRepo interface {
	Add(msg *dto.Message) (id string, err error)
}

type MessageVerifier interface {
	VerifySignature(message, signature string) (valid bool, err error)
}

type Handler struct {
	Repo     MessageRepo
	Verifier MessageVerifier
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGET(w, r)
	case "POST":
		h.handlePOST(w, r)
	default:
		w.WriteHeader(400)
	}
}

func (h *Handler) handleGET(w http.ResponseWriter, r *http.Request) {
	h.showRequestMetadata(r)
	headers := w.Header()
	headers.Add("Content-Type", "text/html; charset=utf-8")
	data := map[string]interface{}{
		"URL": r.URL.String(),
	}
	if err := tmpl.Execute(w, data); err != nil {
		fmt.Printf("error=%q\n", err.Error())
		w.WriteHeader(500)
	}
}

func (h *Handler) handlePOST(w http.ResponseWriter, r *http.Request) {
	h.showRequestMetadata(r)
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(400)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error=%q\n", err.Error())
		w.WriteHeader(500)
		return
	}

	data := string(body)
	signature := r.Header.Get("Clubhouse-Signature")
	valid, err := h.Verifier.VerifySignature(data, signature)
	switch {
	case err != nil:
		fmt.Printf("error=%q\n", err.Error())
		w.WriteHeader(500)
		return
	case !valid:
		w.WriteHeader(403)
		return
	}

	msg, err := message.Parse(strings.NewReader(data))
	if err != nil {
		fmt.Printf("error=%q\n", err.Error())
		w.WriteHeader(500)
	}

	id, err := h.Repo.Add(&dto.Message{
		Timestamp: time.Time(msg.ChangedAt),
		RawData:   data,
	})
	if err != nil {
		fmt.Printf("error=%q\n", err.Error())
		w.WriteHeader(500)
	}
	fmt.Printf("added=%q\n", id)

	w.WriteHeader(200)
}

func (h *Handler) showRequestMetadata(r *http.Request) {
	fmt.Printf("method=%q url=%q\n", r.Method, r.URL.String())
	headers := r.Header
	for k, vals := range headers {
		for _, v := range vals {
			fmt.Printf("%s=%q\n", k, v)
		}
	}
}

var tmpl = template.Must(template.New("response").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>historian</title>
</head>
<body>
  <h1>{{ .URL }}</h1>

  <form id="form" action="" method="post" style="padding:25px;">
<textarea id="data" name="data" style="width:300px;height:100px;border:1px solid;">{"foo": "bar"}</textarea>
    <div style="margin-top:10px;width:300px;text-align:center;">
      <input type="submit" id="submit">
    </div>
  </form>

<script>
(function(window){
  var data = window.document.getElementById("data");
  var submit = window.document.getElementById("submit");
  data.oninput = function(){
    try {
      JSON.parse(data.value);
      submit.disabled = false;
    } catch {
      submit.disabled = true;
    }
  }
  if (self.fetch) {
    var form = document.getElementById('form');
    form.onsubmit = function(){
      var data = new FormData(form);
      var msg = data.get("data");
      fetch("/messages", {
        method: "POST",
        headers: {
  	  'Content-Type': 'application/json'
        },
        body: msg,
      }).then(function(res){
        console.log(res);
      });
      return false;
    };
  }
})(window);
</script>
</body>
</html>
`))

package tplug

import (
	"bytes"
	"net/http"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type TPlug struct {
	Next httpserver.Handler
	*Config
}

func (tp *TPlug) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if httpserver.Path(r.URL.Path).Matches(tp.Config.Endpoint) {
		return tp.TPTest(w, r)
	}

	return tp.Next.ServeHTTP(w, r)
}

func (tp *TPlug) TPTest(w http.ResponseWriter, r *http.Request) (int, error) {
	q := r.URL.Query().Get("q")

	data := make(map[string]string)
	data["q"] = q

	var buf bytes.Buffer
	err := tp.Config.Template.Execute(&buf, data)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	buf.WriteTo(w)
	return http.StatusOK, nil
}

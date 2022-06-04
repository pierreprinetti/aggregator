package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
)

const html = `<html>
<head><title>My list</title></head>
<body>
<ul>{{ range . }}<li>{{ .Text }}</li>{{ end }}</ul>
</body>
</html>`

var listTemplate *template.Template

type store struct {
	data []Item
}

func NewStore() *store {
	return new(store)
}

type Item struct {
	Text string
}

func (s store) serveList(res http.ResponseWriter, req *http.Request) {
	if err := listTemplate.Execute(res, s.data); err != nil {
		fmt.Fprintf(res, "%v", err)
	}
}

func (s *store) serveAdd(res http.ResponseWriter, req *http.Request) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "%v", err)
	}
	s.data = append(s.data, Item{string(b)})
	res.WriteHeader(http.StatusOK)
}

func (s *store) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		s.serveList(res, req)
	case http.MethodPost:
		s.serveAdd(res, req)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func init() {
	t := template.New("list")
	t, err := t.Parse(html)
	if err != nil {
		panic(err)
	}
	listTemplate = t
}

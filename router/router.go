package router

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/ninnemana/cvcc-go-app/quotes"
	"github.com/pkg/errors"
)

var (
	listView = `
	<html><body>
		<h3>Current Quotes ( {{len .Quotes}} )</h3>
		<a href="/add">Add Quote</a>
		<div class="quotes">
		{{range .Quotes}}
			<blockquote>
				<p>{{.Quote}}</p>
				<footer>— {{.Author}}</footer>
			</blockquote>
		{{end}}
		</div>
	</body></html>
	`

	addView = `
	<html><body>
		<h3>Add Quote</h3>
		<div class="quotes">
			<form action="/put" method="post">
				<div>
					<label>Author</label>
					<input type="text" placeholder="Enter author name.." name="author" required />
				</div>
				<div>
					<label>Quote</label>
					<textarea rows="5" cols="10" placeholder="Enter quote.." name="quote" required></textarea>
				</div>
				<div>
					<button type="submit">Submit</button>
				</div>
			</form>
		</div>
	</body></html>
	`
)

type Router interface {
	Index(http.ResponseWriter, *http.Request)
	Add(http.ResponseWriter, *http.Request)
	Put(http.ResponseWriter, *http.Request)
	API(http.ResponseWriter, *http.Request)
}

type BasicRouter struct {
	quotes quotes.Interactor
}

func NewBasic() (Router, error) {

	quoteSvc, err := quotes.NewService()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create quotes service")
	}

	return &BasicRouter{
		quotes: quoteSvc,
	}, nil
}

func (b *BasicRouter) Index(w http.ResponseWriter, req *http.Request) {
	results, err := b.quotes.List(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	t := template.New("index")
	t, err = t.Parse(listView)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = t.Execute(w, map[string]interface{}{
		"Quotes": results,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (b *BasicRouter) Add(w http.ResponseWriter, req *http.Request) {
	t, err := template.New("add").Parse(addView)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := t.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (b *BasicRouter) Put(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	quote := &quotes.Quote{
		Author: req.FormValue("author"),
		Quote:  req.FormValue("quote"),
	}

	_, err := b.quotes.Put(context.Background(), quote)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, req, "/", http.StatusFound)
}

func (b *BasicRouter) API(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		b.apiCreate(w, req)
		return
	case http.MethodPost:
		b.apiUpdate(w, req)
		return
	case http.MethodDelete:
		b.apiDelete(w, req)
		return
	}

	id := req.URL.Query().Get("id")

	var res interface{}
	var err error
	switch id {
	case "":
		res, err = b.quotes.List(context.Background())
	default:
		res, err = b.quotes.Get(context.Background(), id)
	}

	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, PATCH, POST")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (b *BasicRouter) apiCreate(w http.ResponseWriter, req *http.Request) {
	var q quotes.Quote
	err := json.NewDecoder(req.Body).Decode(&q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := b.quotes.Put(context.Background(), &q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, PATCH, POST")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (b *BasicRouter) apiUpdate(w http.ResponseWriter, req *http.Request) {
	var q quotes.Quote
	err := json.NewDecoder(req.Body).Decode(&q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := b.quotes.Update(context.Background(), &q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, PATCH, POST")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (b *BasicRouter) apiDelete(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	if err := b.quotes.Delete(context.Background(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, PATCH, POST")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

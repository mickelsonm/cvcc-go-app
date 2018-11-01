package router

import (
	"context"
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

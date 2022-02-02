package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/asggo/random"
	"github.com/gorilla/mux"
)

const (
	notFoundError    = "Page not found"
	fibError     = "That's not right"
	webServerAddress = "127.0.0.1:3000"
	finalValue = 2111485077978050
	nextPuzzle = "http://127.0.0.1:3000/"
)

var (
	errTmpl  = template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/error.html"))
	idxTmpl  = template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/index.html"))
	fibTmpl  = template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/fib.html"))
	sessions = make(map[string]Fibonacci)
)

type Fibonacci struct {
	curr uint64
	next uint64
}

func (f *Fibonacci) Update() {
	next := f.curr + f.next
	f.curr = f.next
	f.next = next
}

func errorHandler(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	errTmpl.ExecuteTemplate(w, "layout", message)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	errorHandler(w, http.StatusNotFound, notFoundError)
}

func fib(w http.ResponseWriter, r *http.Request) {
	// Get the session and value submitted by the user
	vars := mux.Vars(r)
	sessId := vars["session"]
	val := vars["val"]

	// Ensure we have a valid session
	f, ok := sessions[sessId]
	if !ok {
		notFound(w, r)
		return
	}

	// Make sure we have not reached the final value
	if f.curr >= finalValue {
		http.Redirect(w, r, nextPuzzle, http.StatusFound)
		return
	}

	// Convert our value to a number and make sure it is next
	curr, _ := strconv.ParseUint(val, 10, 64)
	if curr != f.curr {
		errorHandler(w, http.StatusBadRequest, fibError)
		return
	}

	retVal := f.next

	// Update the Fibonacci value to take our turn
	f.Update()  // Take our turn
	f.Update()  // Take their turn

	sessions[sessId] = f

	w.Header().Add("next", fmt.Sprintf("%d", retVal))
	fibTmpl.ExecuteTemplate(w, "layout", retVal)
}

func index(w http.ResponseWriter, r *http.Request) {
	// Create a new session
	val, _ := random.Uint64()
	session := fmt.Sprintf("%d", val)

	// Create a new sequence and save it to our session
	f := Fibonacci{0, 1}
	sessions[session] = f

	// Create the starter link
	link := fmt.Sprintf("/%s/%d", session, f.curr)

	idxTmpl.ExecuteTemplate(w, "layout", link)
}

func main() {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/{session:[0-9]+}/{val:[0-9]+}", fib).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	srv := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         webServerAddress,
		Handler:      r,
	}

	err := srv.ListenAndServe()
	log.Fatalf("Web server closed: %s", err)
}

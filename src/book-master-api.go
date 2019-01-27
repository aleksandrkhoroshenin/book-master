package main

import (
	"../src/books"
	"../src/security"
	"../src/transport"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	ver = "v1"
)

var (
	booksService     books.HandlerBooks
	transportService transport.Router
	securityService   security.Security
)

var initFlag = flag.Bool("initial start", false, "Check your service")
var httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")

func main() {
	flag.Parse()

	db, err := sql.Open("postgres",
		"postgresql://booksDB@localhost:5432/db_1?sslmode=disable&user=postgres&password=Aebnm")
	if err != nil {
		return
	}

	initService(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/login", securityService.Login)
	mux.HandleFunc("/book", transportService.Handle)
	mux.HandleFunc("/book/:uuid", transportService.Handle)

	handler := AccessLogMiddleware(mux)

	s := &http.Server{
		Addr:           *httpAddr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func initService(db *sql.DB) {
	sessionService := security.NewSessionManager()
	booksService = books.CreateInstance(db)
	transportService = transport.CreateInstance(sessionService)
	securityService = security.CreateInstance(sessionService)
}

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("accessLogMiddleware", r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("[%s] %s, %s %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
	})
}

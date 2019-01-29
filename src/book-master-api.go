package main

import (
	"../src/books"
	"../src/security"
	"../src/transport"
	"../src/users"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

const (
	ver = "v1"
)

var (
	booksService     books.HandlerBooks
	userService      users.UserHandler
	transportService transport.Router
	securityService  security.Security
)

var initFlag = flag.Bool("initial start", false, "Check your service")
var httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")

func main() {
	flag.Parse()

	// TODO:: add timeout for docker
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		"postgres", "Aebnm", "db_1")
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer db.Close()
	initService(db)

	if *initFlag {
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/login", securityService.Login)
	mux.HandleFunc("/books", securityService.CheckSession(transportService.Handle))
	//mux.HandleFunc("/books/*", securityService.CheckSession(transportService.Handle))

	handler := AccessLogMiddleware(mux)
	//handler = securityService.CheckSession(handler)

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
	sessionService := security.NewSessionManager(db)
	booksService = books.CreateInstance(db)
	userService = users.CreateInstance(db)
	transportService = transport.CreateInstance(sessionService, booksService)
	securityService = security.CreateInstance(sessionService, userService)
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

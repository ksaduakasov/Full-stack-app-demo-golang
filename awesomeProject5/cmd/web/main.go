package main

import (
	"database/sql"
	"flag"
	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"se03/models/postgresql"
	"time"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session *sessions.Session
	snippets *postgresql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	dsn := flag.String("dsn", "user=kalbek password=pass dbname=snippetbox sslmode=disable", "MySQL data source name")
	addr := flag.String("addr", ":4000", "HTTP network address")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		session: session,
		snippets: &postgresql.SnippetModel{DB: db},
		templateCache: templateCache,
	}
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}


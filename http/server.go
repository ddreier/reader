package http

import (
	"embed"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"reader/data"
	"reader/http/assets"
)

const templatesDir = "templates"

var (
	//go:embed templates/* templates/layouts/*
	files     embed.FS
	templates map[string]*template.Template
)

type Server struct {
	ln     net.Listener
	server *http.Server
	router *mux.Router
	sc     *securecookie.SecureCookie

	// Bind address
	Addr string

	// Keys used for secure cookie encryption.
	HashKey  string
	BlockKey string

	// Data services
	Feeds       data.Feeds
	UnreadItems data.UnreadItems
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: mux.NewRouter(),
	}

	if err := LoadTemplates(); err != nil {
		panic(fmt.Sprintf("Failed to load templates: %s", err))
	}

	s.server.Handler = http.HandlerFunc(s.router.ServeHTTP)

	s.router.Use(loggingMiddleware)

	s.router.PathPrefix("/stylesheets/").
		Handler(http.FileServer(http.FS(assets.FS)))
	s.router.PathPrefix("/javascripts/").
		Handler(http.FileServer(http.FS(assets.FS)))
	s.router.PathPrefix("/fonts/").
		Handler(http.FileServer(http.FS(assets.FS)))
	s.router.PathPrefix("/img/").
		Handler(http.FileServer(http.FS(assets.FS)))

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, ok := templates["index.tmpl"]
		if !ok {
			log.Printf("Failed to load template index.tmpl")
			w.WriteHeader(500)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			log.Printf("Failed executing template index.tmpl: %s", err)
			w.WriteHeader(500)
			return
		}
	}).Methods("GET")

	s.router.HandleFunc("/feeds", func(w http.ResponseWriter, r *http.Request) {
		t, ok := templates["feeds.tmpl"]
		if !ok {
			log.Printf("Failed to load template feeds.tmpl")
			w.WriteHeader(500)
			return
		}

		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				log.Printf("POST /form: Unable to parse request body: %s", err)
				w.WriteHeader(500)
				return
			}

			// TODO going to have to figure out how to get a flash working w/ the cookie
			u, err := url.Parse(r.FormValue("feed_url"))
			if err != nil {
				log.Printf("POST /feed: error parsing feed url: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			_, err = s.Feeds.AddFeed("test", *u)
			if err != nil {
				log.Printf("POST /form: error adding feed to DB: %s", err)
				return
			}
		}

		feeds, err := s.Feeds.GetFeedList()
		if err != nil {
			log.Printf("Error getting feed list from DB: %s", err)
			w.WriteHeader(500)
			return
		}

		templateData := struct {
			Test  string
			Feeds []data.Feed
		}{
			Test:  "Feeds list!",
			Feeds: feeds,
		}

		if err := t.Execute(w, templateData); err != nil {
			log.Printf("Failed executing template feeds.tmpl: %s", err)
			w.WriteHeader(500)
			return
		}
	}).Methods("GET", "POST")

	// TODO handle POST /feeds

	s.router.HandleFunc("/feeds/new", func(w http.ResponseWriter, r *http.Request) {
		t, ok := templates["add_feed.tmpl"]
		if !ok {
			log.Printf("Failed to load template add_feed.tmpl")
			w.WriteHeader(500)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			log.Printf("Failed executing template add_feed.tmpl: %s", err)
			w.WriteHeader(500)
			return
		}
	})

	return s
}

func (s *Server) Open() (err error) {
	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
		return err
	}

	go func() {
		err := s.server.Serve(s.ln)
		if err != nil {
			log.Fatalf("Unable to start HTTP server: %s", err)
		}
	}()

	log.Printf("Server opened on %q", s.Addr)

	return nil
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func LoadTemplates() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		return err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(files, templatesDir+"/"+tmpl.Name(), templatesDir+"/layouts/*", templatesDir+"/partials/*")
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}
	return nil
}

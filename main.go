package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Store struct {
	root string
}

func main() {
	storeLoc := flag.String("store", "tmp", "where to store submissions")
	pubAddr := flag.String("pubaddr", ":8080", "public listen address")
	flag.Parse()

	s := Store{*storeLoc}

	pubHandler := http.NewServeMux()
	pubHandler.HandleFunc("/", s.onPost)
	pubS := &http.Server{
		Addr:           *pubAddr,
		Handler:        pubHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		pubS.ListenAndServe()
	}()
	<-make(chan struct{})
}

func (s *Store) onPost(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Origin, X-Session-ID")
	// Return early if it's CORS preflight.
	if req.Method == "OPTIONS" {
		return
	}

	if req.Method == "GET" {
		if req.URL.Query().Get("key") != os.Getenv("KEY") {
			rw.WriteHeader(403)
			return
		}
		files, err := os.ReadDir(s.root)
		if err != nil {
			rw.WriteHeader(500)
			return
		}
		for _, f := range files {
			rw.Write([]byte("<h1>" + f.Name() + "</h1>\n"))
			data, err := os.ReadFile(s.root + "/" + f.Name())
			if err != nil {
				continue
			}
			rw.Write(data)
			rw.Write([]byte("\n\n"))
		}
	} else if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		if err != nil {
			rw.WriteHeader(406)
			return
		}
		sum := sha256.Sum256(data)
		fname := fmt.Sprintf("%s/%x", s.root, sum)
		if err := os.WriteFile(fname, data, 0644); err == nil {
			if redir := os.Getenv("REDIRECT"); redir != "" {
				http.Redirect(rw, req, redir, http.StatusFound)
			} else {
				rw.WriteHeader(200)
			}
			return
		} else {
			rw.WriteHeader(500)
		}
	} else {
		rw.WriteHeader(406)
		return
	}
}

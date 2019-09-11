package spaserver

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/apex/log"
)

type SPAFileServer struct {
	fileSystem http.FileSystem
	fileServer http.Handler
}

func NewSPAFileServer(fileSystem http.FileSystem) *SPAFileServer {
	return &SPAFileServer{
		fileSystem: fileSystem,
		fileServer: http.FileServer(fileSystem),
	}
}

func (s *SPAFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if f, err := s.fileSystem.Open(path); err == nil {
		if err = f.Close(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.fileServer.ServeHTTP(w, r)
	} else if os.IsNotExist(err) {
		r.URL.Path = ""
		s.fileServer.ServeHTTP(w, r)
	} else {
		log.WithError(err).Error("file system open")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

package server

import "net/http"

type Server struct {
	http.Server
}

func handlerMessage(msg string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(msg))
	}
}

func New(addr string) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlerMessage("Hello from base"))

	return &Server{
		http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

package main

import (
	"log"
	"net/http"
)

type APIServer struct {
    listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
    return &APIServer{
        listenAddr: listenAddr,
    }
}

func (s *APIServer) Run() error {
    router := http.NewServeMux()

    router.HandleFunc("/readiness", MakeHandler(s.Readiness))
    router.HandleFunc("/err", MakeHandler(s.Err))
    
    server := http.Server{
        Addr: s.listenAddr,
        Handler: router,
    }

    log.Println("App is running on port:", s.listenAddr)

    return server.ListenAndServe()
}

func(s *APIServer) Readiness(res http.ResponseWriter, req *http.Request) error {
    return writeJSON(res, http.StatusOK, map[string]string{"msg":"All good"})
}

func(s *APIServer) Err(res http.ResponseWriter, req *http.Request) error {
    return InvalidJSON() 
}

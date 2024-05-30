package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIServer struct {
    listenAddr string
    Store Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
    return &APIServer{
        listenAddr: listenAddr,
        Store: store,
    }
}

func (s *APIServer) Run() error {
    router := http.NewServeMux()

    router.HandleFunc("/v1/users", MakeHandler(s.handleUser))
    router.HandleFunc("/readiness", MakeHandler(s.Readiness))
    router.HandleFunc("/err", MakeHandler(s.Err))
    
    server := http.Server{
        Addr: s.listenAddr,
        Handler: router,
    }

    log.Println("App is running on port:", s.listenAddr)

    return server.ListenAndServe()
}

func (s *APIServer) handleUser(res http.ResponseWriter, req *http.Request) error {
	if req.Method == "POST" {
		return s.HandleCreateUser(res, req)
	}

    return writeJSON(res, http.StatusBadRequest, map[string]string{"msg":req.Method + " is not allowed"})
}

func(s *APIServer) HandleCreateUser(res http.ResponseWriter, req *http.Request) error {
    var reqBody CreateUserRequest

    if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
        return InvalidJSON()
    }

    user, err := NewUser(reqBody.Name)
    if err != nil {
        return err
    }

    if err := s.Store.CreateUserToDb(user); err != nil {
        return err
    }

    return writeJSON(res, http.StatusOK, map[string]string{"msg":reqBody.Name + "successfully created"}) 
}

func(s *APIServer) Readiness(res http.ResponseWriter, req *http.Request) error {
    return writeJSON(res, http.StatusOK, map[string]string{"msg":"All good"})
}

func(s *APIServer) Err(res http.ResponseWriter, req *http.Request) error {
    return InvalidJSON() 
}

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wreckitral/RSS-feed-aggregator/internal/auth"
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

    if req.Method == "GET" {
        return s.HandleGetUser(res, req)
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

    createdUser, err := s.Store.CreateUserToDb(user)
    if err != nil {
        return err
    }

    resBody := UserResponse{
        ID: createdUser.ID,
        CreatedAt: createdUser.CreatedAt,
        UpdatedAt: createdUser.UpdatedAt,
        Name: createdUser.Name,
        APIKey: createdUser.ApiKey,
    }

    return writeJSON(res, http.StatusCreated, resBody)
}

func(s *APIServer) HandleGetUser(res http.ResponseWriter, req *http.Request) error {
    apiKey, err := auth.GetAPIKey(req.Header)   
    if err != nil {
        return UnauthorizedError(err)
    }

    user, err := s.Store.GetUserByAPIKey(apiKey)
    if err != nil {
        return err
    }

    resBody := UserResponse{
        ID: user.ID,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        Name: user.Name,
        APIKey: user.ApiKey,
    }

    return writeJSON(res, http.StatusAccepted, resBody)
}

func(s *APIServer) Readiness(res http.ResponseWriter, req *http.Request) error {
    return writeJSON(res, http.StatusOK, map[string]string{"msg":"All good"})
}

func(s *APIServer) Err(res http.ResponseWriter, req *http.Request) error {
    return InvalidJSON() 
}

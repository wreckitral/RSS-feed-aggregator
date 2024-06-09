package main

import (
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

    router.HandleFunc("/users", MakeHandler(s.handleUser))
    router.HandleFunc("/feeds", MakeHandler(s.handleFeed))
    router.HandleFunc("/feed_follows", MakeHandler(s.handleFeedFollows))
    router.HandleFunc("/feed_follows/{feedFollowID}", MakeHandler(s.handleFeedFollowsById))
    router.HandleFunc("/posts", MakeHandler(s.handlePosts))
    
    server := http.Server{
        Addr: "0.0.0.0:" + s.listenAddr,
        Handler: router,
    }

    log.Println("App is running on port:", s.listenAddr)

    return server.ListenAndServe()
}

/*  These following functions ensuring the right usage of http METHOD 
    with JSON response if theres an error, 
    and also act as a wrapper for authorized endpoints
*/

func (s *APIServer) handleUser(res http.ResponseWriter, req *http.Request) error {
    if req.Method == "POST" {
        return s.HandleCreateUser(res, req)
    }

    if req.Method == "GET" {
        s.middlewareAuth(s.HandleGetUser).ServeHTTP(res, req)
        return nil
    }

    return writeJSON(res, http.StatusBadRequest, map[string]string{"msg":req.Method + " is not allowed"})
}

func (s *APIServer) handleFeed(res http.ResponseWriter, req *http.Request) error {
    if req.Method == "GET" {
        return s.HandleGetFeeds(res, req)
    }

    if req.Method == "POST" {
        s.middlewareAuth(s.HandleCreateFeed).ServeHTTP(res, req)
        return nil
    }

    return writeJSON(res, http.StatusBadRequest, map[string]string{"msg":req.Method + " is not allowed"})
}

func (s *APIServer) handleFeedFollows(res http.ResponseWriter, req *http.Request) error {
    if req.Method == "POST" {
        s.middlewareAuth(s.HandleCreateFeedFollows).ServeHTTP(res, req)
        return nil
    }

    if req.Method == "GET" {
        s.middlewareAuth(s.HandleGetFeedFollows).ServeHTTP(res, req)
        return nil
    }

    return writeJSON(res, http.StatusBadRequest, map[string]string{"msg":req.Method + " is not allowed"})
}

func (s *APIServer) handleFeedFollowsById(res http.ResponseWriter, req *http.Request) error {
    if req.Method == "DELETE" {
        s.middlewareAuth(s.HandleDeleteFeedFollows).ServeHTTP(res, req)
        return nil
    }

    return writeJSON(res, http.StatusBadRequest, map[string]string{"msg":req.Method + " is not allowed"})
}

func (s *APIServer) handlePosts(res http.ResponseWriter, req *http.Request) error {
    if req.Method == "GET" {
        s.middlewareAuth(s.HandleGetPostsForUsers).ServeHTTP(res, req)
        return nil
    }

    return writeJSON(res, http.StatusBadRequest, map[string]string{"msg":req.Method + " is not allowed"})
}

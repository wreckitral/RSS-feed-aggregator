package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
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
    router.HandleFunc("/v1/feeds", MakeHandler(s.handleFeed))
    router.HandleFunc("/v1/feed_follows", MakeHandler(s.handleFeedFollows))
    router.HandleFunc("/v1/feed_follows/{feedFollowID}", MakeHandler(s.handleFeedFollowsById))
    
    server := http.Server{
        Addr: s.listenAddr,
        Handler: router,
    }

    log.Println("App is running on port:", s.listenAddr)

    return server.ListenAndServe()
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

    return writeJSON(res, http.StatusCreated, createdUser)
}

func(s *APIServer) HandleGetUser(res http.ResponseWriter, req *http.Request, user *User) error {
    return writeJSON(res, http.StatusOK, user)
}

func(s *APIServer) HandleCreateFeed(res http.ResponseWriter, req *http.Request, user *User) error {
    var reqBody CreateFeedRequest

    if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
        return InvalidJSON()
    }

    feed, err := NewFeed(reqBody.Name, reqBody.Url, user.ID)
    if err != nil {
        return err
    }

    createdFeed, err := s.Store.CreateFeedToDb(feed)
    if err != nil {
        return err
    }

    return writeJSON(res, http.StatusCreated, createdFeed)
}

func(s *APIServer) HandleGetFeeds(res http.ResponseWriter, req *http.Request) error {
    feedsFromDb, err := s.Store.GetFeeds()
    if err != nil {
        return err
    }
    if len(feedsFromDb) == 0 {
        return writeJSON(res, http.StatusNoContent, map[string]any{
            "statusCode": http.StatusNoContent,
            "msg": "no feeds yet",
        })
    }
    
    return writeJSON(res, http.StatusOK, feedsFromDb)
}

func(s *APIServer) HandleCreateFeedFollows(res http.ResponseWriter, req *http.Request, user *User) error {
    var reqBody FeedFollowRequest

    if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
        return InvalidJSON()
    }
    
    feedFollowtoDb := NewFeedFollow(user.ID, reqBody.FeedID)

    feedFollow, err := s.Store.CreateFeedFollows(feedFollowtoDb)
    if err != nil {
        return err
    }

    return writeJSON(res, http.StatusCreated, feedFollow)
}

func(s *APIServer) HandleGetFeedFollows(res http.ResponseWriter, req *http.Request, user *User) error {
    feeds, err := s.Store.GetFeedFollows(user.ID)
    if err != nil {
        return err
    }
    if len(feeds) == 0 {
        return writeJSON(res, http.StatusNoContent, map[string]any{
            "statusCode": http.StatusNoContent,
            "msg": "no feed follows yet",
        })
    }

    return writeJSON(res, http.StatusOK, feeds)
}

func(s *APIServer) HandleDeleteFeedFollows(res http.ResponseWriter, req *http.Request, user *User) error {
    params := req.PathValue("feedFollowID") 

    id := uuid.MustParse(params)

    if err := s.Store.DeleteFeedFollows(id, user.ID); err != nil {
        return err
    }

    return writeJSON(res, http.StatusOK, map[string]any{
        "statusCode": http.StatusOK,
        "msg": "feed with id: " + params + " successfully unfollowed",
    })
}

func (s *APIServer) middlewareAuth(handler authedHandler) http.HandlerFunc {
    return func(res http.ResponseWriter, req *http.Request) {
        apiKey, err := auth.GetAPIKey(req.Header)
        if err != nil {
            writeJSON(res, http.StatusForbidden, UnauthorizedError(err.Error()))
            return
        }

        user, err := s.Store.GetUserByAPIKey(apiKey)
        if err != nil {
            writeJSON(res, http.StatusForbidden, UnauthorizedError("api key is not recognized"))
            return
        }

        if err := handler(res, req, user); err != nil {
            writeJSON(res, http.StatusInternalServerError, err)
        }
    }
}

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

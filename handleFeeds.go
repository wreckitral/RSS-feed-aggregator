package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

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
        return NewAPIError(http.StatusBadRequest, err)
    }

    feedFollows := NewFeedFollow(createdFeed.UserID, createdFeed.ID)
    createdFeedFollows, err := s.Store.CreateFeedFollows(feedFollows)
    if err != nil {
        return err
    }

    return writeJSON(res, http.StatusCreated, map[string]any{
        "feed": createdFeed,
        "feed_follow": createdFeedFollows,
    })
}

func(s *APIServer) HandleGetFeeds(res http.ResponseWriter, req *http.Request) error {
    feedsFromDb, err := s.Store.GetFeeds()
    if err != nil {
        return err
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
        return NewAPIError(http.StatusBadRequest, err)
    }

    return writeJSON(res, http.StatusCreated, feedFollow)
}

func(s *APIServer) HandleGetFeedFollows(res http.ResponseWriter, req *http.Request, user *User) error {
    feeds, err := s.Store.GetFeedFollows(user.ID)
    if err != nil {
        return err
    }

    return writeJSON(res, http.StatusOK, feeds)
}

func(s *APIServer) HandleDeleteFeedFollows(res http.ResponseWriter, req *http.Request, user *User) error {
    params := req.PathValue("feedFollowID") 

    id, err := uuid.Parse(params)
    if err != nil {
        return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid params"))
    }

    if err := s.Store.DeleteFeedFollows(id, user.ID); err != nil {
        return NewAPIError(http.StatusBadRequest, err)
    }

    return writeJSON(res, http.StatusOK, map[string]any{
        "statusCode": http.StatusOK,
        "msg": "feed with id: " + params + " successfully unfollowed",
    })
}

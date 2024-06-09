package main

import (
    "net/http"
)

func (s *APIServer) HandleGetPostsForUsers(res http.ResponseWriter, req *http.Request, user *User) error {
    posts, err := s.Store.GetPostForUsers(user.ID, 10)
    if err != nil {
        return NewAPIError(http.StatusBadRequest, err)
    }

    return writeJSON(res, http.StatusOK, posts)
}

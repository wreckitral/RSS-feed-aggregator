package main

import (
    "net/http"
    "encoding/json"
)

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
        return NewAPIError(http.StatusBadRequest, err)
    }

    return writeJSON(res, http.StatusCreated, createdUser)
}

func(s *APIServer) HandleGetUser(res http.ResponseWriter, req *http.Request, user *User) error {
    return writeJSON(res, http.StatusOK, user)
}

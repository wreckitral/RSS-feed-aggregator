package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/wreckitral/RSS-feed-aggregator/internal/database"
)

type APIError struct {
    StatusCode int `json:"statusCode"`
    Msg any `json:"msg"`
}

func (e APIError) Error() string {
    return fmt.Sprintf("api error: %d", e.StatusCode) 
}

func NewAPIError(statusCode int, err error) APIError {
    return APIError{
        StatusCode: statusCode,
        Msg: err.Error(),
    }
}

func UnauthorizedError(err string) APIError {
    return NewAPIError(http.StatusForbidden, fmt.Errorf("%s", err))
}

func InvalidJSON() APIError {
    return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid request data"))
}

type APIFunc func(res http.ResponseWriter, req *http.Request) error

func MakeHandler(f APIFunc) http.HandlerFunc {
    return func(res http.ResponseWriter, req *http.Request) {
        if err := f(res, req); err != nil {
            if apiErr, ok := err.(APIError); ok {
                writeJSON(res, apiErr.StatusCode, apiErr)     
            } else {
                errResp := map[string]any{
                    "statusCode": http.StatusInternalServerError,
                    "msg": "internal server error",
                }
                writeJSON(res, http.StatusInternalServerError, errResp)
            }
            slog.Error("HTTP API error", "err", err.Error(), "path", req.URL.Path)
        }
    }
}

func writeJSON(res http.ResponseWriter, status int, v any) error {
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(status)

    return json.NewEncoder(res).Encode(v)
}

type authedHandler func(http.ResponseWriter, *http.Request, *database.User) error

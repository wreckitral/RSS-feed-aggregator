package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts API Key from
// headers of an HTTP request
// Example:
// Authorization: ApiKey {insert apiKey here}
func GetAPIKey(h http.Header) (string, error) {
    val := h.Get("Authorization")
    if val == "" {
        return "", errors.New("no authentication info")
    }

    vals := strings.Split(val, " ")
    if len(vals) != 2 {
        return "", errors.New("malformed auth header")
    }
    if vals[0] != "ApiKey" {
        return "", errors.New("malformed first part of auth header")
    }

    return vals[1], nil
}



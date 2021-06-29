package model

import "net/http"

const (
	StatusInternalServerError = "INTERNAL_SERVER_ERROR"
	StatusNotFoundError       = "NOT_FOUND_ERROR"
	StatusOK                  = "OK"
	StatusCreated             = "CREATED"
)

func GetHTTPStatusCodeByResponseStatus(responseStatus string) (httpSatusCode int) {
	switch responseStatus {
	case StatusOK:
		return http.StatusOK
	case StatusNotFoundError:
		return http.StatusNotFound
	case StatusCreated:
		return http.StatusCreated
	case StatusInternalServerError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

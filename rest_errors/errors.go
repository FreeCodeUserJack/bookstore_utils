package rest_errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RestError interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
}

type restError struct {
	ErrMessage string `json:"message"`
	ErrStatus  int    `json:"status"`
	ErrError   string `json:"error"`
	ErrCauses []interface{} `json:"causes"`
}

func (e restError) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: %v",
		e.ErrMessage, e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e restError) Message() string {
	return e.ErrMessage
}

func (e restError) Status() int {
	return e.ErrStatus
}

func (e restError) Causes() []interface{} {
	return e.ErrCauses
}

func NewRestError(message string, status int, err string, causes []interface{}) RestError {
	return restError{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
		ErrCauses:  causes,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestError, error) {
	var apiErr restError
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

func NewBadRequestError(message string) RestError {
	return restError{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "bad_request",
	}
}

func NewNotFoundError(message string) RestError {
	return restError{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not_found",
	}
}

func NewUnauthorizedError(message string) RestError {
	return restError{
		ErrMessage: message,
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   "unauthorized",
	}
}

func NewInternalServerError(message string, err error) RestError {
	result := restError{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "internal_server_error",
	}
	if err != nil {
		result.ErrCauses = append(result.ErrCauses, err.Error())
	}
	return result
}
package goatling

import "net/http"

type ServerResponse struct {
	Status  int
	Content interface{}
}

func InternalServerError(any interface{}) *ServerResponse {
	return &ServerResponse{
		Status:  http.StatusInternalServerError,
		Content: any,
	}
}

func BadGateway(any interface{}) *ServerResponse {
	return &ServerResponse{
		Status:  http.StatusBadGateway,
		Content: any,
	}
}

func Conflict(any interface{}) *ServerResponse {
	return &ServerResponse{
		Status:  http.StatusConflict,
		Content: any,
	}
}

func Created(any interface{}) *ServerResponse {
	return &ServerResponse{
		Status:  http.StatusCreated,
		Content: any,
	}
}

func Accepted(any interface{}) *ServerResponse {
	return &ServerResponse{
		Status:  http.StatusAccepted,
		Content: any,
	}
}

func Forbidden(any interface{}) *ServerResponse {
	return &ServerResponse{
		Status:  http.StatusForbidden,
		Content: any,
	}
}

func NoContent(any interface{}) *ServerResponse {
	return &ServerResponse{
		Status:  http.StatusNoContent,
		Content: any,
	}
}

func OK(any interface{}) *ServerResponse {
	return &ServerResponse{
		Status:  http.StatusOK,
		Content: any,
	}
}

func BadRequest(any interface{}) *ServerResponse {
	return &ServerResponse{
		Status:  http.StatusBadRequest,
		Content: any,
	}
}

func Unauthorized(any interface{}) *ServerResponse {
	return &ServerResponse{
		Status:  http.StatusUnauthorized,
		Content: any,
	}
}

func NotFound(any interface{}) *ServerResponse {
	return &ServerResponse{
		Status:  http.StatusNotFound,
		Content: any,
	}
}

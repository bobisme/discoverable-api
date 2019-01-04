package errors

import (
	"net/http"
)

func HttpStatus(e error) int {
	class := Internal
	if err, ok := e.(Error); ok {
		class = err.Class()
	}
	switch class {
	case Authentication:
		return http.StatusUnauthorized
	case Permission:
		return http.StatusForbidden
	case Client:
		return http.StatusBadRequest
	case Validation:
		return http.StatusUnprocessableEntity
	}
	return http.StatusInternalServerError
}

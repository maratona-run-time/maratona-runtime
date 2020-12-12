package utils

import (
	"fmt"
	"net/http"
)

// WriteResponse is used to write a HTTP response status in case of an error
func WriteResponse(rs http.ResponseWriter, status int, msg string, err error) {
	rs.WriteHeader(status)
	errorMsg := msg
	if err != nil {
		errorMsg = fmt.Sprintf("%v:\n%v", msg, err.Error())
	}
	rs.Write([]byte(errorMsg))
}

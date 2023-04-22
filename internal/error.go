package internal

import (
	"fmt"
	"log"
	"strings"
)

var errorCodes = map[int]string{
	502: "Failed requesting timetable data. Requested time might be too far in the future or past or other request parameters might be invalid.",
	503: "Failed requesting timetable data. Possibly too many requests. Try again later.",
}

type ErrorWithCode interface {
	Error() string
	ErrorCode() int
}
type ErrorCode struct {
	Msg  string
	Code int
}

func (e *ErrorCode) Error() string {
	if e.Msg != "" {
		return fmt.Sprintf("%d: %v", e.Code, e.Msg)
	}
	if val, ok := errorCodes[e.Code]; ok {
		return fmt.Sprintf("%d: %v", e.Code, val)
	}
	return fmt.Sprintf("%d: Unknown Error", e.Code)
}

func (e *ErrorCode) ErrorCode() int {
	return e.Code
}

func HandleError(e error) *ErrorCode {
	log.Print("Error: ", e)
	if strings.Contains(e.Error(), "HAFAS") {
		return &ErrorCode{Code: 502}
	}
	return &ErrorCode{Code: 503}
}
